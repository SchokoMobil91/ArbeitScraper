package main

import (
	"arbeit-scraper/database"
	"arbeit-scraper/scraper"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/jobs", func(c *gin.Context) {
		log.Println("API call: /jobs")
		jobs, err := scraper.ScrapeJobs()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		err = database.SaveJobsToDB(jobs)
		if err != nil {
			log.Printf("DB error (this does not prevent JSON response): %v", err)
		}

		c.JSON(200, jobs)
	})

	log.Println("Backend listening on :8080")
	r.Run(":8080")
}
