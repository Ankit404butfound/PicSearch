package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/datatypes"
)

type User struct {
	ID        int       `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	// Associations
	Files []File `gorm:"foreignKey:UserID"`
}

type File struct {
	ID         int             `gorm:"primaryKey"`
	Name       string          `gorm:"type:varchar(255);not null"`
	Metadata   datatypes.JSON  `gorm:"type:json"`
	UploadedAt time.Time       `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	Embedding  pgvector.Vector `pg:"type:vector(512)"`
	Url        string          `gorm:"type:varchar(512);not null"`
	Size       float32         `gorm:"type:float"`
	UserId     int             `gorm:"foreignKey;not null"`

	// Associations
	User  User   `gorm:"constraint:OnDelete:CASCADE;"`
	Faces []Face `gorm:"foreignKey:FileId"`
	Jobs  []Job  `gorm:"foreignKey:FileId"`
}

type Face struct {
	ID           int       `gorm:"primaryKey"`
	FileId       int       `gorm:"foreignKey;not null"`
	UniqueFaceID int       `gorm:"foreignKey;not null"`
	Coordinates  []float32 `gorm:"type:float[];not null"`

	// Associations
	UniqueFace UniqueFace `gorm:"constraint:OnDelete:CASCADE;"`
	File       File       `gorm:"constraint:OnDelete:CASCADE;"`
}

type UniqueFace struct {
	ID        int `gorm:"primary_key"`
	Name      string
	Embedding []float32 `gorm:"type:vector(128);not null"`

	// Associations
	Faces []Face `gorm:"foreignKey:UniqueFaceID"`
}

type Job struct {
	ID        int       `gorm:"primary_key"`
	FileId    int       `gorm:"foreignKey;not null"`
	File      File      `gorm:"constraint:OnDelete:CASCADE;"`
	Status    string    `gorm:"type:varchar(50);not null"`
	StartedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	EndedAt   time.Time `gorm:"type:timestamp"`
}

type Devices struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `gorm:"foreignKey;not null"`
	DeviceID  string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	// Associations
	User User `gorm:"constraint:OnDelete:CASCADE;"`
}
