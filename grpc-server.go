package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"

	"github.com/Selvakarthikeyan-7/train-ticketing-assessment/proto"
	"google.golang.org/grpc"
)

type server struct {
	users      map[string]*proto.UserDetails
	userMutex  sync.Mutex
	sections   map[string][]*proto.UserDetails
	sectionMux sync.Mutex
}

// Implement gRPC methods

func (s *server) SubmitPurchase(ctx context.Context, req *proto.PurchaseRequest) (*proto.Receipt, error) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	// Generate a random seat for the user
	seat := generateRandomSeat()

	// Create a new user
	user := &proto.UserDetails{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Seat:      seat,
		PricePaid: 20.0, // Fixed price for the ticket
	}

	// Store user information
	s.users[req.Email] = user

	// Allocate seat to the section
	s.sectionMux.Lock()
	defer s.sectionMux.Unlock()
	section := determineSection(seat)
	s.sections[section] = append(s.sections[section], user)

	// Return the receipt
	return &proto.Receipt{
		From:      req.From,
		To:        req.To,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Seat:      seat,
		PricePaid: 20.0,
	}, nil
}

func (s *server) GetUserDetails(ctx context.Context, req *proto.UserRequest) (*proto.UserDetails, error) {
    s.userMutex.Lock()
    defer s.userMutex.Unlock()

    user, exists := s.users[req.Email]
    if !exists {
        // Return a specific error indicating that the user was not found
        return nil, status.Errorf(codes.NotFound, "user with email %s not found", req.Email)
    }

    return user, nil
}

func (s *server) ViewUsersBySection(ctx context.Context, req *proto.SectionRequest) (*proto.UsersList, error) {
    s.sectionMux.Lock()
    defer s.sectionMux.Unlock()

    users, exists := s.sections[req.Section]
    if !exists {
        // Return a specific error indicating that the section was not found
        return nil, status.Errorf(codes.NotFound, "section %s not found", req.Section)
    }

    return &proto.UsersList{User: users}, nil
}

func (s *server) RemoveUser(ctx context.Context, req *proto.UserRequest) (*proto.Empty, error) {
    s.userMutex.Lock()
    defer s.userMutex.Unlock()

    // Remove user from users map
    delete(s.users, req.Email)

    // Remove user from sections
    s.sectionMux.Lock()
    defer s.sectionMux.Unlock()

    for section, users := range s.sections {
        for i, user := range users {
            if user.Email == req.Email {
                s.sections[section] = append(users[:i], users[i+1:]...)
                break
            }
        }
    }

    // Return a successful response
    return &proto.Empty{}, nil
}

func (s *server) ModifyUserSeat(ctx context.Context, req *proto.ModifySeatRequest) (*proto.UserDetails, error) {
	s.userMutex.Lock()
	defer s.userMutex.Unlock()

	user, exists := s.users[req.Email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Modify user's seat
	user.Seat = req.NewSeat

	// Update user in sections
	s.sectionMux.Lock()
	defer s.sectionMux.Unlock()

	for section, users := range s.sections {
		for _, u := range users {
			if u.Email == req.Email {
				u.Seat = req.NewSeat
				break
			}
		}
	}

	return user, nil
}

func generateRandomSeat() string {
	// Generating a random seat for simplicity
	// In a real-world scenario, you may implement a more sophisticated logic
	return fmt.Sprintf("Seat-%d", rand.Intn(100))
}

// ... (previous code)

func determineSection(seat string) string {
    // In a real-world scenario, you might have more sophisticated rules
    // For example, sections could be determined based on the ticket type,
    // passenger class, or any other business-specific criteria.

    // Placeholder: Check if the seat is in a premium section
    if isPremiumSeat(seat) {
        return "Premium"
    }

    // Default section for other seats
    return "Regular"
}

func isPremiumSeat(seat string) bool {
    // Placeholder: Check if the seat is in a premium section based on business rules
    // You might have more complex logic here, involving database lookups or external services.
    return strings.Contains(seat, "Premium")
}

func main() {
    // Initialize the server
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
