package database

import (
	"github.com/jackc/pgx"
	"log"
)

var conn *pgx.ConnPool

func init() {
	pgxConfig := pgx.ConnConfig{
		Host: "localhost",
		Database: "winterschool",
		User: "trubnikov",
		Password: "pass",
	}
	//pgxConfig := pgx.ConnConfig{
	//	Host: "192.168.1.149",
	//	Database: "master",
	//	User: "student11",
	//	Password: "pass",
	//}

	var err error
	conn, err = pgx.NewConnPool(pgx.ConnPoolConfig {
		ConnConfig: pgxConfig,
		MaxConnections: 50,
	})

	if err != nil {
		log.Fatal(err)
	}
}
