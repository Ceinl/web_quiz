package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"

	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func CreateDatabase(dbPath string) (*Database, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	schema := `
	CREATE TABLE IF NOT EXISTS players (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		score INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS questions (
		id TEXT PRIMARY KEY,
		question TEXT NOT NULL,
		correct_answer TEXT NOT NULL,
		incorrect_answers TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS room (
		id TEXT PRIMARY KEY,
		question_id TEXT NOT NULL,
		player_id TEXT NOT NULL,
		FOREIGN KEY (question_id) REFERENCES questions(id),
		FOREIGN KEY (player_id) REFERENCES players(id)
	);`

	_, err = db.Exec(schema)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Database{DB: db}, nihttps://www.youtube.com/watch?v=RK5ZTYCPYg8&ab_channel=ThePrimeTimel
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) CreatePlayer(id, name string, score int) error {
	_, err := d.DB.Exec("INSERT INTO players (id, name, score) VALUES (?, ?, ?)", id, name, 0)
	return err
}

func (d *Database) CreateQuestion(q Question) error {
	_, err := d.DB.Exec("INSERT INTO questions (id, question, correct_answer, incorrect_answers) VALUES (?, ?, ?, ?)",
		q.Id, q.Question, q.CorrectAnswer, q.IncorrectAnswer)
	return err
}

func (d *Database) CreateRoom(questionId, playerId string) (string, error) {
	id, _ := d.CreateUniqueRoomId()
	_, err := d.DB.Exec("INSERT INTO room (id, question_id, player_id) VALUES (?, ?, ?)", id, questionId, playerId)
	return id, err
}

func (d *Database) GetPlayer(id string) (string, int, error) {
	row := d.DB.QueryRow("SELECT name, score FROM players WHERE id = ?", id)
	var name string
	var score int
	err := row.Scan(&name, &score)
	if err != nil {
		return "", 0, err
	}
	return name, score, nil
}

func (d *Database) IncreasePlayerScore(id string, score int) error {
	_, err := d.DB.Exec("UPDATE players SET score = score + ? WHERE id = ?", score, id)
	return err
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func (d *Database) CreateUniqueRoomId() (string, error) {
	for {
		key := make([]rune, 5)
		for i := range key {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
			key[i] = letterRunes[n.Int64()]
		}
		roomId := string(key)
		var exists bool
		err := d.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM room WHERE id = ?)", roomId).Scan(&exists)
		if err != nil {
			return "", fmt.Errorf("failed to check room ID uniqueness: %w", err)
		}
		if !exists {
			return roomId, nil
		}
	}
}

// Add these new methods after existing methods

func (d *Database) ValidateContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "multipart/form-data")
}

func (d *Database) GetQuestionsByRoomId(roomId string) ([]Question, error) {
	rows, err := d.DB.Query(`
        SELECT q.id, q.question, q.correct_answer, q.incorrect_answers 
        FROM questions q 
        JOIN room r ON r.question_id = q.id 
        WHERE r.id = ?`, roomId)
	if err != nil {
		return nil, fmt.Errorf("failed to get questions: %w", err)
	}
	defer rows.Close()

	var questions []Question
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.Id, &q.Question, &q.CorrectAnswer, &q.IncorrectAnswer); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}
	return questions, nil
}

func (d *Database) ValidateRoom(roomId string) (bool, error) {
	var exists bool
	err := d.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM room WHERE id = ?)", roomId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check room existence: %w", err)
	}
	return exists, nil
}

// Add this method to your Database struct

func (d *Database) ProcessFileUpload(contentType string, fileSize int64) error {
	if contentType == "" {
		return fmt.Errorf("content type is missing")
	}

	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return fmt.Errorf("invalid content type: %s, expected multipart/form-data", contentType)
	}

	if fileSize > 10<<20 { // 10 MB limit
		return fmt.Errorf("file too large: %d bytes", fileSize)
	}

	return nil
}
