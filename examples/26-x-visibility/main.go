package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserPublic represents a public user.
type UserPublic struct {
	ID   int    `example:"1"        json:"id"`
	Name string `example:"John Doe" json:"name"`
}

// UserPrivate represents internal user details.
type UserPrivate struct {
	ID       int    `example:"1"                json:"id"`
	Name     string `example:"John Doe"         json:"name"`
	Email    string `example:"john@example.com" json:"email"`
	Password string `example:"hashed_password"  json:"password"`
	Role     string `example:"admin"            json:"role"`
}

// ErrorResponse represents an error.
type ErrorResponse struct {
	Message string `example:"Error message" json:"message"`
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
// @x-visibility public.
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
// @x-visibility private.
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
// @Router       /users [post].
func CreateUser(c *gin.Context) {
	c.JSON(http.StatusCreated, UserPublic{ID: 2, Name: "Jane Doe"})
}
