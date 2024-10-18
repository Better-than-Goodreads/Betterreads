package models

type UserRecord struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Location  string `json:"location"`
	Gender    string `json:"gender"`
	AboutMe   string `json:"about_me"`
	Id        int    `json:"id"`
	Age       int    `json:"age"`
}

type UserResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Location  string `json:"location"`
	Gender    string `json:"gender"`
	AboutMe   string `json:"about_me"`
	Id        int    `json:"id"`
	Age       int    `json:"age"`
}

type UserLoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type UserRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Location  string `json:"location"`
	Gender    string `json:"gender"`
	AboutMe   string `json:"about_me"`
	Age       int    `json:"age"`
}
