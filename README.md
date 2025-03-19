# Blockchain Client

A simple blockchain client that interfaces with the Polygon blockchain via JSON-RPC.

## Features

- JSON-RPC client for interacting with Polygon blockchain
- Exposes a REST API for simplified blockchain operations
- Docker containerization
- Terraform configuration for AWS ECS Fargate deployment

## API Endpoints

- `GET /api/block-number` - Get the latest block number
- `GET /api/block/:number` - Get details of a specific block by its number (in hex format)
- `GET /health` - Health check endpoint

## Prerequisites

- Go 1.21 or later
- Docker
- Terraform (for deployment)
- AWS CLI (for deployment)

## Local Development

### Clone the repository

```bash
git clone https://github.com/segunjkf/blockchain-client.git
cd blockchain-client
```

### Run the application locally

```bash
go mod download
go run main.go
```

The server will start at `http://localhost:8080`.

### Environment Variables

- `PORT` - Port to run the server on (default: 8080)
- `RPC_URL` - Blockchain JSON-RPC endpoint URL (default: https://polygon-rpc.com/)

### Run tests

```bash
go test ./...
```

## Docker Build

Build the Docker image:

```bash
docker build -t blockchain-client .
```

Run the container:

```bash
docker run -p 8080:8080 blockchain-client
```

## Usage Examples

### Get the latest block number

```bash
curl http://localhost:8080/api/block-number
```

Example response:
```json
{
  "jsonrpc": "2.0",
  "result": "0x134e82a",
  "id": 2
}
```

### Get block by number

```bash
curl http://localhost:8080/api/block/30000
```

Example response:
```json

  "jsonrpc": "2.0",
  "result": {
    "difficulty": "0x7",
    "extraData": "0xd58301090083626f7286676f312e3133856c696e757800000000000000000000a0cbdedf08fef9afe203a5c7e51a0dafb338a79588449216b9bdad543c83134944110b1b6e2112096a6b43b2cac1d50ed18ef97e0351624d36a2bf1eaa1b9efe00",
    "gasLimit": "0x1312d00",
    "gasUsed": "0x0",
    "hash": "0xdf7e8995b8f7d3b80ee36007622ee9c1438837ebef9bc6f4c73583f00b7d1dde",
    "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
    "miner": "0x0000000000000000000000000000000000000000",
    "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "nonce": "0x0000000000000000",
    "number": "0x7530",
    "parentHash": "0xcfdcfbc8a6fbf5840168fc4acff559ca6d419737eec4a2dcfae8181c92116866",
    "receiptsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
    "size": "0x261",
    "stateRoot": "0x1f37a8f688f30f6713aad6755d22c8aeccc277dc1930d91d6312d720ad3b93c4",
    "timestamp": "0x5ed37bb6",
    "transactions": [],
    "transactionsRoot": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
    "uncles": []
  },
  "id": 2
}
```

## Deployment to AWS ECS

The application can be deployed to AWS ECS Fargate using the provided Terraform configuration.

### Initialize Terraform

```bash
cd terraform
terraform init
```

### Plan and apply the infrastructure

```bash
terraform plan -out=tfplan
terraform apply tfplan
```

## Production Readiness Enhancements

To make this application production-ready, the following enhancements could be implemented:

### Security Enhancements

1. **API Authentication and Authorization**
   - Implement JWT-based authentication
   - Add role-based access control
   - Rate limiting to prevent abuse

2. **HTTPS Support**
   - Configure TLS/SSL certificates
   - Redirect HTTP to HTTPS
   - Implement proper security headers

3. **Input Validation**
   - Strict validation of block numbers and other inputs
   - Protection against malicious inputs

### Reliability Enhancements

1. **Multi-RPC Endpoints**
   - Support for multiple RPC endpoints with failover
   - Implement exponential backoff for retries
   - Health checks for RPC endpoints

2. **Monitoring and Observability**
   - Implement comprehensive logging with structured logs
   - Set up monitoring with Prometheus and Grafana
   - Create dashboards for key metrics (request latency, error rates)
   - Configure alerts for critical issues

3. **High Availability**
   - Multi-AZ deployment
   - Auto-scaling based on load metrics
   - Implement circuit breakers for RPC calls

### Performance Enhancements

1. **Caching**
   - Cache block information for frequently accessed blocks
   - Implement Redis or Memcached for distributed caching
   - Properly configure cache invalidation strategies

2. **Connection Pooling**
   - Implement HTTP connection pooling for RPC calls
   - Configure appropriate timeouts and keep-alive settings

3. **Load Testing**
   - Perform regular load testing to identify bottlenecks
   - Optimize based on findings

### Additional Features

1. **Extended API Endpoints**
   - Support for additional JSON-RPC methods
   - Batch requests for multiple blocks
   - Transaction filtering and search
   - Web3 compatibility layer

2. **Metrics and Analytics**
   - Track blockchain statistics over time
   - Generate reports on chain activity
   - Alert on abnormal chain behavior

### CI/CD Pipeline

1. **Continuous Integration**
   - Add GitHub Actions automated testing
   - Implement code quality checks (linting, static analysis)
   - Security scanning for dependencies

2. **Continuous Deployment**
   - Automated deployments to staging and production
   - Blue/green or canary deployment strategies
   - Automated rollbacks if issues are detected
