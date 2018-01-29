package main

import (
	"./services"
	"github.com/valyala/fasthttp"
)

func main() {

	router := services.InitRouter()

	fasthttp.ListenAndServe(":5000", router.Handler)
}

