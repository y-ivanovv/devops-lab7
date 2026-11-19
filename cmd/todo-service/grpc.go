package main

import (
	"context"
	"errors"

	"example.com/devops-lab6/internal/api"
	pb "example.com/devops-lab6/proto"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcTransport struct {
	pb.UnimplementedTodoServer
	create, getByID, markDone kitgrpc.Handler
}

func newGRPC(e endpoints) *grpcTransport {
	return &grpcTransport{
		create: kitgrpc.NewServer(e.create, func(_ context.Context, request any) (any, error) {
			req := request.(*pb.CreateRequest)
			return Task{Title: req.Title, Description: req.Description, Priority: int(req.Priority)}, nil
		}, func(_ context.Context, value any) (any, error) {
			return &pb.IDRequest{Id: value.(api.IDRequest).ID}, nil
		}),
		getByID: kitgrpc.NewServer(e.getByID, decodeID, func(_ context.Context, value any) (any, error) {
			task := value.(Task)
			return &pb.Task{Id: task.ID, Title: task.Title, Description: task.Description,
				Done: task.Done, CreatedOn: task.CreatedOn, Priority: int32(task.Priority)}, nil
		}),
		markDone: kitgrpc.NewServer(e.markDone, decodeID, func(_ context.Context, value any) (any, error) {
			return &pb.DoneResponse{Ok: value.(api.OK).OK}, nil
		}),
	}
}

func decodeID(_ context.Context, request any) (any, error) {
	return api.IDRequest{ID: request.(*pb.IDRequest).Id}, nil
}

func grpcError(err error) error {
	switch {
	case errors.Is(err, api.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, api.ErrInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (g *grpcTransport) Create(ctx context.Context, req *pb.CreateRequest) (*pb.IDRequest, error) {
	_, value, err := g.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcError(err)
	}
	return value.(*pb.IDRequest), nil
}

func (g *grpcTransport) GetByID(ctx context.Context, req *pb.IDRequest) (*pb.Task, error) {
	_, value, err := g.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcError(err)
	}
	return value.(*pb.Task), nil
}

func (g *grpcTransport) MarkDone(ctx context.Context, req *pb.IDRequest) (*pb.DoneResponse, error) {
	_, value, err := g.markDone.ServeGRPC(ctx, req)
	if err != nil {
		return nil, grpcError(err)
	}
	return value.(*pb.DoneResponse), nil
}
