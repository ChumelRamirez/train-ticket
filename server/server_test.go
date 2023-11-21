package main

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/ChumelRamirez/train-ticket/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestTrainTicketService(t *testing.T) {
	client := createClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	testEmail := "mscott@dundermifflin.com"

	tr := testPurchaseTicket(ctx, client, testEmail, t)
	testGetReceipt(ctx, client, testEmail, t)
	testGetSectionUsers(ctx, client, t)
	testModifySeat(ctx, tr, client, testEmail, t)
	testRemoveUser(ctx, client, testEmail, t)
}

func createClient(t *testing.T) pb.TrainTicketClient {
	bufListener := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		bufListener.Close()
	})
	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := server{}
	pb.RegisterTrainTicketServer(srv, &svc)
	go func() {
		if err := srv.Serve(bufListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	ctx := context.Background()
	dialer := func(context.Context, string) (net.Conn, error) {
		return bufListener.Dial()
	}
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	return pb.NewTrainTicketClient(conn)
}

func testPurchaseTicket(ctx context.Context, client pb.TrainTicketClient, testEmail string, t *testing.T) *pb.TicketReceipt {
	tr, err := client.PurchaseTicket(ctx, &pb.TicketRequest{
		From:      "London",
		To:        "France",
		FirstName: "Michael",
		LastName:  "Scott",
		Email:     testEmail,
	})
	if err != nil {
		t.Fatal("Fatal error testing PurchaseTicket")
	}
	if tr.SeatSection != "A" && tr.SeatSection != "B" {
		t.Errorf("PurchaseTicket Seat Section not assigned correctly.")
	}
	return tr
}

func testGetReceipt(ctx context.Context, client pb.TrainTicketClient, testEmail string, t *testing.T) {
	r, err := client.GetReceipt(ctx, &pb.UserEmail{
		Email: testEmail,
	})
	if err != nil {
		t.Fatal("Fatal error testing GetReceipt")
	}
	if r.GetEmail() != testEmail {
		t.Errorf("GetReceipt test failed.")
	}
}

func testGetSectionUsers(ctx context.Context, client pb.TrainTicketClient, t *testing.T) {
	users, err := client.GetSectionUsers(ctx, &pb.Section{
		SeatSection: "A",
	})
	if err != nil {
		t.Fatal("Fatal error testing GetSectionUsers")
	}
	if users == nil {
		t.Errorf("GetSectionUsers test failed.")
	}
}

func testModifySeat(ctx context.Context, tr *pb.TicketReceipt, client pb.TrainTicketClient, testEmail string, t *testing.T) {
	seat := tr.SeatSection
	mt, err := client.ModifyUserSeat(ctx, &pb.UserEmail{
		Email: testEmail,
	})
	if err != nil {
		t.Fatal("Fatal error testing ModifyUserSeat")
	}
	if seat == mt.SeatSection {
		t.Errorf("ModifyUserSeat test failed.")
	}
}

func testRemoveUser(ctx context.Context, client pb.TrainTicketClient, testEmail string, t *testing.T) {
	msg, err := client.RemoveUser(ctx, &pb.UserEmail{
		Email: testEmail,
	})
	if err != nil {
		t.Fatal("Fatal error testing RemoveUser")
	}
	if msg == nil {
		t.Errorf("RemoveUser test failed.")
	}
}
