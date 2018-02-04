package database

import (
	"../services"
	"github.com/jackc/pgx"
	"log"
	"github.com/valyala/fasthttp"
	"math/rand"
	"github.com/satori/go.uuid"
	"fmt"
	"github.com/pkg/errors"
	"github.com/jackc/pgx/pgtype"
)

var (
	master1 *pgx.ConnPool
	master2 *pgx.ConnPool
	slave1 *pgx.ConnPool
	slave2 *pgx.ConnPool

	masterConnectionPool []*pgx.ConnPool
	slaveConnectionPool []*pgx.ConnPool
)

var ErrNotUnique = errors.New("Inserting statement violates the consistency")

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

	masterConnectionPool = append(masterConnectionPool, master1)
	masterConnectionPool = append(masterConnectionPool, master2)

	slaveConnectionPool = append(slaveConnectionPool, slave1)
	slaveConnectionPool = append(slaveConnectionPool, slave2)
}

func initPostgresConnectionPool(config pgx.ConnConfig) *pgx.ConnPool {

	pgConnPool, err := pgx.NewConnPool(pgx.ConnPoolConfig {
		ConnConfig: config,
		MaxConnections: serv.MaxConnections,
	})

	if err != nil {
		log.Fatal(err)
	}
	return pgConnPool
}


func checkError(err error) *serv.ErrorCode {
	if err == pgx.ErrNoRows {
		return &serv.ErrorCode {
			Code:    fasthttp.StatusNotFound,
			Message: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}
	}
	log.Print(err)
	return &serv.ErrorCode {
		Code: fasthttp.StatusInternalServerError,
		Message: fasthttp.StatusMessage(fasthttp.StatusInternalServerError),
	}
}

func sharedKeyForWriteByID(uuid uuid.UUID) *pgx.ConnPool {
	return masterConnectionPool[uuid[0] % serv.NumberOfShards]
}

func sharedKeyForReadByID(uuid uuid.UUID) *pgx.ConnPool {
	dbID := uuid[0] % serv.NumberOfShards
	return choiceMasterSlave(masterConnectionPool[dbID], slaveConnectionPool[dbID])
}

func choiceMasterSlave(masterN *pgx.ConnPool, slaveN *pgx.ConnPool) *pgx.ConnPool {
	if rand.Int31n(serv.SlaveToMasterReadRate) == 0 {
		return masterN
	} else {
		return slaveN
	}
}

func getID() uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Print(err)
	}
	return id
}

func verifyUnique(sql string, ptrDest interface{}, args string) error {

	for _, master := range masterConnectionPool {
		err := master.QueryRow(sql, args).Scan(ptrDest)
		if err != pgx.ErrNoRows {
			if err == nil {
				return ErrNotUnique
			} else {
				fmt.Print(err)
				return err
			}
		}
	}

	return nil
}

func castUUID(pgUUID pgtype.UUID) uuid.UUID {
	return uuid.FromBytesOrNil(pgUUID.Bytes[:])
}
