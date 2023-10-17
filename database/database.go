package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDB(dns string) (db *gorm.DB, err error) {
	db, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

