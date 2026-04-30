package main

import (
	"fmt"
	"log"

	"github.com/ehitelrc/slsdk"
)

func main() {
	// 1. Initialize Connection
	conn := slsdk.NewConnection(slsdk.Config{
		BaseURL:  "https://localhost:50000/b1s/v2",
		Company:  "SBODEMO",
		UserName: "manager",
		Password: "1234",
	})

	// 2. Perform Login
	fmt.Println("Attempting to login...")
	if err := conn.Login(); err != nil {
		log.Fatalf("Login failed: %v", err)
	}
	fmt.Println("Login successful!")

	// 3. Create a Business Partner Object
	bp := slsdk.NewBusinessPartner(conn)

	// Configure the fields
	bp.CardCode("c001").
		CardName("c001").
		CardType("C") // C = Customer

	// 4. Execute the Addition
	fmt.Println("Adding Business Partner...")
	resp, err := bp.Add()
	if err != nil {
		if resp != nil && resp.Error != nil {
			log.Fatalf("SAP Error [%d]: %s", resp.Error.Code, resp.Error.Message)
		}
		log.Fatalf("HTTP/Network Error: %v", err)
	}

	fmt.Printf("Business Partner Added Successfully! Result: %+v\n", resp.Data)
}
