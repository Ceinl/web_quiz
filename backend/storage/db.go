package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"math/big"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DB *sql.DB
}

func CreaateDatabase(dbPath string) (*Database, error) {
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

	return &Database{DB: db}, nil
}

func (d *Database) Close() {
	d.DB.Close()
}

func (d *Database) CreatePlayer(id, name string, score int) error {
	_, err := d.DB.Exec("INSERT INTO players (id, name, score) VALUES (?, ?, ?)", id, name, 0)
	return err
}

func (d *Database) CreateQuestion(id, question, correctAnswer, incorrectAnswers string) error {
	_, err := d.DB.Exec("INSERT INTO questions (id, question, correct_answer, incorrect_answers) VALUES (?, ?, ?, ?)",
		id, question, correctAnswer, incorrectAnswers)
	return err
}

func (d *Database) CreateRoom(roomId, questionId, playerId string) error {
	_, err := d.DB.Exec("INSERT INTO room (id, question_id, player_id) VALUES (?, ?, ?)", roomId, questionId, playerId)
	return err
}

func (d *Database) GetPlayer(id string) (string,int,error) {
	row := d.DB.QueryRow("SELECT name, score FROM players WHERE id = ?", id)
	var name string
	var score int
	err := row.Scan(&name,&score)
	if err != nil {
		return "",0, err
	}
	return name,score, nil 
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



