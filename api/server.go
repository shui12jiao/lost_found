package api

import (
	"fmt"
	"log"
	"lost_found/middleware"
	"lost_found/token"
	"lost_found/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if _, ok := binding.Validator.Engine().(*validator.Validate); !ok {
		log.Fatal("failed to init validator")
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users/login", server.loginUser)

	//普通用户使用功能
	authRoutes := router.Group("/").Use(middleware.AuthMiddleware(server.tokenMaker))
	//用户账户
	authRoutes.GET("/users/:id", server.getUser)
	authRoutes.DELETE("/users/:id", server.deleteUser)
	authRoutes.PATCH("/users/:id", server.updateUser)
	//位置和类型
	authRoutes.GET("/locations", server.listLocation)
	authRoutes.GET("/types", server.listType)
	//拾取物品
	authRoutes.GET("/founds", server.listFound)
	authRoutes.GET("/founds/:id", server.getFound)
	authRoutes.POST("/founds/add", server.addFound)
	authRoutes.DELETE("/founds/delete/:id", server.deleteFound) //只能删除自己发布的拾取物
	//遗失物品
	authRoutes.GET("/losts", server.listLost)
	authRoutes.GET("/losts/:id", server.getLost)
	authRoutes.POST("/losts/add", server.addLost)
	authRoutes.DELETE("/losts/delete/:id", server.deleteLost) //只能删除自己的发布的遗失物
	//归还物品
	authRoutes.GET("/matches", server.listMatch) //自己遗失或拾取的已归还物品

	//管理员系统
	manRoutes := router.Group("/manager").Use(middleware.ManagerMiddleware(server.tokenMaker))
	//用户账户
	manRoutes.POST("/users/add", server.addUser)
	//位置和类型
	manRoutes.GET("/locations", server.listLocation)
	manRoutes.POST("/locations/add", server.addLocation)
	manRoutes.DELETE("/locations/delete/:id", server.deleteLocation)
	manRoutes.GET("/types", server.listType)
	manRoutes.POST("/types/add", server.addType)
	manRoutes.POST("/types/delete/:id", server.deleteType)
	//拾取物品
	manRoutes.GET("/founds", server.listFound)
	manRoutes.GET("/founds/:id", server.getFound)
	manRoutes.POST("/founds/add", server.addFound)
	manRoutes.DELETE("/founds/delete/:id", server.deleteFound)
	//遗失物品
	manRoutes.GET("/losts", server.listLost)
	manRoutes.GET("/losts/:id", server.getLost)
	manRoutes.POST("/losts/add", server.addLost)
	manRoutes.DELETE("/losts/delete/:id", server.deleteLost)
	//归还物品
	manRoutes.GET("/matches", server.listMatch)
	manRoutes.GET("/matches/:id", server.getMatch)
	manRoutes.POST("/matches/add", server.addMatch)
	manRoutes.DELETE("/matches/delete/:id", server.deleteMatch)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
