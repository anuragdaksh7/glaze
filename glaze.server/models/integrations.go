package models

import (
	"errors"
	"glaze/pkg/crypto"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var encryptionKey []byte

func SetEncryptionKey(key string) {
	encryptionKey = []byte(key)
}

type Integration struct {
	Base
	WorkspaceID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Provider       string    `json:"provider"`   // "github"
	ProviderID     string    `json:"providerId"` // GitHub User/Org ID
	AccessToken    string    `json:"-"`          // ENCRYPT THIS
	RefreshToken   string    `json:"-"`          // ENCRYPT THIS
	ExpiresAt      time.Time `json:"expiresAt"`
	InstallationID string    `json:"installationId"` // If using GitHub Apps
	Workspace      Workspace `gorm:"foreignKey:WorkspaceID"`
}

func (i *Integration) BeforeSave(tx *gorm.DB) (err error) {
	if len(encryptionKey) == 0 {
		return errors.New("encryption key is not configured")
	}

	if i.AccessToken != "" {
		i.AccessToken, err = crypto.Encrypt(i.AccessToken, encryptionKey)
		if err != nil {
			return err
		}
	}

	if i.RefreshToken != "" {
		i.RefreshToken, err = crypto.Encrypt(i.RefreshToken, encryptionKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// AfterFind is called whenever a record is retrieved from the DB
func (i *Integration) AfterFind(tx *gorm.DB) (err error) {
	if len(encryptionKey) == 0 {
		return errors.New("encryption key is not configured")
	}

	if i.AccessToken != "" {
		i.AccessToken, err = crypto.Decrypt(i.AccessToken, encryptionKey)
		if err != nil {
			// If decryption fails, the key might be wrong or data is corrupt
			return err
		}
	}

	if i.RefreshToken != "" {
		i.RefreshToken, err = crypto.Decrypt(i.RefreshToken, encryptionKey)
		if err != nil {
			return err
		}
	}
	return nil
}
