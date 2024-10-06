package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

// Constants for wallet and connection profile paths
const (
	walletPath = "path/to/wallet"           // Update with your wallet path
	ccpPath    = "path/to/connection.yaml"   // Update with your connection profile path
)

// Structs for request payloads
type CreateProductRequest struct {
	ProductID   string `json:"productID"`
	ProductName string `json:"productName"`
}

type SupplyProductRequest struct {
	ProductID string `json:"productID"`
	Status    string `json:"status"`
}

type WholesaleProductRequest struct {
	ProductID string `json:"productID"`
	Status    string `json:"status"`
}

type QueryProductRequest struct {
	ProductID string `json:"productID"`
}

type SellProductRequest struct {
	ProductID string `json:"productID"`
	BuyerInfo string `json:"buyerInfo"`
}

// Function to connect to the Fabric network
func connectToGateway(channelName, identity string) (*gateway.Gateway, *gateway.Network, error) {
	os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	// ccp, err := gateway.NewConfigFromFile(ccpPath)
	// if err != nil {
	// 	return nil, nil, fmt.Errorf("failed to create gateway from connection profile: %w", err)
	// }

	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		fmt.Printf("Failed to create wallet: %s\n", err)
		os.Exit(1)
	}

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, identity),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get network: %w", err)
	}

	return gw, network, nil
}

// Handle POST /createProduct
func CreateProductHandler(c *gin.Context) {
	var request CreateProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	gw, network, err := connectToGateway("Channel2", "ProducerUser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gw.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("createProduct", request.ProductID, request.ProductName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Handle POST /supplyProduct
func SupplyProductHandler(c *gin.Context) {
	var request SupplyProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	gw, network, err := connectToGateway("Channel2", "SupplierUser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gw.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("supplyProduct", request.ProductID, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Handle POST /wholesaleProduct
func WholesaleProductHandler(c *gin.Context) {
	var request WholesaleProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	gw, network, err := connectToGateway("Channel3", "WholesalerUser")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gw.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("wholesaleProduct", request.ProductID, request.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// Handle GET /queryProduct
func QueryProductHandler(c *gin.Context) {
	productID := c.Query("productID")

	gw, network, err := connectToGateway("Channel1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gw.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.EvaluateTransaction("queryProduct", productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": string(result)})
}

// Handle POST /sellProduct
func SellProductHandler(c *gin.Context) {
	var request SellProductRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	gw, network, err := connectToGateway("Channel1", "User1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer gw.Close()

	contract := network.GetContract("mychaincode")
	result, err := contract.SubmitTransaction("sellProduct", request.ProductID, request.BuyerInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}