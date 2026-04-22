package utils

import (
	"errors"
	"glaze/models"

	"github.com/gin-gonic/gin"
)

func ExtractUser(c *gin.Context) (models.User, error) {
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, errors.New("user not found in context")
	}

	currentUser := user.(models.User)

	return currentUser, nil
}
