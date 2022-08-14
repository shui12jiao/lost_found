package api

import (
	"fmt"
	"log"
	"lost_found/middleware"
	"lost_found/middleware/session"
	"lost_found/middleware/token"
	"lost_found/service"
	"lost_found/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/silenceper/wechat/v2/miniprogram"
)

const (
	cookieName = middleware.CookieName
	appid      = "appid"
	secret     = "secret"
	domain     = "localhost:8080"
)

type Server struct {
	config         util.Config
	store          service.Store
	tokenMaker     token.Maker
	sessionManager *session.Manager
	wxMini         *miniprogram.MiniProgram
	router         *gin.Engine
}

func NewServer(config util.Config, store service.Store) (*Server, error) {
	//令牌maker
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}
	//session会话管理
	sessionManager := session.NewSessionManager(cookieName, session.NewMemorySessionStore(), config.SessionLifeTime)

	server := &Server{
		config:         config,
		store:          store,
		tokenMaker:     tokenMaker,
		sessionManager: sessionManager,
		wxMini:         service.GetWechatMiniProgram(),
		router:         gin.Default(),
	}
	//设置验证器
	if _, ok := binding.Validator.Engine().(*validator.Validate); !ok {
		log.Fatal("failed to init validator")
	}
	//开始路由
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := server.router

	router.POST("/users/login", server.loginUser)

	//--普通用户使用功能--
	authRoutes := router.Group("/").Use(middleware.AuthSessionMiddleware(server.sessionManager))
	//用户账户
	authRoutes.GET("/users/info", server.getUserSelf)
	authRoutes.DELETE("/users", server.deleteUserSelf)
	authRoutes.PATCH("/users", server.updateUserSelf)
	//用户物品
	authRoutes.GET("/users/found", server.listMyFound)
	authRoutes.GET("/users/lost", server.listMyLost)
	authRoutes.GET("/users/match", server.listMyMatch)
	authRoutes.POST("/users/found/:id", server.completeMyFound) //自己发布的拾取物已归还给失主
	authRoutes.POST("/users/lost/:id", server.completeMyLost)   //自己发布的遗失物已被寻回
	//拾取物品
	authRoutes.GET("/founds", server.listFound)
	authRoutes.GET("/founds/:id", server.getFound)
	authRoutes.POST("/founds/add", server.addFound)
	authRoutes.DELETE("/founds/:id", server.deleteFound) //只能删除自己发布的拾取物
	//遗失物品
	authRoutes.GET("/losts", server.listLost)
	authRoutes.GET("/losts/:id", server.getLost)
	authRoutes.POST("/losts/add", server.addLost)
	authRoutes.DELETE("/losts/:id", server.deleteLost) //只能删除自己的发布的遗失物
	//归还物品
	authRoutes.GET("/matches", server.listMatch)
	//位置和类型
	authRoutes.GET("/locations", server.listLocation)
	authRoutes.GET("/types", server.listType)
	//管理员
	authRoutes.GET("/managers", server.listManager)

	//--管理员系统--
	manRoutes := router.Group("/admin").Use(middleware.ManagerMiddleware(server.sessionManager, server.store))
	//用户账户
	manRoutes.GET("/users", server.listUser)
	manRoutes.POST("/users/add", server.addUser)
	manRoutes.GET("/users/:id", server.getUser)
	manRoutes.GET("/users/some", server.searchUser)
	manRoutes.DELETE("/users/:id", server.deleteUser)
	manRoutes.PATCH("/users", server.updateUser)
	//位置和类型
	manRoutes.GET("/locations", server.listLocation)
	manRoutes.POST("/locations/wide/add", server.addLocationWide)
	manRoutes.DELETE("/locations/wide/delete/:id", server.deleteLocationWide)
	manRoutes.POST("/locations/narrow/add", server.addLocationNarrow)
	manRoutes.DELETE("/locations/narrow/delete/:id", server.deleteLocationNarrow)
	manRoutes.GET("/types", server.listType)
	manRoutes.POST("/types/wide/add", server.addTypeWide)
	manRoutes.POST("/types/wide/delete/:id", server.deleteTypeWide)
	manRoutes.POST("/types/wide/narrow/add", server.addTypeNarrow)
	manRoutes.POST("/types/wide/narrow/delete/:id", server.deleteTypeNarrow)
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
	//管理员账户
	manRoutes.GET("/managers", server.listManager)
	manRoutes.POST("/managers/add", server.addManager)
	manRoutes.GET("/managers/:id", server.getManager)
	manRoutes.DELETE("/managers/:id", server.deleteManager)
	manRoutes.PATCH("/managers", server.updateManager)
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
