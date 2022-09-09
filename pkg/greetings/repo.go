package greetings

import (
	"database/sql"
	"fmt"
	"log"
)

type Repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repo *Repository) FindAll() ([]*Greeting, error) {
	stmt, err := repo.db.Prepare("select template from greetings")
	if err != nil {
		return nil, fmt.Errorf("failed to create prepared statement: %w", err)
	}

	gs := make([]*Greeting, 0)
	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to query db: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	for rows.Next() {
		var tplBody string
		if err := rows.Scan(&tplBody); err != nil {
			return nil, fmt.Errorf("failed to scan template body from db row: %w", err)
		}

		g, err := NewGreeting(tplBody)
		if err != nil {
			return nil, fmt.Errorf("failed to construct greeting from template: %w", err)
		}
		gs = append(gs, g)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed retireve rows: %w", err)
	}

	return gs, nil
}

func (repo *Repository) FindRandom() (*Greeting, error) {
	return nil, nil
}
