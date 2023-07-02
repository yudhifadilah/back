package notecontroller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yudhifadilah/back/models"
	"gorm.io/gorm"
)

func Index(c *gin.Context) {

	var notes []models.Note

	models.DB.Find(&notes)
	c.JSON(http.StatusOK, gin.H{"notes": notes})

}

func Show(c *gin.Context) {

	var notes models.Note
	id := c.Param("id")

	if err := models.DB.First(&notes, id).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Data Tidak Ditemukan"})
			return

		default:
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})

}

func Create(c *gin.Context) {

	var notes models.Note

	if err := c.ShouldBindJSON(&notes); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	models.DB.Create(&notes)
	c.JSON(http.StatusOK, gin.H{"notes": notes})

}

func Update(c *gin.Context) {
	var notes models.Note
	id := c.Param("id")

	if err := c.ShouldBindJSON(&notes); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if models.DB.Model(&notes).Where("id = ?", id).Updates(&notes).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Tidak Dapat MengUpdate Data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil diUpdate"})

}

func Delete(c *gin.Context) {
	Id := c.Param("id")
	//var notes models.Note

	// Parse the ID into an integer
	id, err := strconv.ParseUint(Id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Delete the user by ID
	if err := models.DB.Delete(&models.Note{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data deleted successfully"})
}
