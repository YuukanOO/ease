// Simple package to test out ease capabilities.
package todo

//go:generate go run github.com/YuukanOO/ease/cmd github.com/YuukanOO/ease/todo... github.com/YuukanOO/ease-external-example

import (
	contextalias "context"
	"errors"
	"net/http"
	"sync"
)

var ErrNotFound = errors.New("not found")

type (
	SomeInterface interface{}

	TodoService struct {
		todos  []*Todo
		mu     sync.Mutex
		logger Logger
	}
)

// Builds up a new TodoService.
func NewTodoService(l Logger) *TodoService {
	return &TodoService{
		todos:  make([]*Todo, 0),
		logger: l,
	}
}

type TodoCreateCommand struct {
	Text string `json:"text"`
}

// Creates a new todo with the given text content.
//
//ease:api method=POST path=/api/todos
func (s *TodoService) Create(ctx contextalias.Context, cmd TodoCreateCommand) (*Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
func (s *TodoService) List(ctx contextalias.Context) ([]*Todo, error) {
	return s.todos, nil
}

type TodoUpdateCommand struct {
	Completed bool `json:"completed"`
}

// Updates the todo with the given id.
//
//ease:api method=PUT path=/api/todos/:id
func (s *TodoService) Update(ctx contextalias.Context, id uint, cmd TodoUpdateCommand) (*Todo, error) {
	for _, todo := range s.todos {
		if todo.ID == id {
			todo.Completed = cmd.Completed
			return todo, nil
		}
	}

	return nil, ErrNotFound
}

var ErrOperationNotImplemented = NewAppError("operation_not_supported")

// ease:api method=DELETE path=/api/todos/:id
func (s *TodoService) Delete(id uint) error {
	return ErrOperationNotImplemented
}

// ease:api path=/api/without-params
func (s *TodoService) WithoutParams() {
	s.logger.Log("without params nor return value")
}

// ease:api method=GET path=/api/raw
func (s *TodoService) RawEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

// ease:api method=GET path=/api/raw-without-receiver
func RawWithoutReceiver(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
