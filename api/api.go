package api

import (
	v1 "blog_project/api/v1"
	"blog_project/config"
	"blog_project/storage"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	_ "blog_project/api/docs"
)

type RouterOptions struct {
	Cfg      *config.Config
	Storage  storage.StorageI
	InMemory storage.InMemoryStorageI
}

// @title           Swagger for blog api
// @version         1.0
// @description     This is a blog service api.
// @BasePath  /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @Security ApiKeyAuth
func New(opt *RouterOptions) *gin.Engine {
	router := gin.Default()

	handlerV1 := v1.New(&v1.HabdlerV1Options{
		Cfg:      opt.Cfg,
		Storage:  opt.Storage,
		InMemory: opt.InMemory,
	})

	router.Static("/media", "./media")

	apiV1 := router.Group("/v1")

	apiV1.GET("/users/:id", handlerV1.GetUser)
	apiV1.POST("/users", handlerV1.CreateUser)
	apiV1.GET("/users", handlerV1.GetAllUsers)
	// apiV1.PUT("/user/:id", handlerV1.UpdateUser)
	apiV1.DELETE("/user/:id", handlerV1.DeleteUser)

	apiV1.GET("/categories/:id", handlerV1.GetCategory)
	apiV1.POST("/categories", handlerV1.CreateCategory)
	apiV1.GET("/categories", handlerV1.GetAllCategories)
	// apiV1.PUT("/category/:id", handlerV1.UpdateCategory)
	// apiV1.DELETE("/category/:id", handlerV1.DeleteCategory)

	apiV1.GET("/posts/:id", handlerV1.GetPost)
	apiV1.POST("/posts", handlerV1.CreatePost)
	apiV1.GET("/posts", handlerV1.GetAllPosts)
	// apiV1.PUT("/post/:id", handlerV1.UpdatePost)
	// apiV1.DELETE("/post/:id", handlerV1.DeletePost)

	apiV1.POST("/comments", handlerV1.AuthMiddleWare, handlerV1.CreateComment)
	apiV1.GET("/comments", handlerV1.GetAllComments)
	apiV1.GET("/comment/:id", handlerV1.GetComment)
	apiV1.DELETE("/comment/:id", handlerV1.DeleteComment)

	apiV1.POST("/auth/register", handlerV1.Register)
	apiV1.POST("/auth/login", handlerV1.Login)
	apiV1.POST("/auth/verify", handlerV1.Verify)
	apiV1.POST("/auth/forgot-password", handlerV1.ForgotPassword)
	apiV1.POST("/auth/verify-forgot-password", handlerV1.VerifyForgotPassword)
	apiV1.POST("/auth/update-password", handlerV1.AuthMiddleWare, handlerV1.UpdatePassword)

	apiV1.POST("/file-upload", handlerV1.UploadFile)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
