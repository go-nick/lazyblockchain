package node

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

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

// GetBestBlockHash returns the hash of the best (tip) block in the most-work fully-validated chain.
func (r *RPC) GetBestBlockHash() (map[string]interface{}, error) {

	result, err := r.Client.RawRequest(r.ctx, "getbestblockhash", nil)
	if err != nil {
		return nil, err
	}

	var blockHash string
	if err := json.Unmarshal(result, &blockHash); err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"bestblockhash": blockHash,
	}

	return response, nil
}

// GetBlock fetches details of a specific block by its hash
func (r *RPC) GetBlock(args ...string) (map[string]interface{}, error) {
	blockhash := args[0]
	// Prepare the parameters for the request
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, blockhash)),
	}

	result, err := r.Client.RawRequest(r.ctx, "getblock", params)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, err
}

// GetBlockchainInfo returns an object containing various state info regarding blockchain processing.
func (r *RPC) GetBlockchainInfo() (map[string]interface{}, error) {

	result, err := r.Client.RawRequest(r.ctx, "getblockchaininfo", nil)
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := json.Unmarshal(result, &response); err != nil {
		return nil, err
	}

	return response, err
}

// GetBlockCount returns the height of the most-work fully-validated chain.
func (r *RPC) GetBlockCount() (map[string]interface{}, error) {

	result, err := r.Client.RawRequest(r.ctx, "getblockcount", nil)
	if err != nil {
		return nil, err
	}

	var blockCount int
	if err := json.Unmarshal(result, &blockCount); err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"blockCount": blockCount,
	}

	return response, err
}

// GetBlockFilter retrieve a BIP 157 content filter for a particular block.
func (r *RPC) GetBlockFilter(args ...string) (map[string]interface{}, error) {
	blockHash, filterType := args[0], args[1]
	var response map[string]interface{} = make(map[string]interface{})
	// Prepare the parameters for the request
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, blockHash)),
		json.RawMessage(fmt.Sprintf(`"%s"`, filterType)),
	}

	result, err := r.Client.RawRequest(r.ctx, "getblockfilter", params)
	if err != nil {
		response["error"] = "fail"
		return response, err
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return response, err
	}

	return response, err
}

// GetBlockHash returns hash of block in best-block-chain at height provided
func (r *RPC) GetBlockHash(args ...string) (map[string]interface{}, error) {
	height := args[0]
	var response map[string]interface{} = make(map[string]interface{})

	heightInt, err := strconv.ParseInt(height, 10, 64)
	if err != nil {
		log.Fatal(err)
		// TODO: think of a better handling here
	}

	// Prepare the parameters for the request
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`%d`, heightInt)),
	}

	result, err := r.Client.RawRequest(r.ctx, "getblockhash", params)
	if err != nil {
		response["error"] = "fail"
		return response, err
	}

	var blockHash string
	if err := json.Unmarshal(result, &blockHash); err != nil {
		return nil, err
	}

	response["blockHash"] = blockHash

	return response, err
}

// GetBlockHeader returns an Object with information about blockheader ‘hash’.
func (r *RPC) GetBlockHeader(args ...string) (map[string]interface{}, error) {
	blockHash := args[0]
	var response map[string]interface{} = make(map[string]interface{})

	// Prepare the parameters for the request
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, blockHash)),
	}

	result, err := r.Client.RawRequest(r.ctx, "getblockheader", params)
	if err != nil {
		response["error"] = "fail"
		return response, err
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return response, err
	}

	return response, err
}

// TODO: Implement optional argument #2
// GetBlockStats Compute per block statistics for a given window. All amounts are in satoshis.
// It won’t work for some heights with pruning.
func (r *RPC) GetBlockStats(args ...string) (map[string]interface{}, error) {
	hashOrHeight := args[0]
	var response map[string]interface{} = make(map[string]interface{})

	// Prepare the parameters for the request
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`"%s"`, hashOrHeight)),
	}

	result, err := r.Client.RawRequest(r.ctx, "getblockstats", params)
	if err != nil {
		response["error"] = "fail"
		return response, err
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return response, err
	}

	return response, err
}

// TODO: find a better solution for dealing with arrays.
// GetChainTips return information about all known tips in the block tree, including the main chain as well as orphaned branches.
func (r *RPC) GetChainTips() (map[string]interface{}, error) {
	response := make(map[string]interface{})
	chainTips := make([]map[string]interface{}, 0)

	result, err := r.Client.RawRequest(r.ctx, "getchaintips", nil)
	if err != nil {
		response["error"] = "fail"
		return response, err
	}

	if err := json.Unmarshal(result, &chainTips); err != nil {
		return response, err
	}

	response["result"] = chainTips

	return response, err
}

// GetChainTxStats Compute statistics about the total number and rate of transactions in the chain.
func (r *RPC) GetChainTxStats(args ...string) (map[string]interface{}, error) {
	nblocks := args[0]
	response := make(map[string]interface{})

	// TODO: add option of no parameter provided passing nil as params
	nblocksInt, err := strconv.ParseInt(nblocks, 10, 64)
	if err != nil {
		log.Fatal(err)
		// TODO: think of a better handling here
	}

	// Prepare the parameters for the request
	params := []json.RawMessage{
		json.RawMessage(fmt.Sprintf(`%d`, nblocksInt)),
	}

	result, err := r.Client.RawRequest(r.ctx, "getchaintxstats", params)
	if err != nil {
		response["error"] = "fail"
		return response, err
	}

	if err := json.Unmarshal(result, &response); err != nil {
		return response, err
	}

	return response, err
}

// ScanTxOutSet EXPERIMENTAL!
// Scans the unspent transaction output set for entries that match certain output descriptors.
func (r *RPC) ScanTxOutSet(args ...string) (map[string]interface{}, error) {
	address := args[0]
	response := make(map[string]interface{})

	desc := fmt.Sprintf(`"addr(%s)"`, address)
	params := []json.RawMessage{
		json.RawMessage(`"start"`),                           // Action
		json.RawMessage(fmt.Sprintf(`[{"desc": %s}]`, desc)), // Descriptor as JSON
	}

	result, err := r.Client.RawRequest(r.ctx, "scantxoutset", params)
	if err != nil {
		err = fmt.Errorf("addr: %s - error: %w", address, err)
		return nil, err
	}

	// Parse and print the result for debugging
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
