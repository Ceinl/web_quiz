/*
	CREATE TABLE IF NOT EXISTS questions (
		id TEXT PRIMARY KEY,
		question TEXT NOT NULL,
		correct_answer TEXT NOT NULL,
		incorrect_answers TEXT NOT NULL
	);



	func (d *Database) CreateQuestion(id, question, correctAnswer, incorrectAnswers string) error {
		_, err := d.DB.Exec("INSERT INTO questions (id, question, correct_answer, incorrect_answers) VALUES (?, ?, ?, ?)",
			id, question, correctAnswer, incorrectAnswers)
		return err
	}

*/
// pasce csv file to struct of questions and add them to the db

package storage

import (
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
)

type Question struct {
	Id              string
	Question        string
	CorrectAnswer   string
	IncorrectAnswer string
}

// Update the Reader function to return the count of imported questions
func Reader(file multipart.File, db *Database) (int, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("error reading CSV: %w", err)
	}

	log.Printf("CSV file contains %d records (including header)", len(records))

	// Count of successfully imported questions
	importCount := 0

	for i, record := range records {
		if i == 0 { // Skip header row
			log.Printf("CSV header: %v", record)
			continue
		}

		if len(record) < 3 {
			log.Printf("Warning: Record at line %d has insufficient fields: %v", i+1, record)
			continue
		}

		question := Question{
			Id:              fmt.Sprintf("q_%d", i),
			Question:        record[0],
			CorrectAnswer:   record[1],
			IncorrectAnswer: strings.Join(record[2:], "|"),
		}

		log.Printf("Importing question %d: %s", i, truncateString(question.Question, 30))

		if err := db.CreateQuestion(question); err != nil {
			log.Printf("Failed to import question %d: %v", i, err)
			continue
		}

		importCount++
	}

	log.Printf("Import summary: %d of %d questions successfully imported", importCount, len(records)-1)

	return importCount, nil
}

// Helper function to truncate long strings for logging
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
