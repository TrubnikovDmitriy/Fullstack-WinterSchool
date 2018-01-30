package database

import (
	"../services"
	"github.com/jackc/pgx"
	"log"
	"github.com/valyala/fasthttp"
)

var (
	master1 *pgx.ConnPool
	master2 *pgx.ConnPool
	slave1 *pgx.ConnPool
	slave2 *pgx.ConnPool
)

func init() {
	//master1 = initPostgresConnectionPool(pgx.ConnConfig{
	//	Host: "192.168.1.26",
	//	Database: "master",
	//	User: "student11",
	//	Password: "pass",
	//})
	//slave1 = initPostgresConnectionPool(pgx.ConnConfig{
	//	Host: "192.168.1.26",
	//	Database: "slave",
	//	User: "student11",
	//	Password: "pass",
	//})
	//master2 = initPostgresConnectionPool(pgx.ConnConfig{
	//		Host: "192.168.1.149",
	//		Database: "master",
	//		User: "student11",
	//		Password: "pass",
	//})
	//slave2 = initPostgresConnectionPool(pgx.ConnConfig{
	//		Host: "192.168.1.149",
	//		Database: "slave",
	//		User: "student11",
	//		Password: "pass",
	//})

	master1 = initPostgresConnectionPool(pgx.ConnConfig{
		Host: "localhost",
		Port: 5434,
		Database: "master",
		User: "student11",
		Password: "pass",
	})
	slave1 = initPostgresConnectionPool(pgx.ConnConfig{
		Host: "localhost",
		Port: 5434,
		Database: "slave",
		User: "student11",
		Password: "pass",
	})
	master2 = initPostgresConnectionPool(pgx.ConnConfig{
		Host: "localhost",
		Port: 5433,
		Database: "master",
		User: "student11",
		Password: "pass",
	})
	slave2 = initPostgresConnectionPool(pgx.ConnConfig{
		Host: "localhost",
		Port: 5433,
		Database: "slave",
		User: "student11",
		Password: "pass",
	})

}

func initPostgresConnectionPool(config pgx.ConnConfig) *pgx.ConnPool {

	pgConnPool, err := pgx.NewConnPool(pgx.ConnPoolConfig {
		ConnConfig: config,
		MaxConnections: 50,
	})

	if err != nil {
		log.Fatal(err)
	}
	return pgConnPool
}

func checkError(err error) *services.ErrorCode {
	if err == pgx.ErrNoRows {
		return services.CreateNew(fasthttp.StatusNotFound)
	}
	log.Print(err)
	return services.CreateNew(fasthttp.StatusInternalServerError)
}