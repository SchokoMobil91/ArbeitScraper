package database

import (
	"arbeit-scraper/scraper"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func SaveJobsToDB(jobs []scraper.Job) error {
	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		connStr = "user=postgres password=123456 dbname=postgres host=localhost port=5432 sslmode=disable"
		log.Println("Using fallback connection string for local development.")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("could not connect to the database: %w", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return fmt.Errorf("database connection failed (check Docker/password): %w", err)
	}

	log.Println("Database connection successful.")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS jobs (
		id SERIAL PRIMARY KEY,
		profession TEXT,
		salary TEXT,
		company TEXT,
		location TEXT,
		start_date TEXT,
		telephone TEXT,
		email TEXT,
		short_description TEXT,
		full_description TEXT,
		ref_nr TEXT UNIQUE,
		external_link TEXT,
		application_link TEXT
	);`

	if _, err = db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("error creating jobs table: %w", err)
	}
	log.Println("Table 'jobs' is ready.")

	truncateSQL := `TRUNCATE TABLE jobs RESTART IDENTITY;`
	if _, err = db.Exec(truncateSQL); err != nil {
		return fmt.Errorf("error truncating jobs table: %w", err)
	}
	log.Println("Existing data in 'jobs' table truncated.")

	log.Printf("Starting insertion of %d jobs into the database...", len(jobs))

	insertedCount := 0
	skippedCount := 0

	for _, job := range jobs {
		insertSQL := `
		INSERT INTO jobs (
			profession, salary, company, location, start_date, telephone, email, 
			short_description, full_description, ref_nr, external_link, application_link
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		) ON CONFLICT (ref_nr) DO NOTHING;`

		result, err := db.Exec(
			insertSQL,
			job.Profession, job.Salary, job.Company, job.Location, job.StartDate, job.Telephone, job.Email,
			job.ShortDescription, job.FullDescription, job.RefNr, job.ExternalLink, job.ApplicationLink,
		)

		if err != nil {
			log.Printf("Error inserting job %s (RefNr: %s): %v", job.Profession, job.RefNr, err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()

		if rowsAffected > 0 {
			insertedCount++
		} else {
			skippedCount++
		}
	}
	log.Printf("Data insertion complete. Total Jobs: %d. Inserted: %d. Skipped (Conflict): %d.",
		len(jobs), insertedCount, skippedCount)

	return nil
}
