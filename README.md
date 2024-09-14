# TX Parser

Ethereum blockchain parser that allows querying transactions for addresses.

## Usage

1. Build the project:
   `make build`

1. Run the parser:
   `make up`

1. API Endpoints:

-  GET /current-block - get the current parsed block (curl http://localhost:8080/current_block)
-  POST /subscribe - subscribe to an address (curl -X POST -H "Content-Type: application/json" -d '{"address":"0x742d35Cc6634C0532925a3b844Bc454e4438f44e"}' http://localhost:8080/subscribe)
-  GET /transactions/0x123 - get transactions for a address (curl http://localhost:8080/transactions/0x742d35Cc6634C0532925a3b844Bc454e4438f44e)

## Configuration

The parser uses the following Ethereum RPC endpoint:
https://ethereum-rpc.publicnode.com

## Development

To run tests:
`make test`

To run lint:
`make lint`

## License

This project is licensed under the MIT License.
