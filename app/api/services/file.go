package services

import (
	"PicSearch/app/api/utils"
	"PicSearch/app/db"
	"PicSearch/app/db/models"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

func UploadFiles(userId int, files []*multipart.FileHeader) (bool, error) {

	for _, file := range files {

		base := strings.Split(filepath.Base(file.Filename), ".")[0]
		ext := filepath.Ext(file.Filename)
		//fileSaveName := base + time.Now().Format("20060102150405") + ext
		fileSaveName := fmt.Sprintf("%s_%d%s", base, time.Now().UnixNano(), ext)
		savePath := filepath.Join("uploads", fileSaveName)

		dst, err := os.Create(savePath)
		if err != nil {
			return false, fmt.Errorf("failed to create file: %w", err)
		}
		defer dst.Close()

		var uploadFile models.File
		uploadFile.Url = savePath
		uploadFile.Name = fileSaveName
		uploadFile.Size = float32(file.Size)
		uploadFile.UserId = userId

		user_err := db.DB.Create(&uploadFile).Error

		if user_err != nil {
			return false, nil
		}

	}

	return true, nil

}

func GetFiles(id int, query string, faceIds []int) ([]models.File, error) {
	var files []models.File

	// Initialize base query with JOIN
	dbQuery := db.DB.Model(&models.Face{}).
		Select("files.*").
		Joins("JOIN files ON faces.file_id = files.id")

	// Add WHERE condition if faceIds is provided
	if len(faceIds) > 0 {
		dbQuery = dbQuery.Where("faces.unique_face_id IN ?", faceIds)
	}

	queryEmbedding, err := utils.GetEmbeddings(query)
	if err != nil {
		return nil, err
	}

	// Add ORDER BY clause if query (embedding) is provided
	if query != "" {
		// Assuming you have pgvector extension and the embedding is a vector type
		// Use raw SQL for the vector distance operator
		dbQuery = dbQuery.Order(gorm.Expr("files.embedding <=> ?", queryEmbedding))
	}

	// Add LIMIT and execute
	err = dbQuery.Limit(10).Find(&files).Error
	if err != nil {
		return nil, err
	}

	return files, nil
}

// func GetFiles(id int, querry *string, faced_id *[]int) ([]models.File, error) {

// 	var files []models.File
// 	var uniqueFace models.UniqueFace

// 	if querry!=nil && faced_id != nil {

// 		// Example: simple contains search

// 		db.DB.Find(&files).Order()
// 		// querryEmbedding := getEmbeddings(querry)
// 		// err := db.DB.Find(files, querryEmbedding).Error

// 		// if err != nil {
// 		// 	return nil, err
// 		// }

// 		// return files, nil

// 		subQuery := db.DB.
// 		Select("DISTINCT file_id").
// 		Table("faces").
// 		Where("unique_face_id IN ?", faced_id)

// 	// Main query: get files with embedding similarity and face filter
// 		err := db.DB.
// 			Where("id IN (?)", subQuery).
// 			Order("embedding <=> ?", querryEmbedding).
// 			Limit(limit).
// 			Offset(offset).
// 			Find(&files).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Eager load associations if needed
// 	err = r.DB.Preload("Faces").Preload("Faces.UniqueFace").Find(&files).Error

// 	return files, err

// 	}

// 	if faced_id != nil {
// 		err := db.DB.Find(uniqueFace, faced_id).Error

// 		if err != nil {
// 			file_err := db.DB.Find(files, uniqueFace.Embedding).Error

// 			if file_err != nil {
// 				return files, nil
// 			}
// 		}
// 	}

// 	return files, nil

// }
