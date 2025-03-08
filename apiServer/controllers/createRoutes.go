package controllers

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
	authController     *AuthController
	middleware         *Middleware
	executor           *ExecutorController
	functionController *FunctionController
	userController     *UserController
}

func NewRoutesHandler(
	authController *AuthController,
	middleware *Middleware,
	executor *ExecutorController,
	functionController *FunctionController,
	userController *UserController) *Routes {
	return &Routes{
		authController:     authController,
		middleware:         middleware,
		executor:           executor,
		functionController: functionController,
		userController:     userController,
	}
}

func (r *Routes) SetRoutes(ro *gin.Engine) {
	ro.POST("/register", r.authController.Register)
	ro.POST("/login", r.authController.Login)
	ro.GET("/validate", r.middleware.RequireAuth, r.authController.Validate)
	ro.POST("/execute", r.middleware.RequireAuth, r.executor.ExecuteFunction)
	ro.POST("/registerFunction", r.middleware.RequireAuth, r.functionController.RegisterFunction)
	ro.DELETE("/deleteFunction/:id", r.middleware.RequireAuth, r.functionController.DeleteFunction)
	ro.GET("/getFunctions", r.middleware.RequireAuth, r.functionController.GetFunctions)
	ro.GET("/users", r.middleware.RequireAuth, r.userController.GetAllUsers)
}
