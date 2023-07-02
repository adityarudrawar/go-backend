package models

type User struct {
	Username string `json:"username"`
	Password []byte `json:"-"`
}

type UserSessions struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func (UserSessions) TableName() string {
	return "usersessions"
}