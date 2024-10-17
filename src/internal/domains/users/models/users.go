package models

type UserRecord struct {
	Id			int    `json:"id"`
	Email		string `json:"email"`
	Password	string `json:"password"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Username	string `json:"username"`
	Location	string `json:"location"`
	Gender		string `json:"gender"`
	Age			int    `json:"age"`
	AboutMe		string `json:"about_me"`
}

type UserResponse struct {
	Id			int    `json:"id"`
	Email		string `json:"email"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Username	string `json:"username"`
	Location	string `json:"location"`
	Gender		string `json:"gender"`
	Age			int    `json:"age"`
	AboutMe		string `json:"about_me"`
}

type UserRequest struct {
	Email		string `json:"email"`
	Password	string `json:"password"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	Username	string `json:"username"`
	Location	string `json:"location"`
	Gender		string `json:"gender"`
	Age			int    `json:"age"`
	AboutMe		string `json:"about_me"`
}


