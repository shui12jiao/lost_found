package api

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoSession         = errors.New("no session")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrPermissionDenied  = errors.New("permission denied")
	ErrInvalidPermission = errors.New("invalid permission")
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
