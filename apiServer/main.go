package main

import (
	"faas/controllers"
	"faas/initializers"
	"faas/initializers/nclient"
	"faas/repository"
	"faas/services/auth"
	"faas/services/executor"
	"faas/services/functions"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	nclient.ConnectToNats()
	nclient.CreateJetStream()
	nclient.CreateUserKV()
	nclient.CreateFunctionKV()
	nclient.SubscribeFunctions()
	nclient.CreateJobStream()
}

func main() {
	router := gin.Default()
	conn := nclient.Client.Conn
	js := nclient.JS
	jwtKey := os.Getenv("SECRET")
	consumerKey := "faas_jwt_consumer"

	defer conn.Drain()

	functionRepo := repository.NewFunctionRepository()
	functionService, err := functions.NewService(*functionRepo)
	if err != nil {
		panic(err)
	}
	functionController := controllers.NewFunctionController(functionService)

	// Se inicializan los servicios de executor
	executorService := executor.NewExecutorService(conn, js)
	executorController := controllers.NewExecutorController(executorService)

	// Se inicializan los servicios de auth
	userRepo := repository.NewUserRepository()
	userController := controllers.NewUserController(userRepo)
	authService := auth.NewAuthService(jwtKey, consumerKey, *userRepo)
	authController := controllers.NewAuthController(authService)

	middleware := controllers.NewMiddleware(authService)

	routesHandler := controllers.NewRoutesHandler(authController, middleware, executorController, functionController, userController)

	routesHandler.SetRoutes(router)

	conn.Flush()

	router.Run(":8080")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
