// grpc-client.go
package main

import (
	"context"
	"fmt"
	"log"

	"https://github.com/Selvakarthikeyan-7/train-ticketing-assessment/blob/main/train.proto"
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
}
