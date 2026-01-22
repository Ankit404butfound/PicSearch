package services

import (
	"PicSearch/app/api/schemas"
	"PicSearch/app/api/utils"
	"PicSearch/app/db"
	"PicSearch/app/db/models"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pgvector/pgvector-go"
)

func UploadFile(userId int, file *multipart.FileHeader) (bool, error) {

	base := strings.Split(filepath.Base(file.Filename), ".")[0]
	ext := filepath.Ext(file.Filename)
	//fileSaveName := base + time.Now().Format("20060102150405") + ext
	fileSaveName := fmt.Sprintf("%s_%d%s", base, time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", fileSaveName)

	src, err := file.Open()
	if err != nil {
		fmt.Println("error opening file:", err)
		return false, err
	}
	defer src.Close()

	out, err := os.Create(savePath)
	if err != nil {
		fmt.Println("error creating file:", err)
		return false, err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		fmt.Println("error saving file:", err)
		return false, err
	}

	var uploadFile models.File
	uploadFile.Url = fmt.Sprintf("%s/api/files/download/%s", os.Getenv("SERVER_HOST"), fileSaveName)
	uploadFile.Name = fileSaveName
	uploadFile.Size = float32(file.Size)
	uploadFile.UserId = userId
	uploadFile.Embedding = nil

	user_err := db.DB.Create(&uploadFile).Error

	if user_err != nil {
		fmt.Println("error saving file record to db:", user_err)
		return false, nil
	}

	utils.TriggerImageProcessingJob(uploadFile.ID)

	return true, nil

}

// func GetFiles(id int, query string, faceIds []int) ([]models.File, error) {
// 	var files []models.File

// 	// Check if faceIds is empty (if there's no `id` check or something related)
// 	// If id is relevant to the query, ensure you're passing it correctly into the function
// 	// Initialize base query with JOIN
// 	dbQuery := db.DB.Model(&models.File{})

// 	// Add JOIN condition to include faces if faceIds are provided
// 	if len(faceIds) > 0 {
// 		// Correcting the join to use the correct relation between `faces` and `files`
// 		dbQuery = dbQuery.Joins("JOIN faces ON faces.file_id = files.id").Where("faces.unique_face_id IN ?", faceIds)
// 	}

// 	// Add ORDER BY clause for vector similarity calculation
// 	if query != "" {
// 		// Calculate the query embedding here
// 		queryEmbedding, err := utils.GetEmbeddings(query)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get embeddings: %v", err)
// 		}
// 		// Use raw SQL for the vector distance operator with ORDER BY
// 		dbQuery = dbQuery.Clauses(clause.OrderBy{
// 			Expression: clause.Expr{
// 				SQL: "embedding <=> ?", Vars: []interface{}{pgvector.NewVector(queryEmbedding)},
// 			},
// 		})
// 	}

// 	// Execute the query with LIMIT and retrieve the files
// 	err := dbQuery.Limit(10).Find(&files).Error
// 	if err != nil {
// 		return nil, fmt.Errorf("error retrieving files: %v", err)
// 	}

// 	return files, nil
// }

func GetFiles(userId int, query string, faceIds []int) ([]schemas.FileResponse, error) {
	var files []schemas.FileResponse
	var queryEmbedding []float32
	var err error

	if query != "" && len(faceIds) > 0 {
		queryEmbedding, err = utils.GetEmbeddings(query)
		if err != nil {
			return nil, fmt.Errorf("failed to get embeddings: %v", err)
		}
		err = db.DB.Raw(`
			SELECT
				embedding <=> ? AS distance, files.id, name, metadata, url, size, uploaded_at
			FROM
				files
			JOIN
				faces ON faces.file_id = files.id
			WHERE
				faces.unique_face_id IN ?
			ORDER BY
				embedding <=> ?
			LIMIT ? OFFSET ?;
		`, pgvector.NewVector(queryEmbedding), faceIds, pgvector.NewVector(queryEmbedding), 10, 0).Scan(&files).Error

	} else if query != "" {
		var err error
		queryEmbedding, err = utils.GetEmbeddings(query)
		if err != nil {
			return nil, fmt.Errorf("failed to get embeddings: %v", err)
		}
		err = db.DB.Raw(`
			SELECT
				embedding <=> ? AS distance, id, name, metadata, url, size, uploaded_at
			FROM
				files
			ORDER BY
				embedding <=> ?
			LIMIT ?;
		`, pgvector.NewVector(queryEmbedding), pgvector.NewVector(queryEmbedding), 10).Scan(&files).Error

	} else if len(faceIds) > 0 {
		err = db.DB.Raw(`
			SELECT
				files.id, name, metadata, url, size, uploaded_at
			FROM
				files
			JOIN
				faces ON faces.file_id = files.id
			WHERE
				faces.unique_face_id IN ?
			LIMIT ? OFFSET ?;
		`, faceIds, 10, 0).Scan(&files).Error

	} else {
		return nil, fmt.Errorf("either query or faceIds must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("error retrieving files: %v", err)
	}
	fmt.Println("files retrieved:", files)
	return files, nil
}
