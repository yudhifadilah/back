package models

type Note struct {
	Id        int64  `gorm:"primaryKey" json:"id"`
	Title     string `gorm:"type:varchar(250)" json:"title"`
	Deskripsi string `gorm:"type:text" json:"deskripsi"`
}
