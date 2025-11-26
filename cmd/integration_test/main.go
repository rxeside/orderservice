package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"orderservice/api/clients/productinternal"
	"orderservice/api/clients/userinternal"
	"orderservice/api/server/orderinternal"
)

func main() {
	ctx := context.Background()

	// 1. Создаем подключение к ProductService (Port 8081)
	fmt.Println("Connecting to ProductService (localhost:8081)")
	prodConn, err := grpc.NewClient("localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to product: %v", err)
	}
	defer func() {
		_ = prodConn.Close()
	}()
	prodClient := productinternal.NewProductInternalServiceClient(prodConn)

	prodResp, err := prodClient.StoreProduct(ctx, &productinternal.StoreProductRequest{
		Product: &productinternal.Product{
			Name:  "Super Laptop",
			Price: 150000, // 1500.00 рубасов
		},
	})
	if err != nil {
		log.Fatalf("Failed to create product: %v", err)
	}
	fmt.Printf("Product Created. ID: %s\n", prodResp.ProductID)

	// 2. Создаем подключение к UserService (Port 8082)
	fmt.Println("\n Step 2: Connecting to UserService (localhost:8082)")
	userConn, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user: %v", err)
	}
	defer func() {
		_ = userConn.Close()
	}()
	userClient := userinternal.NewUserInternalServiceClient(userConn)

	// Создаем пользователя (Active = 1)
	userResp, err := userClient.StoreUser(ctx, &userinternal.StoreUserRequest{
		User: &userinternal.User{
			Login:    "test_buyer",
			Status:   userinternal.UserStatus_Active,
			Email:    toPtr("buyer@example.com"),
			Telegram: toPtr("@test_buyer"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("User Created. ID: %s\n", userResp.UserID)

	// 3. Создаем заказ через OrderService (Port 8084)
	fmt.Println("\n--- Step 3: Creating Order via OrderService (localhost:8084) ---")
	orderConn, err := grpc.NewClient("localhost:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order: %v", err)
	}
	defer func() {
		_ = orderConn.Close()
	}()
	orderClient := orderinternal.NewOrderInternalServiceClient(orderConn)

	orderResp, err := orderClient.CreateOrder(ctx, &orderinternal.CreateOrderRequest{
		UserID:    userResp.UserID,
		ProductID: prodResp.ProductID,
		Price:     0,
	})

	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	fmt.Printf("успех! Order Created. OrderID: %s\n", orderResp.OrderID)
}

func toPtr(s string) *string {
	return &s
}
