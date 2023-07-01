package models

type User struct {
	Id       uint   `json:"id" gorm:"autoIncrement;not null"`
	Username string `json:"username" gorm:"primaryKey,unique"`
	Password []byte `json:"-"`
}

type Session struct {
	Username string `json:"username" gorm:"primaryKey,unique"`
	Jwt      string `json:"jwt"`
	Expires  int64  `json:"expires_at"`
}