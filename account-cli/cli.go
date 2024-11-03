package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/grpc"

	pbf "bonus-service/pb"
)

const (
	address         = "localhost:8080"
	defaultFilename = "account.json"
)

func parseFile(file string) (*pbf.CreateAccountRequest, error) {
	var account *pbf.CreateAccountRequest
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &account)
	return account, err
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("don't connect: %v", err)
	}
	defer conn.Close()
	client := pbf.NewBonusServiceClient(conn)

	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	account, err := parseFile(file)

	if err != nil {
		log.Fatalf("don't parse file %v: ", err)
	}

	r, err := client.CreateAccount(context.Background(), account)
	if err != nil {
		log.Fatalf("don't create %v: ", err)
	}
	log.Printf("Create %v: ", r.Created)
}
