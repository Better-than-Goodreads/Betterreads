package models

type UserRecord struct {
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
	Location  string `json:"location" db:"location"`
	Gender    string `json:"gender" db:"gender"`
	AboutMe   string `json:"about_me" db:"about_me"`
	Id        string    `json:"id" db:"id"`
	Age       int    `json:"age" db:"age"`
}

type UserStageRecord struct {
    Id string `json:"id" db:"id"`
    Email string `json:"email" db:"email"`
    Username string `json:"username" db:"username"`
    Password string `json:"password" db:"password"`
    FirstName string `json:"first_name" db:"first_name"`
    LastName string `json:"last_name" db:"last_name"`
}

type UserStageResponse struct {
    Email string `json:"email"`
    Username string `json:"username"`
    First_name string `json:"first_name"`
    Last_name string `json:"last_name"`
    Id string `json:"id"`
}

type UserStageRequest struct {
    Email    string `json:"email" binding:"required" db:"email"`
    Username string `json:"username" binding:"required" db:"username"`
    Password string `json:"password" binding:"required" db:"password"`
    FirstName string `json:"first_name" binding:"required" db:"first_name"`
    LastName string `json:"last_name" binding:"required" db:"last_name"`
}

type UserAdditionalRequest struct {
    Username string `json:"username" binding:"required"`
    Location string `json:"location"`
    Gender string `json:"gender"`
    AboutMe string `json:"about_me"`
    Id string `json:"id"`
    Age int `json:"age"`
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
