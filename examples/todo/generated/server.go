// Code generated by ease; DO NOT EDIT
package main

import (
	context_ea7792 "context"
	todo_ca7678 "github.com/YuukanOO/ease/todo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Server struct {
	Router             *gin.Engine
	TodoService_9abf69 *todo_ca7678.TodoService
}

func NewServer() (s *Server, err error) {
	s = &Server{
		Router: gin.Default(),
	}
	s.TodoService_9abf69 = todo_ca7678.NewTodoService()

	s.Router.POST("/api/todos", s.Create_9e7e22)
	s.Router.GET("/api/todos", s.List_765091)
	s.Router.PUT("/api/todos/:id", s.Update_0cdc62)
	s.Router.DELETE("/api/todos/:id", s.Delete_df53ad)
	s.Router.GET("/api/_health", s.HealthCheck_06868b)
	s.Router.GET("/api/without-params", s.WithoutParams_5fd7e9)
	s.Router.GET("/api/raw", gin.WrapF(s.TodoService_9abf69.RawEndpoint))
	s.Router.GET("/api/raw-without-receiver", gin.WrapF(todo_ca7678.RawWithoutReceiver))

	return s, nil
}

func (s *Server) Listen() {
	s.Router.Run()
}

func main() {
	s, err := NewServer()

	if err != nil {
		panic(err)
	}

	s.Listen()
}

func (s *Server) Create_9e7e22(c *gin.Context) {
	var ctx context_ea7792.Context = c.Request.Context()
	var cmd todo_ca7678.TodoCreateCommand
	if !bind(c, &cmd) {
		return
	}
	result_5a2298, err := s.TodoService_9abf69.Create(
		ctx,
		cmd,
	)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result_5a2298)
}

func (s *Server) List_765091(c *gin.Context) {
	var ctx context_ea7792.Context = c.Request.Context()
	result_5a2298, err := s.TodoService_9abf69.List(
		ctx,
	)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result_5a2298)
}

func (s *Server) Update_0cdc62(c *gin.Context) {
	var ctx context_ea7792.Context = c.Request.Context()
	var id uint = paramToInt[uint](c, "id")
	var cmd todo_ca7678.TodoUpdateCommand
	if !bind(c, &cmd) {
		return
	}
	result_5a2298, err := s.TodoService_9abf69.Update(
		ctx,
		id,
		cmd,
	)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, result_5a2298)
}

func (s *Server) Delete_df53ad(c *gin.Context) {
	var id uint = paramToInt[uint](c, "id")
	err := s.TodoService_9abf69.Delete(
		id,
	)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (s *Server) HealthCheck_06868b(c *gin.Context) {
	var ctx context_ea7792.Context = c.Request.Context()
	result_5a2298 := todo_ca7678.HealthCheck(
		ctx,
	)
	c.JSON(http.StatusOK, result_5a2298)
}

func (s *Server) WithoutParams_5fd7e9(c *gin.Context) {
	s.TodoService_9abf69.WithoutParams()
	c.Status(http.StatusNoContent)
}

func paramToInt[T int | uint](c *gin.Context, name string) T {
	value, _ := strconv.Atoi(c.Param(name))
	return T(value)
}

type HttpError interface {
	error
	Status() int
}

func handleError(c *gin.Context, err error) {
	c.Error(err)

	httpErr, implementHttpErr := err.(HttpError)

	if !implementHttpErr {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(httpErr.Status(), err)
}

func bind[T any](c *gin.Context, target *T) bool {
	if err := c.ShouldBind(target); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return false
	}

	return true
}
