package sqlite

import (
	"database/sql"

	"github.com/mahimtalukder/go-student-api/internal/config"
	_ "modernc.org/sqlite"
)

type SQlite struct {
	//All interfacer for use database
	Db *sql.DB
}

func New(cfg *config.Config) (*SQlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    email TEXT,
    age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &SQlite{
		Db: db,
	}, nil
}

func (s *SQlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}
