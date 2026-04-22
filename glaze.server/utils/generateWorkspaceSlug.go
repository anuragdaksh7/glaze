package utils

import (
	"fmt"
	"glaze/models"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

func GenerateWorkspaceSlug(name string) string {
	// lowercase
	slug := strings.ToLower(name)

	// replace non-alphanumeric with dash
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")

	// trim leading/trailing dashes
	slug = strings.Trim(slug, "-")

	return slug
}

func GenerateUniqueSlug(db *gorm.DB, name string) string {
	base := GenerateWorkspaceSlug(name)
	slug := base

	var count int64
	i := 1

	for {
		db.Model(&models.Workspace{}).
			Where("slug = ?", slug).
			Count(&count)

		if count == 0 {
			break
		}

		slug = fmt.Sprintf("%s-%d", base, i)
		i++
	}

	return slug
}
