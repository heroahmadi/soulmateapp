
# Soulmate App

  

Tinder Bumble like app

  

## Table of Contents

  

- [Getting Started](#getting-started)

- [Prerequisites](#prerequisites)

- [Running the dependency services](#running-the-dependency-services)

- [Running app](#running-app)

- [Running the Tests](#running-the-tests)

  

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.  

### Prerequisites

 - Go 1.22.4
 - Docker & Docker Compose

### Running the dependency services

Go to the `deployment/local` folder, and run docker compose

    cd deployment/local
    docker-compose

This will create 2 new container (MongoDB & Redis) and 1 docker volume to store database data.

### Running app

Go to the project root, and run go run

    go run cmd/app/main.go

### Running the tests
The tests included in this project is a postman test. The exported file is in the root directory named `SoulmateApp.postman_collection`.
After exporting the collection, use postman runner and check the tests

![Postman Screenshot](https://i.ibb.co.com/7xhRTG4/SCR-20240610-hany.png)

## Project Structure
```go
project-root/
    ├── api/
    │   ├── common/.            # Stores some widely used variables
    │   │   ├── context_key.go  
    │   ├── handler/			# Files to handle incoming requests
    │   │   ├── user_handler.go 
    │   │   └── ...                
    │   ├── middleware/			# Authorization logic
    │   │   ├── authz_middleware.go  
    │   │   └── ...             
    │   ├── model/				# Structs that represents models used in application
    │   │   ├── entity/			# Structs that represents db model
    │   │   │	├── user.go  
    │   │   │	├── ...  
    │   │   ├── login_request.go  
    │   │   └── ...             
    │   ├── routes/				# Application routing logic
    │   │   ├── routes.go  
    ├── cmd/
    │   ├── app/
    │   │   ├── main.go  		# Stores some widely used variables
    ├── db/						# Resource for application inital data
    │   ├── user.json
    ├── deployment/				# Yaml files for deployment
    │   ├── local/				
    │   │   ├── docker-compose.yml  	
    ├── internal/
    │   ├── config/				# Configuration of connected services
    │   │   ├── mongo.go  
    │   │   └── ...             
    │   ├── migration/			# Migration script for development
    │   │   ├── migration.go  	
    │   │   └── ...             
    ├── pkg/
    │   ├── redis/				# Helper to access redis		
    │   │   ├── redis.go  
    │   │   └── ...             
    ├── .gitignore               # Gitignore file
    ├── go.mod                   # Go module file
    ├── go.sum                   # Go module dependencies file
    └── README.md                # Project README
    └── SoulmateApp.postman_collection # Postman test collection
```