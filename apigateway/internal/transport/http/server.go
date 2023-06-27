package http

import (
	"apigateway/internal/pkg"
	"errors"
	"github.com/valyala/fasthttp"
	"strings"
)

// TODO: при сборке на среды (дев, тест, прод) это все вынести в динамичный конфиг
var RouteMap = map[string]string{
	"/api/auth": "http://auth:9091",
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) RoutingHttp(ctx *fasthttp.RequestCtx) {
	resp := pkg.InitResp()

	defer func() {
		ctx.Response.Header.Add("Content-Type", "application/json")
	}()
	url, prefix, err := getServicePath(string(ctx.Path()))
	if err != nil {
		resp.SetError(21001, err.Error())
		ctx.Write(resp.FormResponse().Json())
		return
	}

	k, exist := RouteMap[url]
	if !exist {
		resp.SetError(22001, "Неизвестный путь")
		ctx.Write(resp.FormResponse().Json())
		return
	}
	targetURL := k + prefix

	targetReq := fasthttp.AcquireRequest()
	targetReq.Header.SetMethodBytes(ctx.Method())
	targetReq.Header.SetContentTypeBytes(ctx.Request.Header.ContentType())
	targetReq.SetRequestURI(targetURL)
	targetReq.SetBody(ctx.PostBody())

	client := &fasthttp.Client{}
	targetResp := fasthttp.AcquireResponse()
	err = client.Do(targetReq, targetResp)
	if err != nil {
		ctx.Error("Ошибка при обработке запроса+ "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	ctx.Response.SetStatusCode(targetResp.StatusCode())
	ctx.Response.SetBody(targetResp.Body())
	ctx.Response.Header.SetContentTypeBytes(targetResp.Header.ContentType())
	fasthttp.ReleaseRequest(targetReq)
	fasthttp.ReleaseResponse(targetResp)
}
func getServicePath(requestPath string) (string, string, error) {
	path := strings.Split(requestPath, "/")
	var cleanPath []string
	for _, p := range path {
		if p != "" {
			cleanPath = append(cleanPath, p)
		}
	}

	prefix := ""
	for i := 2; i < len(cleanPath); i++ {
		prefix += "/" + cleanPath[i]
	}
	if len(cleanPath) == 0 {
		return "", "", nil
	}
	if len(cleanPath) <= 2 {
		return "", "", errors.New("not enough url prefix")
	}
	url := "/" + cleanPath[0] + "/" + cleanPath[1]

	return url, prefix, nil
}
