package main

import (
	"fmt"
	"log"

	"github.com/ehitelrc/slsdk"
)

func main() {
	// 1. Initialize Connection
	conn := slsdk.NewConnection(slsdk.Config{
		BaseURL:  "https://server:50000/b1s/v2",
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

	// 3. Create a Stock Transfer Object
	tran := slsdk.NewStockTransfer(conn)

	// Configure Header
	tran.Header().
		FromWarehouse("01").
		ToWarehouse("02").
		Comments("Test transfer via SLSDK")

	// Add Lines
	tran.AddLine().
		ItemCode("A0001").
		Quantity(10).
		FromWarehouse("01").
		ToWarehouse("02").
		Add()

	tran.AddLine().
		ItemCode("A0002").
		Quantity(5).
		FromWarehouse("01").
		ToWarehouse("02").
		Add()

	// 4. Execute the Addition
	fmt.Println("Adding Stock Transfer...")
	resp, err := tran.Add()
	if err != nil {
		// Output the normalized error model
		if resp != nil && resp.Error != nil {
			log.Fatalf("SAP Error [%d]: %s", resp.Error.Code, resp.Error.Message)
		}
		log.Fatalf("HTTP/Network Error: %v", err)
	}

	fmt.Printf("Stock Transfer Added Successfully! Result: %+v\n", resp.Data)

	// 5. Example usage of the Query Builder
	fmt.Println("Querying first 5 Items...")
	items, err := conn.Items().
		Select("ItemCode", "ItemName").
		Top(5).
		Get()

	if err != nil {
		log.Fatalf("Query Failed: %v", err)
	}
	fmt.Printf("Items Result: %+v\n", items)
}
