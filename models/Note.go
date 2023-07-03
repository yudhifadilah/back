package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Note struct {
	Id        int64  `gorm:"primaryKey" json:"id"`
	Title     string `gorm:"type:varchar(250)" json:"title"`
	Deskripsi string `gorm:"type:text" json:"deskripsi"`
}

type User struct {
	IdUser    uint   `gorm:"primaryKey;column:id_user;autoIncrement" json:"-"`
	Nama      string `gorm:"column:nama" json:"nama"`
	Username  string `gorm:"uniqueIndex" json:"username"`
	Password  string `json:"-"`
	CreatedAt []uint8
	UpdatedAt []uint8
	DeletedAt *time.Time
}

// jwt struct
type JWTClaims struct {
	jwt.StandardClaims
	IdUser uint `json:"id_user"`
}
