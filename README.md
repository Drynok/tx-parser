# TX Parser

Ethereum blockchain parser that allows querying transactions for addresses.

## Usage

1. Build the project:
   `make build`

1. Run the parser:
   `make up`

1. API Endpoints:

-  GET /current-block - get the current parsed block
-  POST /subscribe - subscribe to an address
-  GET /transactions?address=0x123 - get transactions for a address

## Configuration

The parser uses the following Ethereum RPC endpoint:
https://ethereum-rpc.publicnode.com

## Development

To run tests:
`make test`

## License

This project is licensed under the MIT License.
