package database

import (
	"../services"
	"github.com/jackc/pgx"
	"log"
	"github.com/valyala/fasthttp"
	"math/rand"
	"github.com/satori/go.uuid"
	"github.com/pkg/errors"
	"github.com/jackc/pgx/pgtype"
	"crypto/md5"
	"strings"
)

var (
	masterConnectionPool []*pgx.ConnPool
	slaveConnectionPool []*pgx.ConnPool
)

var ErrNotUnique = errors.New("Inserting statement violates the consistency")

func init() {

	shards := serv.GetConfig().Shards
	for i := range shards {
		masterConnectionPool = append(masterConnectionPool, initPgConnPool(shards[i].Master))
		slaveConnectionPool = append(slaveConnectionPool, initPgConnPool(shards[i].Slave))
	}

}


func initPgConnPool(config serv.DataBase) *pgx.ConnPool {

	connectionConfig := pgx.ConnConfig{
		Host:     config.Host,
		Database: config.DBName,
		User:     config.User,
		Password: config.Pass,
		Port:     config.Port,
	}

	pgConnPool, err := pgx.NewConnPool(pgx.ConnPoolConfig {
		ConnConfig: connectionConfig,
		MaxConnections: serv.GetConfig().MaxConnectionsToDB,
	})

	if err != nil {
		log.Fatal(err)
	}
	return pgConnPool

}



func checkError(err error, suspicion *pgx.ConnPool) *serv.ErrorCode {
	if err == pgx.ErrNoRows {
		return &serv.ErrorCode {
			Code:    fasthttp.StatusNotFound,
			Message: fasthttp.StatusMessage(fasthttp.StatusNotFound),
		}
	}
	log.Print(err)
	if err == pgx.ErrDeadConn || strings.Contains(err.Error(), "connection refused") {
		mark: for i := range masterConnectionPool {
			if masterConnectionPool[i] == suspicion {
				masterConnectionPool[i] = slaveConnectionPool[i]
				log.Printf("Master replaced by slave (%d)", i)
				break mark
			}
			if slaveConnectionPool[i] == suspicion {
				slaveConnectionPool[i] = masterConnectionPool[i]
				log.Printf("Slave replaced by master (%d)", i)
				break mark
			}
		}
	}
	return serv.NewServerError(err)
}


func sharedKeyForWriteByID(uuid uuid.UUID) *pgx.ConnPool {
	return masterConnectionPool[int(uuid[0]) % serv.GetConfig().NumberOfShards]
}

func sharedKeyForReadByID(uuid uuid.UUID) *pgx.ConnPool {
	dbID := int(uuid[0]) % serv.GetConfig().NumberOfShards
	return choiceMasterSlave(masterConnectionPool[dbID], slaveConnectionPool[dbID])
}


func sharedKeyForWriteByMail(email string) *pgx.ConnPool {
	shardKey := md5.New().Sum([]byte(email))[0]
	return masterConnectionPool[int(shardKey) % serv.GetConfig().NumberOfShards]
}

func sharedKeyForReadByMail(email string) *pgx.ConnPool {
	shardKey := md5.New().Sum([]byte(email))[0]
	dbID := int(shardKey) % serv.GetConfig().NumberOfShards
	return choiceMasterSlave(masterConnectionPool[dbID], slaveConnectionPool[dbID])
}


func choiceMasterSlave(masterN *pgx.ConnPool, slaveN *pgx.ConnPool) *pgx.ConnPool {
	if rand.Int31n(serv.GetConfig().SlaveToMasterReadRatio + 1) == 0 {
		return masterN
	} else {
		return slaveN
	}

}

func getID() uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		log.Print(err)
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
				log.Print(err)
				return err
			}
		}
	}

	return nil
}

func castUUID(pgUUID pgtype.UUID) uuid.UUID {
	return uuid.FromBytesOrNil(pgUUID.Bytes[:])
}
