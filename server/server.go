package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/ChumelRamirez/train-ticket/proto"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50505, "The server port")
)

type server struct {
	pb.UnimplementedTrainTicketServer
}

// purchase ticket function
func (s *server) PurchaseTicket(ctx context.Context, in *pb.TicketRequest) (*pb.TicketReceipt, error) {
	log.Printf("Ticket request for: %v", in.GetFirstName()+" "+in.GetLastName())
	return &pb.TicketReceipt{
		From:      in.GetFrom(),
		To:        in.GetTo(),
		FirstName: in.FirstName,
		LastName:  in.LastName,
		Email:     in.Email,
		PricePaid: 20.00,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTrainTicketServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
