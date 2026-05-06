package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"glaze/internal/tasks"
	"glaze/logger"
	"glaze/models"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/hibiken/asynq"
	"github.com/moby/moby/client"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//"github.com/moby/go-archive"
//"github.com/moby/moby/api/types"
//"github.com/moby/moby/client"

type BuildWorker struct {
	db     *gorm.DB
	docker *client.Client
}

func NewBuildWorker(db *gorm.DB) (*BuildWorker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &BuildWorker{
		db:     db,
		docker: cli,
	}, nil
}

func generateSubdomain(repoFullName string) string {
	// 1. Lowercase everything
	slug := strings.ToLower(repoFullName)

	// 2. Replace slashes and dots with dashes
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, ".", "-")

	// Result: "anuragdaksh7/gh-sdk" becomes "anuragdaksh7-gh-sdk"
	return slug
}

func (w *BuildWorker) failDeployment(deploymentID string, errMsg string) error {
	w.db.Model(&models.Deployment{}).Where("id = ?", deploymentID).
		Updates(map[string]interface{}{
			"status":      models.StatusFailed,
			"logs":        errMsg,
			"finished_at": time.Now(),
		})
	return nil
}

func (w *BuildWorker) cleanupOldContainer(ctx context.Context, containerName string) {
	// 1. Try to stop the container (giving it a 10-second grace period)
	timeout := 10
	stopOptions := container.StopOptions{Timeout: &timeout}

	err := w.docker.ContainerStop(ctx, containerName, stopOptions)
	if err != nil {
		// If it's already stopped or doesn't exist, we don't care, just log it
		logger.Logger.Info("Container stop skipped or failed", zap.String("name", containerName))
	}

	// 2. Remove the container and its associated volumes
	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true, // This kills it even if the stop failed
	}

	if err := w.docker.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		logger.Logger.Info("Container removal skipped", zap.String("name", containerName))
	}
}

func getSafeContainerName(projectID string) string {
	return "glaze-app-" + projectID
}

func resolveBuildRoot(buildDir, rootDir string) string {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" || rootDir == "/" || rootDir == "." {
		return buildDir
	}

	rootDir = strings.TrimPrefix(rootDir, "./")
	rootDir = strings.TrimPrefix(rootDir, ".\\")
	rootDir = strings.TrimLeft(rootDir, "/\\")
	rootDir = filepath.Clean(rootDir)

	if rootDir == "." || rootDir == string(filepath.Separator) {
		return buildDir
	}

	return filepath.Join(buildDir, rootDir)
}

func (w *BuildWorker) deployContainer(ctx context.Context, imageName string, projectID string, deploymentID string, subdomain string, logBuffer bytes.Buffer) error {
	labels := map[string]string{
		"caddy":               fmt.Sprintf("%s.localhost", subdomain),
		"caddy.reverse_proxy": "{{upstreams 3000}}",
	}

	config := &container.Config{
		Image:  imageName,
		Labels: labels,
	}

	hostConfig := &container.HostConfig{
		NetworkMode: "glaze-network",
		// PortBindings: nil,
	}

	resp, err := w.docker.ContainerCreate(ctx, config, hostConfig, nil, nil, getSafeContainerName(projectID))
	if err != nil {
		w.failDeployment(deploymentID, "Container creation failed: "+err.Error())
		return err
	}

	if err := w.docker.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		w.failDeployment(deploymentID, "Container start failed: "+err.Error())
		return err
	}

	w.db.Model(&models.Deployment{}).Where("id = ?", deploymentID).Updates(map[string]interface{}{
		"status":       models.StatusSuccess,
		"logs":         logBuffer.String(),
		"finished_at":  time.Now(),
		"container_id": resp.ID,
		"image_name":   imageName,
	})

	return nil
}

