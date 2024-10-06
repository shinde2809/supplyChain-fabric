package main

import (
    "encoding/json"
    "fmt"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// CreateProduct adds a new product to the ledger
func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, productID string, name string, description string, manufacturingDate string, batchNumber string) error {
    product := Product{
        ProductID: productID,
        Name: name,
        Description: description,
        ManufacturingDate: manufacturingDate,
        BatchNumber: batchNumber,
        Status: "Created",
    }

    productAsBytes, err := json.Marshal(product)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(productID, productAsBytes)
}

// SupplyProduct updates the product status with supply details
func (s *SmartContract) SupplyProduct(ctx contractapi.TransactionContextInterface, productID string, supplyDate string, warehouseLocation string) error {
    product, err := s.QueryProduct(ctx, productID)
    if err != nil {
        return err
    }

    product.SupplyDate = supplyDate
    product.WarehouseLocation = warehouseLocation
    product.Status = "Supplied"

    productAsBytes, _ := json.Marshal(product)
    return ctx.GetStub().PutState(productID, productAsBytes)
}

// WholesaleProduct updates the product with wholesale details
func (s *SmartContract) WholesaleProduct(ctx contractapi.TransactionContextInterface, productID string, wholesaleDate string, wholesaleLocation string, quantity int) error {
    product, err := s.QueryProduct(ctx, productID)
    if err != nil {
        return err
    }

    product.WholesaleDate = wholesaleDate
    product.WholesaleLocation = wholesaleLocation
    product.Quantity = quantity
    product.Status = "Wholesaled"

    productAsBytes, _ := json.Marshal(product)
    return ctx.GetStub().PutState(productID, productAsBytes)
}

// QueryProduct retrieves a product from the ledger by productID
func (s *SmartContract) QueryProduct(ctx contractapi.TransactionContextInterface, productID string) (*Product, error) {
    productAsBytes, err := ctx.GetStub().GetState(productID)
    if err != nil {
        return nil, fmt.Errorf("Failed to read from world state: %s", err.Error())
    }
    if productAsBytes == nil {
        return nil, fmt.Errorf("Product %s does not exist", productID)
    }

    var product Product
    err = json.Unmarshal(productAsBytes, &product)
    if err != nil {
        return nil, err
    }

    return &product, nil
}

// UpdateProductStatus updates the status of a product (e.g., sold)
func (s *SmartContract) UpdateProductStatus(ctx contractapi.TransactionContextInterface, productID string, status string) error {
    product, err := s.QueryProduct(ctx, productID)
    if err != nil {
        return err
    }

    product.Status = status
    productAsBytes, _ := json.Marshal(product)
    return ctx.GetStub().PutState(productID, productAsBytes)
}
