package auth

import (
	"auth/internal/transport/http"
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

	if err := fasthttp.ListenAndServe(":9091", s.RoutingHttp); err != nil {
		log.Fatal(err)
	}
}
