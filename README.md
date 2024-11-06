## Learning

```(.go)
	result, err := r.client.RawRequest(context.Background(), "scantxoutset", params)
	if err != nil {
		log.Fatalf("Failed to scan UTXO set: %v", err)
	}
```
- `success`: A boolean indicating whether the scan completed successfully (true in your case).
- `txouts`: The total number of UTXOs (Unspent Transaction Outputs) currently in the UTXO set. For example, 102996765 shows the total UTXOs that your node has indexed up to the scanned height.
- `height`: The block height up to which the scan was completed. In your case, 793230 is the current height of the blockchain that your node has scanned.
    - `Height`: This is the number of blocks from the very first block (called the "genesis block") up to the latest block in the chain. It represents the position of a block in the blockchain.
- `bestblock`: The block hash of the highest block (height 793230) scanned by your node, which represents the most recent block included in the UTXO set.
    - The best block hash is unique to this particular block and allows nodes to verify it against others in the network to ensure they’re synchronized with the most recent chain state.
    - Explorers typically display both the best block’s hash and height, giving users a reference to the latest verified block in the chain.
- `unspents`: An array of UTXOs found for the specific address(es) you scanned. Each element in this array is a UTXO associated with the address you specified. Here’s what each UTXO entry typically includes:
    - `txid`: The transaction ID where the UTXO was created.
        - The transaction ID (TXID) of the transaction that created this UTXO.
        - This is a unique 64-character hexadecimal string identifying the transaction in which the output was created.
    - `vout`: The output index within the transaction for the UTXO.
        - The output index within the transaction (txid) where this UTXO resides.
        - Each transaction can have multiple outputs, and vout specifies which output in that transaction this UTXO is.
        - Example: "vout": 1 indicates the UTXO is the second output (0-based index) of the transaction with the given txid.
    - `desc`: The descriptor that matches this UTXO, as specified in your scan parameters.
        - Since you provided a descriptor in the form of addr(address), this field confirms the match with that descriptor.
        - Example: "desc": "addr(1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa)#xyzabc"
    - `scriptPubKey`: The locking script (or scriptPubKey) that matches the scanned address.
        - This is the locking script or "scriptPubKey" of the UTXO.
        - It specifies the conditions under which the UTXO can be spent, commonly including the address of the receiver or a public key hash.
        - Example: "scriptPubKey": "76a914..." (in P2PKH scripts, this would be the OP_DUP OP_HASH160 format for Bitcoin addresses).
    - `amount`: The amount of Bitcoin in this specific UTXO, usually in BTC.
        - The amount of Bitcoin associated with this UTXO, usually expressed in BTC.
        - You’ll need to sum these values across all entries to calculate the total BTC balance.
        - Example: "amount": 0.0002
    - `height`: The block height at which the transaction containing this UTXO was confirmed.
        - The block height in which the transaction containing this UTXO was confirmed.
        - This helps determine the age of the UTXO, which can be relevant for checking maturity or prioritizing spending based on confirmation depth.
        - Example: "height": 650000
    - `confirmations`: The number of blocks that have been added on top of the block containing this UTXO, indicating how many confirmations it has.
        - For example, if the current blockchain height is 793230, and this UTXO’s height is 650000, then it has 793230 - 650000 = 143230 confirmations.
        - Example: "confirmations": 143230

By iterating over unspents and summing the amount field, you can calculate the total BTC held by the address. Let me know if you’d like help writing code to parse and calculate this!