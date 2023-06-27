package apigateway

import (
	"apigateway/internal/transport/http"
	"github.com/valyala/fasthttp"
	"log"
)

func Run() {

	s := http.NewServer()
	if err := fasthttp.ListenAndServe(":9090", s.RoutingHttp); err != nil {
		log.Fatal(err)
	}
}
