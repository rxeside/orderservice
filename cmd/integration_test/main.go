package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"orderservice/api/clients/paymentinternal"
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
			Price: 150000,
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

	userResp, err := userClient.StoreUser(ctx, &userinternal.StoreUserRequest{
		User: &userinternal.User{
			Login:    "rich_buyer",
			Status:   userinternal.UserStatus_Active,
			Email:    toPtr("rich@example.com"),
			Telegram: toPtr("@rich_buyer"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	fmt.Printf("✅ User Created. ID: %s\n", userResp.UserID)

	// 3. Пополняем баланс через PaymentService (Port 8083)
	fmt.Println("\n Step 3: Depositing Money via PaymentService (localhost:8083)")
	payConn, err := grpc.NewClient("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to payment: %v", err)
	}
	defer func() { _ = payConn.Close() }()
	payClient := paymentinternal.NewPaymentInternalServiceClient(payConn)

	_, err = payClient.StoreBalance(ctx, &paymentinternal.StoreBalanceRequest{
		UserID: userResp.UserID,
		Amount: 200000, // Кладем 2000.00 (больше, чем стоит ноутбук)
	})
	if err != nil {
		log.Fatalf("Failed to deposit money: %v", err)
	}
	fmt.Printf("Balance Deposited. User now has funds.\n")

	// 4. Создаем Заказ через OrderService (Port 8084)
	fmt.Println("\n 4: Creating Order via OrderService (localhost:8084)")
	orderConn, err := grpc.NewClient("localhost:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to order: %v", err)
	}
	defer func() { _ = orderConn.Close() }()
	orderClient := orderinternal.NewOrderInternalServiceClient(orderConn)

	orderResp, err := orderClient.CreateOrder(ctx, &orderinternal.CreateOrderRequest{
		UserID:    userResp.UserID,
		ProductID: prodResp.ProductID,
		Price:     0,
	})

	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	fmt.Printf("Order Created. OrderID: %s\n", orderResp.OrderID)
	fmt.Println("   Now check 'paymentservice' logs. You should see 'funds withdrawn successfully'.")
}

func toPtr(s string) *string {
	return &s
}
