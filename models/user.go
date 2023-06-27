package models

type User struct {
	Id       uint   `json:"id" gorm:"autoIncrement;not null"`
	Username string `json:"username" gorm:"primaryKey,unique"`
	Password []byte `json:"-"`
}