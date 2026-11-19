package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

var ErrNotFound = errors.New("record not found")
var ErrInvalid = errors.New("invalid request")

// Memory stores values safely across concurrent requests.
type Memory[T any] struct {
	mu    sync.RWMutex
	items map[string]T
}

func NewMemory[T any]() *Memory[T] {
	return &Memory[T]{items: make(map[string]T)}
}

func (m *Memory[T]) Put(id string, value T) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[id] = value
}

func (m *Memory[T]) Get(id string) (T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok := m.items[id]
	if !ok {
		return value, ErrNotFound
	}
	return value, nil
}

func (m *Memory[T]) Change(id string, change func(*T)) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.items[id]
	if !ok {
		return ErrNotFound
	}
	change(&value)
	m.items[id] = value
	return nil
}

type IDRequest struct {
	ID string `json:"id"`
}
type OK struct {
	OK bool `json:"ok"`
}

func Route(r *mux.Router, method, path string, e endpoint.Endpoint, decode kithttp.DecodeRequestFunc) {
	r.Handle(path, kithttp.NewServer(e, decode, encode,
		kithttp.ServerErrorEncoder(encodeError))).Methods(method)
}

func ID(_ context.Context, r *http.Request) (any, error) {
	return IDRequest{ID: mux.Vars(r)["id"]}, nil
}

func ReadJSON(r *http.Request, target any) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		return ErrInvalid
	}
	return nil
}

func JSON[T any](_ context.Context, r *http.Request) (any, error) {
	var value T
	err := ReadJSON(r, &value)
	return value, err
}

func encode(_ context.Context, w http.ResponseWriter, value any) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrInvalid):
		code = http.StatusBadRequest
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = encode(ctx, w, map[string]string{"error": err.Error()})
}
