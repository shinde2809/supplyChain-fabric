package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing the product lifecycle
type SmartContract struct {
    contractapi.Contract
}

// Product describes the details of the product in the supply chain
type Product struct {
    ProductID     string `json:"productID"`
    Name          string `json:"name"`
    Description   string `json:"description"`
    ManufacturingDate string `json:"manufacturingDate"`
    BatchNumber   string `json:"batchNumber"`
    Status        string `json:"status"`
    SupplyDate    string `json:"supplyDate"`
    WarehouseLocation string `json:"warehouseLocation"`
    WholesaleDate string `json:"wholesaleDate"`
    WholesaleLocation string `json:"wholesaleLocation"`
    Quantity      int    `json:"quantity"`
}

// InitLedger initializes the ledger with some sample products
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    products := []Product{
        {ProductID: "P001", Name: "Product1", Description: "Description for Product1", ManufacturingDate: "2023-09-25", BatchNumber: "B001", Status: "Created"},
        {ProductID: "P002", Name: "Product2", Description: "Description for Product2", ManufacturingDate: "2023-09-26", BatchNumber: "B002", Status: "Created"},
    }

    for _, product := range products {
        productAsBytes, _ := json.Marshal(product)
        err := ctx.GetStub().PutState(product.ProductID, productAsBytes)

        if err != nil {
            return fmt.Errorf("Failed to initialize ledger")
        }
    }
    return nil
}

func main() {
    chaincode, err := contractapi.NewChaincode(&SmartContract{})
    if err != nil {
        fmt.Printf("Error creating supply chain smart contract: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting supply chain smart contract: %s", err.Error())
    }
}
