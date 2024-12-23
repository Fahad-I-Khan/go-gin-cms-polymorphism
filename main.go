package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "CMS-Polymorphism/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type ErrorResponse struct {
	Message string `json:"message"`
}

type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime"` // In Postman it failed to set value and because of this ID, Taggable, Commentable also failed when it was set to string type.
	UpdatedAt time.Time `gorm:"autoUpdateTime"` // In Postman it failed to set value and because of this ID, Taggable, Commentable also failed when it was set to string type.
}

// Article model
type Article struct {
	BaseModel
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Tags     []Tag     `gorm:"polymorphic:Taggable;"`
	Comments []Comment `gorm:"polymorphic:Commentable;"`
}

// Video model
type Video struct {
	BaseModel
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Tags     []Tag     `gorm:"polymorphic:Taggable;"`
	Comments []Comment `gorm:"polymorphic:Commentable;"`
}

// Tag model
type Tag struct {
	BaseModel
	Name         string `json:"name"`
	TaggableID   uint   `gorm:"index"`
	TaggableType string `gorm:"index"`
}

// Comment model
type Comment struct {
	BaseModel
	Content         string `json:"content"`
	CommentableID   uint   `gorm:"index"`
	CommentableType string `gorm:"index"`
}

// Initialize DB connection
func initDB() {

	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database", err)
	}

	// Auto-migrate models
	db.AutoMigrate(&Article{}, &Video{}, &Tag{}, &Comment{})
}

// @title CMS API
// @version 1.0
// @description API for managing articles, videos, tags, and comments.
// @host localhost:8000
// BasePath /api/v1
// @contact.name API Support
// @contact.url http://localhost:8000/support   // Local URL for your development environment
// @contact.email support@localhost.com
func main() {
	// Initialize the DB
	initDB()

	r := gin.Default()
	r.Use(cors.Default())
	// Serve Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Define other routes here...
	r.GET("/api/v1/articles", getarticles)
	r.POST("/api/v1/articles", createarticles)
	r.PUT("/api/v1/articles/:id", updatearticles)
	r.GET("/api/v1/videos", getvideos)
	r.POST("/api/v1/videos", createvideos)
	r.PUT("/api/v1/videos/:id", updatevideos)

	// Start the server
	if err := r.Run(":8000"); err != nil {
		log.Fatal("Failed to start the server:", err)
	}
}

// @Summary List all articles
// @Description Retrieve a list of all articles, including their tags and comments.
// @Tags Articles
// @Produce json
// @Success 200 {array} Article
// @Router /api/v1/articles [get]
func getarticles(c *gin.Context)  {
	var articles []Article
	if err := db.Preload("Tags").Preload("Comments").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error fetching articles"})
		return
	}
	c.JSON(200, articles)
}

// @Summary Create a new article
// @Description Add a new article to the database.
// @Tags Articles
// @Accept json
// @Produce json
// @Param article body Article true "Article Data"
// @Success 201 {object} Article
// @Router /api/v1/articles [post]
func createarticles(c *gin.Context)  {
	var article Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
		return
	}
	db.Create(&article)
	c.JSON(201, article)
}

// @Summary Update an article
// @Description Update the details of an article by its ID.
// @Tags Articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param article body Article true "Article Data"
// @Success 200 {object} Article
// @Router /api/v1/articles/{id} [put]
func updatearticles(c *gin.Context)  {
	id := c.Param("id")
	var article Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Article not found"})
		return
	}

	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
		return
	}

	db.Save(&article)
	c.JSON(200, article)
		
}

// @Summary List all videos
// @Description Retrieve a list of all videos, including their tags and comments.
// @Tags Videos
// @Produce json
// @Success 200 {array} Video
// @Router /api/v1/videos [get]
func getvideos(c *gin.Context)  {
	var videos []Video
	if err := db.Preload("Tags").Preload("Comments").Find(&videos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error fetching users"})
		return
	}
	c.JSON(200, videos)
}

// @Summary Create a new video
// @Description Add a new video to the database.
// @Tags Videos
// @Accept json
// @Produce json
// @Param video body Video true "Video Data"
// @Success 201 {object} Video
// @Router /api/v1/videos [post]
func createvideos(c *gin.Context)  {
	var video Video
	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
		return
	}
	db.Create(&video)
	c.JSON(201, video)
}

