package main

import (
	"context"

	"example.com/devops-lab6/internal/api"
	"github.com/go-kit/kit/endpoint"
)

type endpoints struct {
	create, getByID, markDone endpoint.Endpoint
}

func makeEndpoints(s Service) endpoints {
	return endpoints{
		create: func(ctx context.Context, request any) (any, error) {
			id, err := s.Create(ctx, request.(Task))
			return api.IDRequest{ID: id}, err
		},
		getByID: func(ctx context.Context, request any) (any, error) {
			return s.GetByID(ctx, request.(api.IDRequest).ID)
		},
		markDone: func(ctx context.Context, request any) (any, error) {
			err := s.MarkDone(ctx, request.(api.IDRequest).ID)
			return api.OK{OK: err == nil}, err
		},
	}
}
