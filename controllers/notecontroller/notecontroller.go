package notecontroller

import (
	"encoding/json"
	"net/http"

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

	var notes models.Note

	var input struct {
		Id json.Number
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, _ := input.Id.Int64()
	if models.DB.Delete(&notes, id).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Gagal Menghapus Data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data Berhasil di Hapus"})

}
