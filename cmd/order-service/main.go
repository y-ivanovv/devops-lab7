package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"
	"time"

	"example.com/devops-lab6/internal/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	Name        string  `json:"name"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int     `json:"quantity"`
}

type Order struct {
	ID           string      `json:"id"`
	CustomerID   string      `json:"customer_id"`
	RestaurantID string      `json:"restaurant_id"`
	Status       string      `json:"status"`
	CreatedOn    int64       `json:"created_on"`
	OrderItems   []OrderItem `json:"order_items"`
}

type Service interface {
	Create(context.Context, Order) (string, error)
	GetByID(context.Context, string) (Order, error)
	ChangeStatus(context.Context, string, string) error
}

type Repository interface {
	Put(string, Order)
	Get(string) (Order, error)
	Change(string, func(*Order)) error
}

type service struct{ repo Repository }

func (s service) Create(_ context.Context, order Order) (string, error) {
	order.ID = uuid.NewString()
	order.Status = "Pending"
	order.CreatedOn = time.Now().Unix()
	s.repo.Put(order.ID, order)
	return order.ID, nil
}

func (s service) GetByID(_ context.Context, id string) (Order, error) {
	return s.repo.Get(id)
}

func (s service) ChangeStatus(_ context.Context, id, status string) error {
	if strings.TrimSpace(status) == "" {
		return api.ErrInvalid
	}
	return s.repo.Change(id, func(order *Order) { order.Status = status })
}

type statusRequest struct {
	ID     string
	Status string `json:"status"`
}

func handler(s Service) http.Handler {
	router := mux.NewRouter()
	api.Route(router, "POST", "/orders", func(ctx context.Context, request any) (any, error) {
		id, err := s.Create(ctx, request.(Order))
		return api.IDRequest{ID: id}, err
	}, api.JSON[Order])
	api.Route(router, "GET", "/orders/{id}", func(ctx context.Context, request any) (any, error) {
		return s.GetByID(ctx, request.(api.IDRequest).ID)
	}, api.ID)
	api.Route(router, "PATCH", "/orders/{id}/status", func(ctx context.Context, request any) (any, error) {
		req := request.(statusRequest)
		err := s.ChangeStatus(ctx, req.ID, req.Status)
		return api.OK{OK: err == nil}, err
	}, func(_ context.Context, r *http.Request) (any, error) {
		var req statusRequest
		err := api.ReadJSON(r, &req)
		req.ID = mux.Vars(r)["id"]
		return req, err
	})
	return router
}

func main() {
	addr := flag.String("http.addr", ":8080", "HTTP listen address")
	flag.Parse()
	log.Printf("order-service HTTP %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, handler(service{api.NewMemory[Order]()})))
}
