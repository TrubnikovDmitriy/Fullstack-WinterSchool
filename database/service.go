package database

import (
	"../services"
	"github.com/jackc/pgx"
	"log"
	"github.com/valyala/fasthttp"
	"math/rand"
)

var (
	master1 *pgx.ConnPool
	master2 *pgx.ConnPool
	slave1 *pgx.ConnPool
	slave2 *pgx.ConnPool
)

func init() {

	//// Prod
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

	//// SSH
	//master1 = initPostgresConnectionPool(pgx.ConnConfig{
	//	Host: "localhost",
	//	Port: 5434,
	//	Database: "master",
	//	User: "student11",
	//	Password: "pass",
	//})
	//slave1 = initPostgresConnectionPool(pgx.ConnConfig{
	//	Host: "localhost",
	//	Port: 5434,
	//	Database: "slave",
	//	User: "student11",
	//	Password: "pass",
	//})
	//master2 = initPostgresConnectionPool(pgx.ConnConfig{
	//	Host: "localhost",
	//	Port: 5433,
	//	Database: "master",
	//	User: "student11",
	//	Password: "pass",
	//})
	//slave2 = initPostgresConnectionPool(pgx.ConnConfig{
	//	Host: "localhost",
	//	Port: 5433,
	//	Database: "slave",
	//	User: "student11",
	//	Password: "pass",
	//})

	//// Local
	pgxConfig := pgx.ConnConfig{
		Host: "localhost",
		Database: "winterschool",
		User: "trubnikov",
		Password: "pass",
	}
	master1 = initPostgresConnectionPool(pgxConfig)
	slave1 = initPostgresConnectionPool(pgxConfig)
	master2 = initPostgresConnectionPool(pgxConfig)
	slave2 = initPostgresConnectionPool(pgxConfig)
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

func sharedKeyForWriteByTeamID(teamID int) *pgx.ConnPool {
	if teamID % 2 != 0 {
		return master1
	} else {
		return master2
	}
}

func sharedKeyForReadByTeamID(teamID int) *pgx.ConnPool {
	if teamID % 2 != 0 {
		return choiceMasterSlave(master1, slave1)
	} else {
		return choiceMasterSlave(master2, slave2)
	}
}

func choiceMasterSlave(masterN *pgx.ConnPool, slaveN *pgx.ConnPool) *pgx.ConnPool {
	if rand.Int31n(services.SlaveToMasterReadRate) == 0 {
		return masterN
	} else {
		return slaveN
	}
}

func getID(sql string) int {

	var newID int

	// Обращаемся к обеим базам, для инкрементирования счетчика ID
	master1.QueryRow(sql).Scan(&newID)
	master2.QueryRow(sql)

	return newID
}