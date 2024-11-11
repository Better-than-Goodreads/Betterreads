package models


type FriendOfUser struct {
    ID       string `json:"id" db:"id"`
    Username string `json:"username" db:"username"`
}


