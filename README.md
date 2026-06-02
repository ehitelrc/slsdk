# SAP Business One Service Layer SDK (SLSDK)

SLSDK is a clean, extensible, and independent Go module for interacting with the SAP Business One Service Layer (OData REST API).

This SDK allows developers to interact with the Service Layer using native Go structures and fluent APIs, avoiding the complexity of manual JSON construction, session handling, and raw HTTP operations.

## Features & Benefits (Why SLSDK vs Raw Service Layer?)

Interacting directly with the SAP B1 Service Layer via raw HTTP calls can be tedious and error-prone. **SLSDK** acts as a robust middleware that solves common pain points:

- **No Manual Session Management**: Service Layer requires handling `B1SESSION` and `ROUTEID` cookies, managing timeouts, and triggering re-logins. SLSDK's **Connection Manager** tracks this automatically.
- **Fluent API & Type Safety**: Building complex, deeply nested JSON trees manually (for lines, bin allocations, batch numbers) often leads to syntax errors or invalid payloads. SLSDK provides clean, chainable Go methods (e.g., `.ItemCode("A0001").Quantity(10)`).
- **Pagination Handled Natively**: Instead of manually parsing OData's `odata.nextLink` to fetch thousands of records, SLSDK's query engine includes `.GetAll()` to automatically paginate and merge results.
- **Standardized Error Handling**: Transforms obscure, dynamic JSON error structures from SAP into typed, predictable Go `SAPError` models.
- **Native Go Integration**: Reduces boilerplate. You focus on the business logic, while SLSDK handles the HTTP and OData complexities.

## Installation

```bash
go get github.com/ehitelrc/slsdk@v0.1.1
```

## Usage Example

```go
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
	if err := conn.Login(); err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	// 3. Create a Stock Transfer Object
	tran := slsdk.NewStockTransfer(conn)

	// Configure Header and Lines
	tran.Header().
		FromWarehouse("01").
		ToWarehouse("02").
		Comments("Test transfer")

	tran.AddLine().
		ItemCode("A0001").
		Quantity(10).
		FromWarehouse("01").
		ToWarehouse("02").
		Add()

	// 4. Execute POST request
	resp, err := tran.Add()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Success! Result: %+v\n", resp.Data)
	
	// 5. Query Builder Example (OData)
	items, err := slsdk.NewQuery(conn, "Items").
		Select("ItemCode", "ItemName").
		Filter("ItemsGroupCode eq 100").
		Top(50).
		Get()
		
	if err != nil {
		log.Fatalf("Query Failed: %v", err)
	}
	fmt.Printf("Items Result: %+v\n", items.Data)
}
```
## Supported Core Objects & Examples

SLSDK provides fully typed native builders for multiple SAP entities. This prevents typos in property names and ensures data types are correct before the payload leaves your server.

### 1. Delivery Note (Entrega)
```go
delivery := slsdk.NewDeliveryNote(conn)
delivery.Header().
	CardCode("C20000").
	DocDate("2024-06-26").
	Comments("Delivery generated via SLSDK")

delivery.AddLine().
	ItemCode("A0001").
	Quantity(5).
	WarehouseCode("01").
	Add()

resp, err := delivery.Add()
if err != nil {
	log.Fatalf("Error adding delivery note: %v", err)
}
```

### 2. Purchase Delivery (Entrada de Mercancía / Recepción)
```go
pdpo := slsdk.NewPurchaseDelivery(conn)
pdpo.Header().
	CardCode("V10000").
	DocDate("2024-06-26").
	Comments("Purchase Delivery via SLSDK")

pdpo.AddLine().
	ItemCode("A0001").
	Quantity(20).
	UnitPrice(15.50).
	WarehouseCode("01").
	Add()

resp, err := pdpo.Add()
if err != nil {
	log.Fatalf("Error adding purchase delivery: %v", err)
}
```

## Advanced Queries

SLSDK abstracts several ways to read data beyond basic CRUD operations:

### OData Query Builder
The native `QueryBuilder` allows selecting, filtering, ordering, and paginating over standard endpoints. By default, Service Layer paginates results (e.g., returning 20 records at a time). 

- Use `.Get()` to fetch a single page.
- Use `.GetAll()` to automatically follow pagination links (`odata.nextLink`) and fetch all records into a single response array.

```go
// Fetch all matching records, handling pagination automatically
resp, err := slsdk.NewQuery(conn, "BusinessPartners").
    Select("CardCode", "CardName").
    Filter("CardType eq 'cCustomer'").
    OrderBy("CardCode desc").
    GetAll() // Use Get() if you only want the first page
```

### SQL Views (Microsoft SQL Server)
If you're using SAP B1 on SQL Server, you can expose and query custom SQL views.

```go
// 1. Initialize the SQL View object
myView := slsdk.NewSQLView(conn, "B1_MyCustomView")

// 2. Expose the view in Service Layer (only needed once)
myView.Expose()

// 3. Query the view like a regular entity
resp, err := myView.Query().
    Filter("TotalAmount gt 5000").
    Get()

// 4. (Optional) Unexpose when no longer needed
myView.Unexpose()
```

### Cross Joins
For complex queries without creating a view, you can use `CrossJoin` (via `QueryService_PostQuery`):

```go
resp, err := slsdk.NewCrossJoin(conn, "Orders", "Orders/DocumentLines").
    Expand("Orders($select=DocEntry),Orders/DocumentLines($select=ItemCode)").
    Filter("Orders/DocEntry eq 1").
    Get()
```


## Generic Objects (Unmapped Endpoints)

If you need to interact with a Service Layer endpoint that hasn't been explicitly mapped in the SDK yet, you can use `GenericObject` and `slsdk.Map`.

```go
// Initialize a generic object targeting any endpoint, e.g., "StockTransfers"
gen := slsdk.NewGenericObject(conn, "StockTransfers")

// Set header fields
gen.Set("DocDate", "2024-06-26").
    Set("FromWarehouse", "004").
    Set("ToWarehouse", "004")

// Append complex nested lines easily using slsdk.Map
gen.Append("StockTransferLines", slsdk.Map{
    "ItemCode":      "IMP TK 8100",
    "Quantity":      1,
    "WarehouseCode": "004",
    "StockTransferLinesBinAllocations": []slsdk.Map{
        {
            "BinAbsEntry":   12,
            "Quantity":      1,
            "BinActionType": 2,
        },
    },
})

// Execute the POST request
resp, err := gen.Add()
```

## Roadmap

- **v0.1.x** → Connection Engine, Generic Objects, Query Builder, and multiple core SAP objects (StockTransfer, Orders, Invoices, DeliveryNotes, PurchaseDeliveries, BusinessPartners, Items) (Current)
- **v0.2.0** → Advanced Query builder enhancements (expanded OData support, aggregations)
- **v0.3.0** → Extended SAP objects (Payments, Journal Entries, Inventory validations)
- **v1.0.0** → Stable SDK release

## Versioning

This project strictly adheres to [Semantic Versioning](https://semver.org/). 
- **MAJOR** version when making incompatible API changes.
- **MINOR** version when adding functionality in a backwards compatible manner.
- **PATCH** version when making backwards compatible bug fixes.
