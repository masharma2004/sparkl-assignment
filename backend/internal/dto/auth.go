package dto

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StudentSignupRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserResponse struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
}

type MeResponse struct {
	User UserResponse `json:"user"`
}
