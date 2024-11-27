# LazyBlockchain

LazyBlockchain is a **Command-Line Interface (CLI) Terminal User Interface (TUI) application** designed to interact seamlessly with Bitcoin nodes. It allows users to send JSON-RPC commands and view blockchain data in an intuitive and interactive terminal-based interface.

![lazyblockchain](docs/lazyblockchain.jpg)

---

## Features

- **Command-Line Options**:
  - `--host="<node_ip>"`
  - `--port="<node_port(default: 8332)>"`
  - `--user="<node_user>"`
  - `--password="<node_password>"`

  If command-line options are not provided `LazyBlockchain` will check for the configuration at `~/.bitcoin/bitcoin.conf`
- **JSON-RPC Integration**:
  - Supports all JSON-RPC methods exposed by Bitcoin nodes, such as:
    - `getblockchaininfo`
    - `getblockhash`
    - `getchaintips`
    - and more!
- **Interactive TUI**:
  - A terminal interface for real-time interaction and visualization.
  - Displays logs and results in a clean and navigable UI.
  - Shortcuts for working with JSON.


---

## Installation

### Prerequisites
- [Go](https://go.dev/) 1.21.5 or higher
- A running Bitcoin node with RPC enabled (e.g., Bitcoin Core).

### Clone the Repository
```bash
git clone https://github.com/nicholasinatel/lazyblockchain.git
cd lazyblockchain
```

### Running
```bash
make run
```

> Check: make help

---

### Developer's help

VSCode launch.json example for debugging:
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "lazyblockchain",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}",
                        "args": [
                            "--host", "***.***.***.***",
                            "--port", "***",
                            "--user", "***",
                            "--password", "***"
                        ],
            "console": "integratedTerminal"
        }
    ]
}
```