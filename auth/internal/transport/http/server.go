package http

import (
	_ "auth/docs"
	"auth/internal/models"
	"auth/internal/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"time"
)

type Server struct {
	validate *validator.Validate
	db       *pgxpool.Pool
}

func NewServer() (*Server, error) {
	//TODO: необходимо инициализацию подключения к базе вывести отедльный пакет репозитория
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", "test_user", "test_pass", "postgres", "5432", "test_db")
	poolConfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}
	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}
	return &Server{db: db, validate: validator.New()}, nil
}

func (s *Server) RoutingHttp(ctx *fasthttp.RequestCtx) {
	resp := pkg.InitResp()
	path := string(ctx.Path())

	defer func() {
		ctx.Response.Header.Add("Content-Type", "application/json")
	}()
	switch path {
	case "/student/new":
		s.CreateUser(ctx, resp)
	case "/student/balance/add":
		s.ChangeUserBalance(ctx, resp)
	}
	ctx.Write(resp.FormResponse().Json())
}

// CreateUser godoc
//
//	@Tags		Users
//	@Summary	Регистрация.
//	@Accept		json
//	@Produce	json
//	@Param		Body body  models.User	true	"Тело"
//	@Success	200		{object}	models.User
//	@Router		/student/new [post]
func (s *Server) CreateUser(ctx *fasthttp.RequestCtx, resp *pkg.Response) {
	var user models.User
	if err := json.Unmarshal(ctx.PostBody(), &user); err != nil {
		resp.SetError(1, err.Error())
		return
	}

	if err := s.validate.Struct(user); err != nil {
		resp.SetError(2, err.Error())
		return
	}

	user.CreateLogin()
	now := time.Now().UTC()
	if err := s.db.QueryRow(context.Background(),
		`INSERT INTO USERS (first_name, last_name,middle_name,login,group_number,balance, created_at, updated_at) 
			 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)  returning(id)`, user.FirstName, user.LastName, user.MiddleName, user.Login, user.GroupNumber, 0, now, now).Scan(&user.Id); err != nil {
		resp.SetError(3, err.Error())
		return
	}
	user.CreatedAt = now
	user.UpdatedAt = now
	resp.SetValue(user)
}

// ChangeUserBalance godoc
//
//	@Tags		Users
//	@Summary	Изменить баланс (долг) студента. Принимает положительные или отрицательные значения баланса.
//	@Accept		json
//	@Produce	json
//	@Param		Body body  models.UserChangeBalanceRequest	true	"Тело"
//	@Success	200		{object}	models.User
//	@Router		/student/balance/add [post]
func (s *Server) ChangeUserBalance(ctx *fasthttp.RequestCtx, resp *pkg.Response) {
	var userBody models.UserChangeBalanceRequest
	if err := json.Unmarshal(ctx.PostBody(), &userBody); err != nil {
		resp.SetError(1, err.Error())
		return
	}
	if err := s.validate.Struct(userBody); err != nil {
		resp.SetError(2, err.Error())
		return
	}

	var user models.User
	if err := s.db.QueryRow(context.Background(), `
	select users.first_name, users.last_name,
	users.middle_name,users.login,users.group_number,users.balance 
	from users
	WHERE users.id = $1`, userBody.Id).Scan(
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
		&user.Login,
		&user.GroupNumber,
		&user.Balance,
	); err != nil {
		resp.SetError(3, err.Error())
		return
	}
	newBalance := userBody.Balance + user.Balance
	if newBalance > 1000 || newBalance <= 0 {
		resp.SetError(4, "Limit is wrong")
		return
	}

	if _, err := s.db.Exec(context.Background(), `
		UPDATE users
		SET balance = $1,updated_at = $2
		WHERE id = $3
	`, newBalance, time.Now().UTC(), userBody.Id); err != nil {
		resp.SetError(5, err.Error())
		return
	}

	user.Balance = newBalance
	resp.SetValue(user)
}

func (s *Server) ServDocs() {
	router := mux.NewRouter()
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
	if err := http.ListenAndServe("0.0.0.0:24000", router); err != nil {
		log.Println(err)
	}
}
