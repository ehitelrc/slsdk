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

	// Configure the explicit fields
	bp.CardCode("c002").
		CardName("c002").
		CardType("C").
		FederalTaxID("XAXX010101000").
		Currency("USD").
		GroupCode(100).
		Phone1("555-1234").
		EmailAddress("contacto@c002.com")

	// You can still use Set() for completely custom/unmapped fields
	bp.Set("U_MyCustomField", "Hello World")

	// Add an Address (BPAddresses)
	bp.AddAddress().
		AddressName("Main Office").
		Street("123 Tech Avenue").
		City("San Francisco").
		ZipCode("94105").
		Country("US").
		State("CA").
		AddressType("bo_BillTo").
		Add()

	// Add a Contact Employee (ContactEmployees)
	bp.AddContact().
		Name("John Doe").
		E_Mail("john@c002.com").
		Phone1("555-5555").
		Position("CTO").
		Add()

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
