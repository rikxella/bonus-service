package main

import (
	"log"
	"net"

	pbf "bonus-service/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

type IRepository interface {
	Create(*pbf.CreateAccountRequest) (*pbf.CreateAccountRequest, error)
}

type Repository struct {
	accounts []*pbf.CreateAccountRequest
}

func (repo *Repository) Create(account *pbf.CreateAccountRequest) (*pbf.CreateAccountRequest, error) {
	updated := append(repo.accounts, account)
	repo.accounts = updated
	return account, nil
}

type service struct {
	repo IRepository
}

func (s *service) CretaeAccount(ctx context.Context, req *pbf.CreateAccountRequest) (*pbf.CreateAccountResponse, error) {
	account, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pbf.CreateAccountResponse{Created: true, Account: account}, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pbf.RegisterBonusServiceServer(s, &service{repo})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
