# SAP Business One Service Layer SDK (SLSDK)

SLSDK is a clean, extensible, and independent Go module for interacting with the SAP Business One Service Layer (OData REST API).

This SDK allows developers to interact with the Service Layer using native Go structures and fluent APIs, avoiding the complexity of manual JSON construction, session handling, and raw HTTP operations.

## Features
- **Connection Manager**: Automatic session tracking via cookies, auto-relogin support.
- **Fluent API**: Build complex business objects like `StockTransfer` through simple, chainable methods.
- **Query Engine**: Native OData query builder for selecting, filtering, ordering, and paginating.
- **Unified Response & Error Model**: Standardized response handling and typed `SAPError` structures.

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

- **v0.1.0** → Connection Engine + StockTransfer Object (Current)
- **v0.2.0** → Query builder enhancements (expanded OData support)
- **v0.3.0** → More SAP objects (Purchase Orders, Delivery Notes, Inventory)
- **v1.0.0** → Stable SDK release

## Versioning

This project strictly adheres to [Semantic Versioning](https://semver.org/). 
- **MAJOR** version when making incompatible API changes.
- **MINOR** version when adding functionality in a backwards compatible manner.
- **PATCH** version when making backwards compatible bug fixes.
