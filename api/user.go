package api

import "github.com/gin-gonic/gin"

type addUserRequest struct {
}

type addUserResponse struct {
}

type loginUserRequest struct {
	Code string `json:"code" binding:"required,alphanum"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        addUserResponse
}

type getUserRequest struct{}

type getUserResponse struct{}

type listUserRequest struct{}

type listUserResponse struct{}

type updateUserRequest struct{}

type updateUserResponse struct{}

type deleteUserRequest struct{}

type deleteUserResponse struct{}

func (server *Server) loginUser(ctx *gin.Context) {

}

func (server *Server) addUser(ctx *gin.Context) {

}

func (server *Server) getUser(ctx *gin.Context) {

}

func (server *Server) listUser(ctx *gin.Context) {

}

func (server *Server) deleteUser(ctx *gin.Context) {

}

func (server *Server) updateUser(ctx *gin.Context) {

}
