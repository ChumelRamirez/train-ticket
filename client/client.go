package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/ChumelRamirez/train-ticket/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50505", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTrainTicketClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	tr, err := c.PurchaseTicket(ctx, &pb.TicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Michael",
		LastName:  "Scott",
		Email:     "mscott@dundermifflin.com",
	})
	if err != nil {
		log.Fatalf("could not purchase ticket: %v", err)
	}
	log.Println("Ticket purchased:", tr)
}
