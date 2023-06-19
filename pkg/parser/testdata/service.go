package testdata

import (
	contextalias "context"
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type TestService struct {
	mu     sync.Mutex
	models []*TestModel
}

// Builds up a new TestService.
func NewTestService() *TestService {
	return &TestService{
		models: make([]*TestModel, 0),
	}
}

// Retrieve all TestModel.
func (s *TestService) GetAll(contextalias.Context) []*TestModel {
	return s.models
}

type TestCreateModel struct {
	Name string
}

// Create a new TestModel.
func (s *TestService) Create(ctx contextalias.Context, cmd TestCreateModel) *TestModel {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := &TestModel{
		ID:   len(s.models) + 1,
		Name: cmd.Name,
	}

	s.models = append(s.models, m)

	return m
}

// Get a TestModel by its ID.
func (s *TestService) GetByID(ctx contextalias.Context, id int) (*TestModel, error) {
	for _, m := range s.models {
		if m.ID == id {
			return m, nil
		}
	}

	return nil, ErrNotFound
}
