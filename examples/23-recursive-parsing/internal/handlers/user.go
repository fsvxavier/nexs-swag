package handlers

// InternalUser represents internal user with sensitive data
type InternalUser struct {
	ID       int    `json:"id" example:"1"`
	Email    string `json:"email" example:"user@example.com"`
	Password string `json:"password" example:"hashed_password"`
	Role     string `json:"role" example:"admin"`
}

// CreateUser creates a new user
// @Summary Create user (internal handler)
// @Description Creates a new user with internal validation
// @Tags users
// @Accept json
// @Produce json
// @Param user body InternalUser true "User data"
// @Success 201 {object} InternalUser
// @Failure 400 {string} string "Bad request"
// @Router /internal/users [post]
func CreateUser() {}

// UpdateUser updates an existing user
// @Summary Update user (internal handler)
// @Description Updates user information with internal privileges
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body InternalUser true "User data"
// @Success 200 {object} InternalUser
// @Failure 404 {string} string "User not found"
// @Router /internal/users/{id} [put]
func UpdateUser() {}
