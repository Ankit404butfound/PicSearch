package controllers

import (
	"PicSearch/app/api/services"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// file routes

func UploadFiles(c *gin.Context) {

	// userIdAny, _ := c.Get("userId")
	// userId := userIdAny.(int)
	// print(userId)
	form, err := c.MultipartForm()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save photos"})
		return
	}

	ok, err := services.UploadFiles(1, form.File["files"])
	fmt.Println("upload files result:", ok, err)
	if err != nil || !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Files uploaded successfully"})

}

func GetFiles(c *gin.Context) {

	// userIdAny, _ := c.Get("userId")
	userId := 1 //userIdAny.(int)

	query := c.Param("q")
	faceIdsStr := c.QueryArray("face_ids")

	// if query == "" && len(faceIdsStr) == 0 {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Either query or face is required"})
	// 	return
	// }
	var faceIds []int
	if len(faceIdsStr) > 0 {
		faceIds = make([]int, 0, len(faceIdsStr))

		for _, s := range faceIdsStr {
			id, err := strconv.Atoi(s)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("invalid file_id: %q (must be integer)", s),
				})
				return
			}
			faceIds = append(faceIds, id)
		}
	}

	files, err := services.GetFiles(userId, query, faceIds)
	// get all photos

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}
