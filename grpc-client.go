package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Selvakarthikeyan-7/train-ticketing-assessment/proto"
	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := proto.NewTrainTicketServiceClient(conn)

	// Example: Submitting a purchase
	receipt, err := client.SubmitPurchase(context.Background(), &proto.PurchaseRequest{
		From:      "London",
		To:        "France",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
	})
	if err != nil {
    log.Printf("SubmitPurchase failed: %v", err)
    // Handle the error gracefully, e.g., return an error to the calling code, or take appropriate action
    return err
}

// Continue with processing the receipt or other logic
fmt.Printf("Purchase successful! Receipt:\n%s\n", receipt)
	// Example: View user details
	userDetails, err := client.GetUserDetails(context.Background(), &proto.UserRequest{
		Email: "john.doe@example.com",
	})
	if err != nil {
		log.Fatalf("GetUserDetails failed: %v", err)
	}
	fmt.Printf("User Details:\n%s\n", userDetails)

	// Additional client calls can be added for other API methods

	fmt.Println("Client execution completed.")
	s := &server{
        users:    make(map[string]*proto.UserDetails),
        sections: make(map[string][]*proto.UserDetails),
    }

    // Create a gRPC server
    grpcServer := grpc.NewServer()

    // Register the TrainTicketService with the server
    proto.RegisterTrainTicketServiceServer(grpcServer, s)

    // Start the server on port 50051
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
    log.Println("Server is listening on port 50051...")

    // Serve gRPC requests
    if err := grpcServer.Serve(listener); err != nil {
        // Handle server startup error
        log.Fatalf("Failed to serve: %v", err)
    }
}
