package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Sql struct {
	Host     string
	Port     int
	UserName string
	PassWord string
	DbName   string
}

var DB *gorm.DB

func (s *Sql) Connect() error {
	dataSource := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.UserName, s.PassWord, s.DbName)
	db, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})

	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	DB = db
	fmt.Println("Connected to the database successfully")
	return nil
}

func (s *Sql) Close() {
	DB.Statement.ReflectValue.Close()
}
