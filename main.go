package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type dbconfig struct {
	DBUrl string
}

func connection(conf dbconfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.Connect(context.Background(), conf.DBUrl)
	if err != nil {
		return nil, err
	}
	return  pool, nil
}

func main() {
	dbconf := dbconfig{
		DBUrl: "postgresql://localhost/testtask?user=postgres&password=1213",
	}
	data, err := Parse()
	if err != nil {
		log.Fatal(err)
	}
	pool, err := connection(dbconf)
	if err != nil {
		log.Fatal(err)
	}

	dbmodel := Model{
		DB: pool,
	}

	err = dbmodel.InsertAll(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("All rows inserted successfully!")

}


