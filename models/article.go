package models

import "gorm.io/gorm"

type Articles struct {
	ID        uint    `gorm:"primaryKey; autoIncrement" json:"id"`
	Author    *string `json:"author"`
	Title     *string `json:"title"`
	Publisher *string `json:"publisher"`
}

func MigrateArticles(db *gorm.DB) error {
	err := db.AutoMigrate(&Articles{})

	return err
}
