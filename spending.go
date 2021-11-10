package main

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Spending struct {
	ID       uint      `gorm:"primaryKey;not null;autoIncrement"`
	Date     time.Time `gorm:"not null"`
	Name     string    `gorm:"not null"`
	Amount   float32   `gorm:"not null"`
	Currency string    `gorm:"not null"`
}

type Spendings struct {
	db *gorm.DB
}

func New() *Spendings {
	return &Spendings{db: createDb()}
}

func (s *Spendings) Add(spending *Spending) error {
	result := s.db.Create(spending)

	if result.Error != nil {
		return result.Error
	}

	log.Println("Added spending")

	return nil
}

func (s *Spendings) Get(id uint) (*Spending, error) {
	var spending Spending
	result := s.db.First(&spending, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &spending, nil
}

func (s *Spendings) Delete(id uint) error {
	var spending Spending
	result := s.db.Delete(&spending, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *Spendings) Update(spending *Spending) error {
	result := s.db.Save(spending)

	if result.Error != nil {
		return result.Error
	}

	log.Println("Updated spending")

	return nil
}

func (s *Spendings) AllSpendings(year int, month time.Month) (*[]Spending, error) {
	var spendings []Spending

	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	result := s.db.Where("date >= ? and date < ?", firstOfMonth, firstOfNextMonth).Order("date ASC").Find(&spendings)

	if result.Error != nil {
		return nil, result.Error
	}

	return &spendings, nil
}

func createDb() *gorm.DB {
	dbUrl := os.Getenv("MYSQL_URL")
	db, err := gorm.Open(mysql.Open(dbUrl+"?parseTime=true"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&Spending{})

	return db
}
