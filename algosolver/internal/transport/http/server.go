package http

import (
	_ "algosolver/docs"
	"algosolver/internal/models"
	"algosolver/internal/pkg"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/valyala/fasthttp"
	"gopkg.in/go-playground/validator.v9"
	"io"
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
	case "/task/solution":
		s.TaskSolution(ctx, resp)
	case "/task/list":
		s.TaskList(ctx, resp)
	case "/task/price":
		s.TaskPrice(ctx, resp)
	}
	ctx.Write(resp.FormResponse().Json())
}

// TaskSolution godoc
//
//	@Tags		Tasks
//	@Summary	Получить решение задачи.
//	@Accept		json
//	@Produce	json
//	@Param		Body body  models.TaskSolutionRequest	true	"Тело"
//	@Success	200		{object}	models.TaskSolutionResponse
//	@Router		/task/solution [post]
func (s *Server) TaskSolution(ctx *fasthttp.RequestCtx, resp *pkg.Response) {
	var taskBody models.TaskSolutionRequest
	if err := json.Unmarshal(ctx.PostBody(), &taskBody); err != nil {
		resp.SetError(1, err.Error())
		return
	}
	if err := s.validate.Struct(taskBody); err != nil {
		resp.SetError(2, err.Error())
		return
	}

	var task models.Task
	if err := s.db.QueryRow(context.Background(), `
	select tasks.type, tasks.description,
	tasks.cost,tasks.created_at,tasks.updated_at
	from tasks
	WHERE tasks.id = $1`, taskBody.Id).Scan(
		&task.Type,
		&task.Description,
		&task.Cost,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		resp.SetError(3, err.Error())
		return
	}

	res, err := task.GetSolution()
	log.Println(res)
	if err != nil {
		resp.SetError(4, err.Error())
		return
	}

	client := &http.Client{}
	marshalled, err := json.Marshal(models.UserChangeBalanceRequest{
		Id:      taskBody.UserId,
		Balance: task.Cost,
	})
	bodyReader := bytes.NewReader(marshalled)
	req, err := http.NewRequest(http.MethodGet, "http://apigateway:9090/api/auth/student/balance/add", bodyReader)
	if err != nil {
		resp.SetError(4, err.Error())
		return
	}

	response, err := client.Do(req)
	if err != nil {
		resp.SetError(5, err.Error())
		return
	}

	var resultAuth pkg.Response
	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		resp.SetError(7, err.Error())
		return
	}
	if err := json.Unmarshal(resBody, &resultAuth); err != nil {
		log.Println(resultAuth)
		resp.SetError(8, err.Error())
		return
	}

	if !resultAuth.Status {
		//TODO: можно потом привести в удобным читаемый вид все ответы от сервиса ауса
		log.Println(resultAuth)
		resp.SetError(9, "something wrong")
		return
	}

	if _, err := s.db.Query(context.Background(),
		`INSERT INTO histrory (user_id, task_id,created_at) 
			 VALUES ($1,$2,$3)`, taskBody.UserId, taskBody.Id, time.Now().UTC()); err != nil {
		resp.SetError(10, err.Error())
		return
	}
	log.Println(resultAuth)
	resp.SetValue(res)
}

// TaskList godoc
//
//	@Tags		Tasks
//	@Summary	Получить список задач.
//	@Accept		json
//	@Produce	json
//	@Param		Body body  models.UserTaskHistoryRequest	true	"Тело"
//	@Success	200		{array}	 models.UserTaskHistoryRequest
//	@Router		/task/list [post]
func (s *Server) TaskList(ctx *fasthttp.RequestCtx, resp *pkg.Response) {
	var taskBody models.UserTaskHistoryRequest
	if err := json.Unmarshal(ctx.PostBody(), &taskBody); err != nil {
		resp.SetError(1, err.Error())
		return
	}
	if err := s.validate.Struct(taskBody); err != nil {
		resp.SetError(2, err.Error())
		return
	}

	offset := (taskBody.Page - 1) * taskBody.Limit
	var task []models.UserTaskHistory
	rows, err := s.db.Query(context.Background(), `select * from histrory where user_id = $1 limit  $2 OFFSET $3`, taskBody.UserId, taskBody.Limit, offset)
	defer rows.Close()
	if err != nil {
		resp.SetError(1, err.Error())
		return
	}
	for rows.Next() {
		var history models.UserTaskHistory
		err := rows.Scan(&history.UserId, &history.TaskId, &history.CreatedAt)
		if err != nil {
			resp.SetError(2, err.Error())
			return
		}
		task = append(task, history)
	}

	resp.SetValue(task)
}

// TaskPrice godoc
//
//	@Tags		Tasks
//	@Summary	Изменить стоимость задачи
//	@Accept		json
//	@Produce	json
//	@Param		Body body  models.TaskPriceRequest	true	"Тело"
//	@Success	200		{object}	models.TaskPriceRequest
//	@Router		/task/price [post]
func (s *Server) TaskPrice(ctx *fasthttp.RequestCtx, resp *pkg.Response) {
	var taskBody models.TaskPriceRequest
	if err := json.Unmarshal(ctx.PostBody(), &taskBody); err != nil {
		resp.SetError(1, err.Error())
		return
	}
	if err := s.validate.Struct(taskBody); err != nil {
		resp.SetError(2, err.Error())
		return
	}

	if _, err := s.db.Exec(context.Background(), `
		UPDATE tasks
		SET cost = $1,updated_at = $2
		WHERE id = $3
	`, taskBody.Cost, time.Now().UTC(), taskBody.Id); err != nil {
		resp.SetError(5, err.Error())
		return
	}

	resp.SetValue(taskBody)
}

func (s *Server) FindTaskById(id string) (models.Task, error) {
	var task models.Task
	if err := s.db.QueryRow(context.Background(), `
	select tasks.type, tasks.description,
	tasks.cost,tasks.created_at,tasks.updated_at
	from users
	WHERE users.id = $1`, id).Scan(
		&task.Type,
		&task.Description,
		&task.Cost,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return task, err
	}
	return task, nil
}

func (s *Server) ServDocs() {
	router := mux.NewRouter()
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)
	if err := http.ListenAndServe("0.0.0.0:24001", router); err != nil {
		log.Println(err)
	}
}
