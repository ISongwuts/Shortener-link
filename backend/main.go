package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ShortenerLink struct {
	gorm.Model
	OriginalURL string `gorm: "unique"`
	ShortURL string `gorm: "unique"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error .env file", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dbConcat := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"

	db, err := gorm.Open(mysql.Open(dbConcat), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&ShortenerLink{})

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/api/generate", func(c *gin.Context) {
		generatePostHandler(c, db)
	})

	router.GET("/:shortURL", func(c *gin.Context) {
		generateGetHandler(c, db)
	})

	router.Run(":8000")
}

func generateGetHandler (c *gin.Context, db *gorm.DB) {
	shortURL := c.Param("shortURL")
	var link ShortenerLink
	result := db.Where("short_url = ?", shortURL).Find(&link)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error" : "URL not found."})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.Redirect(http.StatusMovedPermanently, link.OriginalURL)
}

func generatePostHandler (c *gin.Context, db *gorm.DB) {
	var data struct {
		URL string `json: "url" binding: "required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
		return
	}

	var link ShortenerLink
	result := db.Where("original_url = ?", data.URL).First(&link)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			shortURL := generateShortURL()
			link = ShortenerLink{
				OriginalURL: data.URL,
				ShortURL: shortURL,
			}
			result := db.Create(&link)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error" : result.Error.Error()})
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"short_url": link.ShortURL})
}

func generateShortURL() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	var shortURL string
	for i := 0; i < 6; i++ {
		shortURL += string(chars[rand.Intn(len(chars))])
	}
	return shortURL
}