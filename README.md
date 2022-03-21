# Wham Platform

### Work In Progress

Platform is responsible for handling requests from the flutter mobile app and svelte web app.

The server authenticates with Google, performs CRUD operations on the firebase firestore, creates PDF's
and sends emails to clients.

# How to run 

Start the redis server

`redis-server`

Start the platform

`go run main.go`

# Tests

`go test ./...`

# Roadmap