// @Summary Update a video
// @Description Update the details of a video by its ID.
// @Tags Videos
// @Accept json
// @Produce json
// @Param id path int true "Video ID"
// @Param video body Video true "Video Data"
// @Success 200 {object} Video
// @Router /api/v1/videos/{id} [put]
func updatevideos(c *gin.Context)  {
	id := c.Param("id")
	var video Video
	if err := db.First(&video, id).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Message: "Video not found"})
		return
	}

	if err := c.ShouldBindJSON(&video); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
		return
	}

	if err := db.Save(&video).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update video"})
		return
	}

	c.JSON(200, video)
}

// func setupRouter(db *gorm.DB) *gin.Engine {
// 	r := gin.Default()
// 	r.Use(cors.Default())

// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// Article Routes
// 	articleGroup := r.Group("/articles")
// 	{
// 		// @Summary List Articles
// 		// @Produce json
// 		// @Success 200 {array} Article
// 		// @Router /articles [get]
// 		articleGroup.GET("", func(c *gin.Context) {
// 			var articles []Article
// 			if err := db.Preload("Tags").Preload("Comments").Find(&articles).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error fetching articles"})
// 				return
// 			}
// 			c.JSON(200, articles)
// 		})

// 		// @Summary Create Article
// 		// @Produce json
// 		// @Param article body Article true "Article Data"
// 		// @Success 201 {object} Article
// 		// @Router /articles [post]
// 		articleGroup.POST("", func(c *gin.Context) {
// 			var article Article
// 			if err := c.ShouldBindJSON(&article); err != nil {
// 				c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
// 				return
// 			}
// 			db.Create(&article)
// 			c.JSON(201, article)
// 		})

// 		// @Summary Update Article
// 		// @Produce json
// 		// @Param id path int true "Article ID"
// 		// @Param article body Article true "Article Data"
// 		// @Success 200 {object} Article
// 		// @Router /articles/{id} [put]
// 		articleGroup.PUT("/:id", func(c *gin.Context) {
// 			id := c.Param("id")
// 			var article Article
// 			if err := db.First(&article, id).Error; err != nil {
// 				c.JSON(http.StatusNotFound, ErrorResponse{Message: "Article not found"})
// 				return
// 			}

// 			if err := c.ShouldBindJSON(&article); err != nil {
// 				c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
// 				return
// 			}

// 			db.Save(&article)
// 			c.JSON(200, article)
// 		})
// 	}

// 	// Video Routes
// 	videoGroup := r.Group("/videos")
// 	{
// 		// @Summary List Videos
// 		// @Produce json
// 		// @Success 200 {array} Video
// 		// @Router /videos [get]
// 		videoGroup.GET("", func(c *gin.Context) {
// 			var videos []Video
// 			if err := db.Preload("Tags").Preload("Comments").Find(&videos).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Error fetching users"})
// 				return
// 			}
// 			c.JSON(200, videos)
// 		})

// 		// @Summary Create Video
// 		// @Produce json
// 		// @Param video body Video true "Video Data"
// 		// @Success 201 {object} Video
// 		// @Router /videos [post]
// 		videoGroup.POST("", func(c *gin.Context) {
// 			var video Video
// 			if err := c.ShouldBindJSON(&video); err != nil {
// 				c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
// 				return
// 			}
// 			db.Create(&video)
// 			c.JSON(201, video)
// 		})

// 		// @Summary Update Video
// 		// @Produce json
// 		// @Param id path int true "Video ID"
// 		// @Param video body Video true "Video Data"
// 		// @Success 200 {object} Video
// 		// @Router /videos/{id} [put]
// 		videoGroup.PUT("/:id", func(c *gin.Context) {
// 			id := c.Param("id")
// 			var video Video
// 			if err := db.First(&video, id).Error; err != nil {
// 				c.JSON(http.StatusNotFound, ErrorResponse{Message: "Video not found"})
// 				return
// 			}

// 			if err := c.ShouldBindJSON(&video); err != nil {
// 				c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid input"})
// 				return
// 			}

// 			if err := db.Save(&video).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Failed to update video"})
// 				return
// 			}

// 			c.JSON(200, video)
// 		})
// 	}

// 	return r
// }