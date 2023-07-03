package jwt_middleware

import (
	"time"

	"github.com/yudhifadilah/back/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterUser godoc
// @Summary Register User.
// @Description Register User.
// @Tags Authentication
// @Accept application/json
// @Produce json
// @Success 200 {object} model.User
// Register handler
func RegisterUser(c *gin.Context) {
	var user models.User

	// Parsing body request ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Request Body Invalid",
		})
		return
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Gagal Register",
		})
		return
	}

	// Simpan user ke database
	user.Password = string(hashedPassword)
	if err := models.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{
			"message": "Gagal Register",
		})
		return
	}

	c.JSON(201, gin.H{
		"message": "Berhasil Register!",
	})
}

// LoginUser godoc
// @Summary Login User.
// @Description Login User.
// @Tags Authentication
// @Accept application/json
// @Produce json
// @Success 200 {object} model.User
// Login handler
func LoginUser(c *gin.Context) {
	var user models.User

	// Parsing body request ke struct User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Request Body Invalid",
		})
		return
	}

	inputPassword := user.Password

	// Cari user di database berdasarkan username
	result := models.DB.Where("username = ?", user.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{
				"message": "Terjadi Kesalahan",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "Failed to login",
		})
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
		c.JSON(401, gin.H{
			"message": "Username atau Password Salah",
		})
		return
	}

	// Buat JWT token
	claims := &models.JWTClaims{
		IdUser: user.IdUser,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret_key"))
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Gagal Mendapatkan Token",
		})
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})
}

// GetMe godoc
// @Summary Data user berdasarkan ID.
// @Description get data user.
// @Tags Authentication
// @Accept application/json
// @Produce json
// @Success 200 {object} model.User
// GetMe handler
func GetMe(c *gin.Context) {
	// Mendapatkan data user yang sedang login melalui JWT token
	user, _ := c.Get("user")
	claims := user.(*jwt.Token).Claims.(*models.JWTClaims)

	// Cari user di database berdasarkan user ID
	var userData models.User
	result := models.DB.First(&userData, claims.IdUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(401, gin.H{
				"message": "User Tidak Ditemukan",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "Gagal untuk Mendapatkan Data User",
		})
		return
	}

	c.JSON(200, gin.H{
		"user": userData,
	})
}

// Middleware untuk otentikasi JWT
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan token dari header Authorization
		authHeader := c.GetHeader("Authorization")
		token := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			token = authHeader[7:]
		} else {
			c.JSON(401, gin.H{
				"message": "Header Otorisasi Salah",
			})
			c.Abort()
			return
		}

		// Verifikasi token
		claims := new(models.JWTClaims)
		tkn, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret_key"), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(401, gin.H{
					"message": "Token Invalid",
				})
				c.Abort()
				return
			}
			c.JSON(401, gin.H{
				"message": "Gagal Mengautentikasi Token",
			})
			c.Abort()
			return
		}

		if !tkn.Valid {
			c.JSON(401, gin.H{
				"message": "Token Salah",
			})
			c.Abort()
			return
		}

		// Menyimpan data user ke local context
		c.Set("user", tkn)

		c.Next()
	}
}

func LogoutUser(c *gin.Context) {
	// Hapus token dari Authorization header
	c.Header("Authorization", "")

	c.JSON(200, gin.H{
		"message": "Logout berhasil",
	})
}
