package database

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/satori/go.uuid"
	"context"
)

func Test(ctx *fasthttp.RequestCtx) {

	//fmt.Printf("int[0]=%d\n\n", id.Bytes()[0])
	master1.Prepare("insert", "INSERT INTO games(id, title, about, uuid) " +
		"VALUES ($1, 'hello', 'world', $2);")
	batch := master1.BeginBatch()


	for i := 0; i < 10; i++ {
		id, _ := uuid.NewV4()
		fmt.Printf("hash=%s\n", id.String())

		batch.Queue("insert", []interface{}{ i, id },
				nil, nil)
	}

	err := batch.Send(context.Background(), nil)
	if err!= nil {
		fmt.Print("Send ")
		fmt.Print(err)
		return
	}
	_, err = batch.ExecResults()
	if err!= nil {
		fmt.Print("\nExec ")
		fmt.Print(err)
		fmt.Print("\n")
		return
	}

	batch.Close()



	//master1.Exec("INSERT INTO games(id, title, about, uuid) VALUES (1, 'hello', 'world', $1);", id)
}

