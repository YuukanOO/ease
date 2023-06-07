package todo

//go:generate go run github.com/YuukanOO/ease/cmd github.com/YuukanOO/ease/todo...

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type TodoService struct {
	todos []*Todo
}

// Builds up a new TodoService.
func NewTodoService() *TodoService {
	return &TodoService{
		todos: make([]*Todo, 0),
	}
}

type TodoCreateCommand struct {
	Text string `json:"text"`
}

// Creates a new todo with the given text content.
//
//ease:api method=POST path=/api/todos
func (s *TodoService) Create(ctx context.Context, cmd *TodoCreateCommand) (*Todo, error) {
	todo := &Todo{
		ID:        uint(len(s.todos) + 1),
		Text:      cmd.Text,
		Completed: false,
	}
	s.todos = append(s.todos, todo)

	return todo, nil
}

// Lists all todos.
//
//ease:api method=GET path=/api/todos
func (s *TodoService) List(ctx context.Context) ([]*Todo, error) {
	return s.todos, nil
}

type TodoUpdateCommand struct {
	Completed bool `json:"completed"`
}

// Updates the todo with the given id.
//
//ease:api method=PUT path=/api/todos/:id
func (s *TodoService) Update(ctx context.Context, id uint, cmd TodoUpdateCommand) (*Todo, error) {
	for _, todo := range s.todos {
		if todo.ID == id {
			todo.Completed = cmd.Completed
			return todo, nil
		}
	}

	return nil, ErrNotFound
}

// Gets the server health
//
//ease:api method=GET path=/api/_health
func HealthCheck(ctx context.Context) (string, error) {
	return "ok", nil
}
