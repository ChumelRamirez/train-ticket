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
	// sectionMap map[string]*pb.Users = map[string]*pb.Users{}
)

type server struct {
	pb.UnimplementedTrainTicketServer
}

// purchase ticket function
func (s *server) PurchaseTicket(ctx context.Context, in *pb.TicketRequest) (*pb.TicketReceipt, error) {
	log.Printf("\nPurchase Ticket request for: %v", in.GetEmail())
	if ticketMap[in.GetEmail()] != nil {
		log.Printf("Ticket already exists for: %v", in.GetEmail())
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
		// todo: delete print all tickets
		// for k, v := range ticketMap {
		// 	log.Printf("Ticket for user %v: \n%v", k, v)
		// }
		log.Printf("Ticket purchased for: %v \n %v", in.GetEmail(), tr)
		return tr, nil
	}
}

// get receipt details function
func (s *server) GetReceipt(ctx context.Context, e *pb.UserEmail) (*pb.TicketReceipt, error) {
	log.Printf("\nReceipt request for: %v", e.GetEmail())
	if ticketMap[e.GetEmail()] != nil {
		return ticketMap[e.GetEmail()], nil
	}
	return &pb.TicketReceipt{}, status.Error(5, "NOT_FOUND: receipt for this email address was not found.")
}

// get users by section
func (s *server) GetSectionUsers(ctx context.Context, sec *pb.Section) (*pb.Users, error) {
	log.Printf("\nSeat Section Users request for section: %v", sec.GetSeatSection())
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
	users.SeatSection = sec.GetSeatSection()
	return users, nil
}

// remove user
func (s *server) RemoveUser(ctx context.Context, e *pb.UserEmail) (*pb.ResponseMsgString, error) {
	log.Printf("\nRemove User request for: %v", e.GetEmail())
	if ticketMap[e.GetEmail()] != nil {
		delete(ticketMap, e.GetEmail())
		log.Printf("User Removed: %v", e.GetEmail())
		return &pb.ResponseMsgString{
			ResponseMsg: "Removed user: " + e.GetEmail(),
		}, nil
	}
	log.Printf("User to remove not found: %v", e.GetEmail())
	return &pb.ResponseMsgString{
		ResponseMsg: "NOT_FOUND: " + e.GetEmail() + " user to remove from the train was not found.",
	}, status.Error(5, "NOT_FOUND: "+e.GetEmail()+" user to remove from the train was not found.")
}

// modify seat
func (s *server) ModifyUserSeat(ctx context.Context, e *pb.UserEmail) (*pb.TicketReceipt, error) {
	log.Printf("\nModify User Seat request for: %v", e.GetEmail())
	if ticketMap[e.GetEmail()] != nil {
		if ticketMap[e.GetEmail()].SeatSection == "A" {
			ticketMap[e.GetEmail()].SeatSection = "B"
		} else {
			ticketMap[e.GetEmail()].SeatSection = "A"
		}
		return ticketMap[e.GetEmail()], nil
	}
	log.Printf("User to modify seat for not found: %v", e.GetEmail())
	return nil, status.Error(5, "NOT_FOUND: "+e.GetEmail()+" user to modify seat for was not found.")
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
