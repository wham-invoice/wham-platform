# Wham Platform

### Work In Progress

Platform is responsible for handling requests from the flutter mobile app and svelte web app.

The server authenticates with Google, performs CRUD operations on the firebase firestore, creates PDF's
and sends emails to clients.

# How to run 

Start the redis server

 - `redis-server`

Start the platform

- `go run main.go`

# Tests

`go test ./...`

# Roadmap

### Packaged into container, deployed to Kubernetes cluster hosted on AWS
- Use Terraform to create the infrastructure
  
### Confirm invoices are paid by polling users bank account
- Open banking API: https://www.akahu.nz

### Handle payments in cryptocurrencies
- Look to Solana Pay