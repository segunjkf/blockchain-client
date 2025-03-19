package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Default RPC endpoint
const defaultRpcURL = "https://polygon-rpc.com/"

// JSON-RPC request structure
type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params,omitempty"`
	ID      int           `json:"id"`
}

// JSON-RPC response structure
type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      int         `json:"id"`
}

// BlockchainClient represents our client for interacting with the blockchain
type BlockchainClient struct {
	rpcURL string
}

// NewBlockchainClient creates a new blockchain client
func NewBlockchainClient(rpcURL string) *BlockchainClient {
	if rpcURL == "" {
		rpcURL = defaultRpcURL
	}
	return &BlockchainClient{rpcURL: rpcURL}
}

// GetBlockNumber returns the latest block number
func (bc *BlockchainClient) GetBlockNumber() (*JSONRPCResponse, error) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		ID:      2,
	}

	return bc.sendRequest(request)
}

// GetBlockByNumber returns block details for a given block number
func (bc *BlockchainClient) GetBlockByNumber(blockNumber string, fullTransactions bool) (*JSONRPCResponse, error) {
	request := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{blockNumber, fullTransactions},
		ID:      2,
	}

	return bc.sendRequest(request)
}

// sendRequest sends a JSON-RPC request to the blockchain node
func (bc *BlockchainClient) sendRequest(request JSONRPCRequest) (*JSONRPCResponse, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(bc.rpcURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response JSONRPCResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func main() {
	// Get RPC URL from environment or use default
	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		rpcURL = defaultRpcURL
	}

	// Set default port or use environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize blockchain client
	client := NewBlockchainClient(rpcURL)

	// Set up Gin router
	router := gin.Default()

	// Define API routes
	router.GET("/api/block-number", func(c *gin.Context) {
		response, err := client.GetBlockNumber()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	router.GET("/api/block/:number", func(c *gin.Context) {
		blockNumber := c.Param("number")
		fullTx := c.DefaultQuery("full", "true") == "true"

		response, err := client.GetBlockByNumber(blockNumber, fullTx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Start the server
	log.Printf("Server starting on port %s, connecting to RPC endpoint: %s", port, rpcURL)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
