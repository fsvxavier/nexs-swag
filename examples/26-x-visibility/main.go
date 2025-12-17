package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserPublic represents a public user
type UserPublic struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"John Doe"`
}

// UserPrivate represents internal user details
type UserPrivate struct {
	ID       int    `json:"id" example:"1"`
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john@example.com"`
	Password string `json:"password" example:"hashed_password"`
	Role     string `json:"role" example:"admin"`
}

// ErrorResponse represents an error
type ErrorResponse struct {
	Message string `json:"message" example:"Error message"`
}

// @title           X-Visibility Example API
// @version         1.0
// @description     API demonstrating @x-visibility public/private separation
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	r := gin.Default()

	r.GET("/api/v1/users/:id", GetUser)
	r.GET("/api/v1/admin/users/:id", GetUserAdmin)
	r.POST("/api/v1/users", CreateUser)

	r.Run(":8080")
}

// GetUser godoc
// @Summary      Get user (public)
// @Description  Get user details for public consumption
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  UserPublic
// @Failure      404  {object}  ErrorResponse
// @Router       /users/{id} [get]
// @x-visibility public
func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, UserPublic{ID: 1, Name: "John Doe"})
}

// GetUserAdmin godoc
// @Summary      Get user (admin)
// @Description  Get full user details including sensitive information
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  UserPrivate
// @Failure      404  {object}  ErrorResponse
// @Router       /admin/users/{id} [get]
// @x-visibility private
func GetUserAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, UserPrivate{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "hashed_password",
		Role:     "admin",
	})
}

// CreateUser godoc
// @Summary      Create user
// @Description  Create a new user (available in both public and private)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      UserPublic  true  "User data"
// @Success      201   {object}  UserPublic
// @Failure      400   {object}  ErrorResponse
// @Router       /users [post]
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, UserPublic{ID: 2, Name: "Jane Doe"})
}
