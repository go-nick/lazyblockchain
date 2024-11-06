package node

import (
	"encoding/json"
	"fmt"
	"log"
)

// RPCCommandFunc TODO: Doc this
type RPCCommandFunc func(args ...string) (map[string]interface{}, error)

// MockResponse ...
func (r *RPC) MockResponse() (map[string]interface{}, error) {
	var mockResponse = map[string]interface{}{
		"chain":                "main",
		"blocks":               797095,
		"headers":              868423,
		"bestblockhash":        "000000000000000000016c9c6dc3c3f394faf0b6599e36ab746c92b510e23a34",
		"difficulty":           50646206431058.09,
		"time":                 1688446849,
		"mediantime":           1688444825,
		"verificationprogress": 0.7777889155252595,
		"initialblockdownload": true,
		"chainwork":            "00000000000000000000000000000000000000004daa54faceea33f8dc7562c0",
		"size_on_disk":         559956038638,
		"pruned":               false,
	}

	return mockResponse, nil
}

// GetBlockchainInfo TODO:
func (r *RPC) GetBlockchainInfo() (map[string]interface{}, error) {

	// Execute RawRequest for getrawtransaction
	result, err := r.Client.RawRequest(r.ctx, "getblockchaininfo", nil)
	if err != nil {
		log.Fatalf("Failed to get transaction: %v", err)
	}

	// Decode the result into a map or structured type
	var txDetails map[string]interface{}
	if err := json.Unmarshal(result, &txDetails); err != nil {
		log.Fatalf("Failed to decode transaction details: %v", err)
	}

	// Print or process transaction details
	fmt.Printf("Transaction Details: %+v\n", txDetails)

	return txDetails, err
}

// CalculateBTCFromAddress ...
func (r *RPC) CalculateBTCFromAddress(addr string) (map[string]interface{}, error) {
	// TODO: work with this later
	// r.client.ValidateAddress(r.ctx, address)

	// Convert params to json.RawMessage
	desc := fmt.Sprintf(`"addr(%s)"`, addr)
	params := []json.RawMessage{
		json.RawMessage(`"start"`),                           // Action
		json.RawMessage(fmt.Sprintf(`[{"desc": %s}]`, desc)), // Descriptor as JSON
	}

	result, err := r.Client.RawRequest(r.ctx, "scantxoutset", params)
	if err != nil {
		err = fmt.Errorf("addr: %s - error: %w", addr, err)
		return nil, err
	}

	// Parse and print the result for debugging
	var response map[string]interface{}
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetTransaction retrieves transaction details by its txHash using RawRequest
func (r *RPC) GetTransaction(txHash string) {
	// Set up the parameters for getrawtransaction
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, txHash)), // Transaction ID
		json.RawMessage(`true`),                      // Verbose flag to get detailed output
	}

	// Execute RawRequest for getrawtransaction
	result, err := r.Client.RawRequest(r.ctx, "getrawtransaction", params)
	if err != nil {
		log.Fatalf("Failed to get transaction: %v", err)
	}

	// Decode the result into a map or structured type
	var txDetails map[string]interface{}
	if err := json.Unmarshal(result, &txDetails); err != nil {
		log.Fatalf("Failed to decode transaction details: %v", err)
	}

	// Print or process transaction details
	fmt.Printf("Transaction Details: %+v\n", txDetails)
}

// func (r *RPC) GetTransaction(txHash string) {
// r.client.GetRawTransactionVerbose(c.ctx, txHash)
// }
