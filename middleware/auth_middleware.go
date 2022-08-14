package middleware

import (
	"errors"
	"fmt"
	"lost_found/middleware/session"
	"lost_found/middleware/token"
	"lost_found/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"

	CookieName              = "session_id"
	AuthorizationPayloadKey = "authorization_payload"
	SessionHeaderKey        = "session"
	ManagerHeaderKey        = "manager"
)

func AuthSessionMiddleware(sessionManager *session.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID := ctx.GetHeader(CookieName)
		if len(sessionID) == 0 {
			err := errors.New("session id is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		session, err := sessionManager.ReadSession(sessionID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(SessionHeaderKey, session)
		ctx.Next()
	}
}

func AuthJWTMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func ManagerMiddleware(sessionManager *session.Manager, store service.Store) gin.HandlerFunc { //TODO
	return func(ctx *gin.Context) {
		sessionID := ctx.GetHeader(CookieName)
		if len(sessionID) == 0 {
			err := errors.New("session id is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		session, err := sessionManager.ReadSession(sessionID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		openId, ok := session.Value["open_id"]
		if !ok {
			err := errors.New("open_id is not provided")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		manager, err := store.GetManagerByOpenid(ctx, openId)
		if err != nil {
			err := errors.New("manager is not found")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(SessionHeaderKey, session)
		ctx.Set(ManagerHeaderKey, manager)
		ctx.Next()
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
