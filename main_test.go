package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock HTTP server for testing JSON-RPC requests
func setupMockRPCServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var response JSONRPCResponse

		switch req.Method {
		case "eth_blockNumber":
			response = JSONRPCResponse{
				JSONRPC: "2.0",
				Result:  "0x1234567",
				ID:      req.ID,
			}
		case "eth_getBlockByNumber":
			blockNumber, ok := req.Params[0].(string)
			if !ok || blockNumber == "" {
				http.Error(w, "Invalid block number", http.StatusBadRequest)
				return
			}

			// Mock block response
			response = JSONRPCResponse{
				JSONRPC: "2.0",
				Result: map[string]interface{}{
					"number":          blockNumber,
					"hash":            "0xabc123...",
					"parentHash":      "0xdef456...",
					"nonce":           "0x1234",
					"timestamp":       "0x6789",
					"transactions":    []interface{}{},
					"transactionsRoot": "0xghi789...",
				},
				ID: req.ID,
			}
		default:
			response = JSONRPCResponse{
				JSONRPC: "2.0",
				Error:   map[string]interface{}{"code": -32601, "message": "Method not found"},
				ID:      req.ID,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
}

func setupRouter(rpcURL string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	
	// Initialize blockchain client
	client := NewBlockchainClient(rpcURL)
	
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

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	
	return router
}

func TestGetBlockNumber(t *testing.T) {
	mockServer := setupMockRPCServer()
	defer mockServer.Close()
	
	router := setupRouter(mockServer.URL)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/block-number", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response JSONRPCResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "0x1234567", response.Result)
}

func TestGetBlockByNumber(t *testing.T) {
	mockServer := setupMockRPCServer()
	defer mockServer.Close()
	
	router := setupRouter(mockServer.URL)
	
	// Test getting a block by number
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/block/0x1234567", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response JSONRPCResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	
	// Verify block data
	blockData, ok := response.Result.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "0x1234567", blockData["number"])
}

func TestHealthCheck(t *testing.T) {
	mockServer := setupMockRPCServer()
	defer mockServer.Close()
	
	router := setupRouter(mockServer.URL)
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, "ok", response["status"])
}
