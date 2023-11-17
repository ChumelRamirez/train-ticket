package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"

	pb "github.com/ChumelRamirez/train-ticket/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	port                                   = flag.Int("port", 50505, "The server port")
	sections                               = []string{"A", "B"}
	ticketMap map[string]*pb.TicketReceipt = map[string]*pb.TicketReceipt{}
	// ticketMap map[*pb.User]*pb.TicketReceipt = map[*pb.User]*pb.TicketReceipt{}
	// sectionMap map[string]*pb.Users = map[string]*pb.Users{}
)

type server struct {
	pb.UnimplementedTrainTicketServer
}

// purchase ticket function
func (s *server) PurchaseTicket(ctx context.Context, in *pb.TicketRequest) (*pb.TicketReceipt, error) {
	log.Printf("\nTicket request for: %v", in.GetEmail())
	if ticketMap[in.GetEmail()] != nil {
		log.Printf("Ticket already purchased for: %v", in.GetEmail())
		return &pb.TicketReceipt{}, status.Error(6, "ALREADY_EXISTS: ticket for this email address already purchased. Please use a different endpoint to update existing tickets.")
	} else {
		// randomizing seat section
		sectionIndex := rand.Intn(len(sections))
		seat := sections[sectionIndex]
		// seat := fmt.Sprint(rand.Intn(30)+1) + sections[sectionIndex]
		tr := &pb.TicketReceipt{
			From:        in.GetFrom(),
			To:          in.GetTo(),
			FirstName:   in.GetFirstName(),
			LastName:    in.GetLastName(),
			Email:       in.GetEmail(),
			PricePaid:   20.00,
			SeatSection: seat,
		}
		ticketMap[in.GetEmail()] = tr
		// todo delete print
		for k, v := range ticketMap {
			log.Printf("Ticket for user %v: \n%v", k, v)
		}
		return tr, nil
	}
}

// get receipt details function
func (s *server) GetReceipt(ctx context.Context, rr *pb.ReceiptRequest) (*pb.TicketReceipt, error) {
	if ticketMap[rr.GetEmail()] != nil {
		return ticketMap[rr.GetEmail()], nil
	} else {
		return &pb.TicketReceipt{}, status.Error(5, "NOT_FOUND: receipt for this email address was not found.")
	}
}

// get users by section
func (s *server) GetSectionUsers(ctx context.Context, sec *pb.Section) (*pb.Users, error) {
	users := &pb.Users{}
	for _, v := range ticketMap {
		if v.GetSeatSection() == sec.GetSeatSection() {
			users.Users = append(users.Users, &pb.Users_User{
				FirstName:   v.GetFirstName(),
				LastName:    v.GetLastName(),
				Email:       v.GetEmail(),
				SeatSection: v.GetSeatSection(),
			})
		}
	}
	return users, nil
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
