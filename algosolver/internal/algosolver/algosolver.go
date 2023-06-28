package algosolver

import (
	"algosolver/internal/transport/http"
	"github.com/valyala/fasthttp"
	"log"
)

func Run() {
	//	Здесь обычно пишется функционал инициализации конфигов, подключение к базе, хэндлеров
	//	сервисов, репозиториев.

	s, err := http.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	go s.ServDocs()

	if err := fasthttp.ListenAndServe(":9092", s.RoutingHttp); err != nil {
		log.Fatal(err)
	}
}
