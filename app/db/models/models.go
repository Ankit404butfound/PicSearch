package models

import "time"

type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time

	// Associations
	Files []File `gorm:"foreignKey:UserID"`
}

type File struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Description string
	UploadedAt  time.Time
	Embedding   []float32 `gorm:"type:vector(128)"`
	Url         string
	Size        float32
	UserID      int

	// Associations
	User  User   `gorm:"constraint:OnDelete:CASCADE;"`
	Faces []Face `gorm:"foreignKey:FileId"`
	Jobs  []Job  `gorm:"foreignKey:FileId"`
}

type Face struct {
	ID           int       `gorm:"primaryKey"`
	FileId       int       `gorm:"foreignKey"`
	UniqueFaceID int       `gorm:"foreignKey"`
	Coordinates  []float32 `gorm:"type:float[]"`

	// Associations
	UniqueFace UniqueFace `gorm:"constraint:OnDelete:CASCADE;"`
	File       File       `gorm:"constraint:OnDelete:CASCADE;"`
}

type UniqueFace struct {
	ID        int `gorm:"primary_key"`
	Name      string
	Embedding []float32 `gorm:"type:vector(128)"`

	// Associations
	Faces []Face `gorm:"foreignKey:UniqueFaceID"`
}

type Job struct {
	ID        int  `gorm:"primary_key"`
	FileId    int  `gorm:"foreignKey"`
	File      File `gorm:"constraint:OnDelete:CASCADE;"`
	Status    string
	StartedAt time.Time
	EndedAt   time.Time
}

type Devices struct {
	ID        int `gorm:"primaryKey"`
	UserID    int `gorm:"foreignKey"`
	DeviceID  string
	CreatedAt time.Time

	// Associations
	User User `gorm:"constraint:OnDelete:CASCADE;"`
}
