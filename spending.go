package main

import (
	"log"
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

func (s *Spendings) AllSpendings() (*[]Spending, error) {
	var spendings []Spending

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, now.Location())
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)

	result := s.db.Where("date >= ? and date < ?", firstOfMonth, firstOfNextMonth).Find(&spendings)

	if result.Error != nil {
		return nil, result.Error
	}

	return &spendings, nil
}

func (s *Spendings) TotalAmount() (float32, error) {
	spendings, err := s.AllSpendings()

	if err != nil {
		return 0.0, err
	}

	var sum float32
	for _, spending := range *spendings {
		sum += spending.Amount
	}

	return sum, nil
}

func createDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:test@/spendings?parseTime=true"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&Spending{})

	return db
}
