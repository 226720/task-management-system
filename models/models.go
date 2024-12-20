package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string
	Priority    string `gorm:"not null"`
	Deadline    string
	Status      string `gorm:"not null"`
	Category    string
	UserID      uint `gorm:"not null"`
}
