package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Model struct {
	DB *pgxpool.Pool
}

func (c *Model) Insert(data JsonData) (int, error) {
	stmt := `INSERT INTO exploit_db (date, title, type, platform, author, code) VALUES($1, $2, $3, $4, $5, $6) returning id`
	var id int
	if err := c.DB.QueryRow(context.Background(), stmt, data.Date, data.Title, data.Type, data.Platform, data.Author, data.Code).Scan(&id); err != nil {
		return 0, nil
	}
	return id, nil
}

func (m *Model) InsertAll(dataSlice []JsonData) error {
	for i := 0; i < len(dataSlice); i++ {
		id, err := m.Insert(dataSlice[i])
		if err != nil {
			return err
		}
		fmt.Printf("Inserted with id : %d\n", id)
	}
	return nil
}
