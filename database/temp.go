package database

import (
	"github.com/valyala/fasthttp"
	"fmt"
	"github.com/satori/go.uuid"
	"context"
	"time"
)


func Test2(ctx *fasthttp.RequestCtx) {

	code := ctx.QueryArgs().Peek("codes")
	fmt.Print(len(code))
	fmt.Printf(string(code) + "\n\n")

	//fasthttp.Post("")
	cookie := &fasthttp.Cookie{}
	cookie.SetKey("ws_auth")
	cookie.SetValue(string(code))
	cookie.SetExpire(time.Now().AddDate(1, 0, 0))

	ctx.Response.Header.SetCookie(cookie)

	ctx.SetStatusCode(201)
	return
}

func Test(ctx *fasthttp.RequestCtx) {


	ctx.Redirect("http://localhost:5000/test2?code=42", 302)
	return

	//fmt.Printf("int[0]=%d\n\n", id.Bytes()[0])
	masterConnectionPool[0].Prepare("insert", "INSERT INTO games(id, title, about, uuid) " +
		"VALUES ($1, 'hello', 'world', $2);")
	batch := masterConnectionPool[0].BeginBatch()


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

