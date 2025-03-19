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
curl http://localhost:8080/api/block/0x134e82a
```

Example response:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "number": "0x134e82a",
    "hash": "0x...",
    "parentHash": "0x...",
    ...
    "transactions": [...]
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
   - Add GitHub Actions or GitLab CI for automated testing
   - Implement code quality checks (linting, static analysis)
   - Security scanning for dependencies

2. **Continuous Deployment**
   - Automated deployments to staging and production
   - Blue/green or canary deployment strategies
   - Automated rollbacks if issues are detected

### Operational Enhancements

1. **Documentation**
   - Comprehensive API documentation with Swagger/OpenAPI
   - Runbooks for common operational tasks
   - Incident response procedures

2. **Health Checks**
   - More thorough health checks including RPC endpoint checks
   - Graceful shutdown handling

3. **Infrastructure as Code**
   - Complete infrastructure defined in Terraform or CloudFormation
   - Immutable infrastructure patterns

4. **Disaster Recovery**
   - Cross-region replication for disaster recovery
   - Regular DR drills

## License

MIT