func (w *BuildWorker) ProcessBuildTask(ctx context.Context, t *asynq.Task) error {
	var p tasks.BuildPayload
	json.Unmarshal(t.Payload(), &p)

	logger.Logger.Info("payload", zap.Any("payload", p))

	startedAt := time.Now()
	w.db.Model(&models.Deployment{}).Where("id = ?", p.DeploymentID).
		Updates(map[string]interface{}{
			"status":     models.StatusCloning,
			"started_at": startedAt,
		})

	buildDir := filepath.Join(os.TempDir(), p.DeploymentID)

	if err := os.RemoveAll(buildDir); err != nil {
		log.Printf("Failed to clean build path: %v", err)
		return err // Let Asynq retry if we can't even delete a folder
	}
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		return err
	}

	var deploy models.Deployment
	w.db.Model(&models.Deployment{}).Where("id = ?", p.DeploymentID).First(&deploy)

	var project models.Project
	err := w.db.Where("id = ?", deploy.ProjectID).First(&project).Error
	if err != nil {
		w.failDeployment(p.DeploymentID, "Associated project not found.")
		return err
	}

	var integration models.Integration
	err = w.db.Where("workspace_id = ? AND provider = ?", project.WorkspaceID, "github").First(&integration).Error
	if err != nil {
		w.failDeployment(p.DeploymentID, "GitHub integration not found. Please connect your account.")
		return nil
	}

	//repoURL := fmt.Sprintf("https://github.com/%s.git", p.RepoFullName)
	repoURL := fmt.Sprintf("https://%s@github.com/%s.git", integration.AccessToken, p.RepoFullName)

	branch := p.Branch
	if branch == "" {
		branch = deploy.Branch
	}
	if branch == "" {
		branch = project.DeployBranch
	}

	cloneArgs := []string{"clone", "--depth", "1"}
	if branch != "" {
		cloneArgs = append(cloneArgs, "--branch", branch)
	}
	cloneArgs = append(cloneArgs, repoURL, buildDir)

	cmd := exec.CommandContext(ctx, "git", cloneArgs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errorMessage := stderr.String()
		//errorMessage = strings.ReplaceAll(errorMessage, , "********

		cleanError := strings.ReplaceAll(errorMessage, integration.AccessToken, "********")

		logger.Logger.Error("Git Clone Error: ", zap.String("error", cleanError))
		w.failDeployment(p.DeploymentID, cleanError)
		return nil
	}

	defer os.RemoveAll(buildDir)

	if p.CommitSHA != "" {
		cmd = exec.CommandContext(ctx, "git", "-C", buildDir, "checkout", p.CommitSHA)
		stderr.Reset()
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			cleanError := strings.ReplaceAll(stderr.String(), integration.AccessToken, "********")
			logger.Logger.Error("Git checkout error", zap.String("error", cleanError))
			w.failDeployment(p.DeploymentID, cleanError)
			return nil
		}
	}

	//tarStream, err := archive.TarWithOptions(buildDir, &archive.TarOptions{})
	//if err != nil {
	//	return err
	//}

	w.db.Model(&models.Deployment{}).Where("id = ?", p.DeploymentID).Update("status", models.StatusBuilding)

	w.db.Model(&models.Deployment{}).Where("id = ?", p.DeploymentID).First(&deploy)
	buildRoot := resolveBuildRoot(buildDir, project.RootDirectory)

	imageName := strings.ReplaceAll(strings.ToLower(p.RepoFullName), "/", "-")

	cmd = exec.CommandContext(ctx, "nixpacks", "build", buildRoot, "--name", imageName)

	var logBuffer bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &logBuffer)
	cmd.Stderr = io.MultiWriter(os.Stderr, &logBuffer)

	if err := cmd.Run(); err != nil {
		logger.Logger.Error("NixPacks Build Error: ", zap.Any("error", err))
		w.failDeployment(p.DeploymentID, "Nixpacks build failed: "+logBuffer.String())
		return err
	}

	w.db.Model(&models.Deployment{}).Where("id = ?", p.DeploymentID).Updates(map[string]interface{}{
		"status":         models.StatusSuccess,
		"logs":           logBuffer.String(),
		"finished_at":    time.Now(),
		"build_duration": int64(time.Since(startedAt).Seconds()),
	})

	safeContainerName := getSafeContainerName(deploy.ProjectID.String())
	logger.Logger.Info("Cleaning up old deployment if it exists...")
	w.cleanupOldContainer(ctx, safeContainerName)

	err = w.deployContainer(ctx, imageName, deploy.ProjectID.String(), p.DeploymentID, generateSubdomain(p.RepoFullName), logBuffer)
	if err != nil {
		logger.Logger.Error("Deployment skipped or failed", zap.String("name", deploy.ProjectID.String()))
		return err
	}

	//containerConfig := &container.Config{
	//	Image: imageName,
	//	ExposedPorts: nat.PortSet{
	//		"3000/tcp": struct{}{},
	//	},
	//	Labels: map[string]string{
	//		"managed_by": "glaze",
	//		"project_id": deploy.ProjectID.String(),
	//	},
	//}

	//extPort := 3000

	//hostConfig := &container.HostConfig{
	//	PortBindings: nat.PortMap{
	//		"3000/tcp": []nat.PortBinding{
	//			{
	//				HostIP:   "0.0.0.0",
	//				HostPort: strconv.Itoa(extPort),
	//			},
	//		},
	//	},
	//	RestartPolicy: container.RestartPolicy{Name: "always"},
	//}

	//res, err := w.docker.ImageBuild(ctx, tarStream, types.ImageBuildOptions{
	//	Tags:        []string{imageName + ":latest"},
	//	Dockerfile:  "Dockerfile",
	//	Remove:      true,
	//	ForceRemove: true,
	//})
	//if err != nil {
	//	w.failDeployment(p.DeploymentID, "Docker build failed: "+err.Error())
	//	return err
	//}
	//defer res.Body.Close()
	//
	//w.db.Model(&models.Deployment{}).Where("id = ?", p.DeploymentID).
	//	Updates(map[string]interface{}{
	//		"status":      models.StatusSuccess,
	//		"finished_at": time.Now(),
	//	})

	return nil
}
