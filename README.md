# Receipt Processor API

A RESTful service that processes receipts and calculates points based on defined rules.

## Features

- Receipt processing endpoint
- Points calculation based on specific rules
- In-memory storage of receipts and points
- Comprehensive validation
- Well-tested components
- Containerized for easy deployment

## Architecture

This service follows a clean architecture with clear separation of concerns:
- **API Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic for calculating points
- **Model Layer**: Defines data structures and validation
- **Repository Layer**: Manages data storage (in-memory for this implementation)

## Points Calculation Rules

Points are calculated according to these rules:
1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up for points
6. 5 points if the total is greater than $10.00
7. 6 points if the purchase day is odd
8. 10 points if the time of purchase is between 2:00pm and 4:00pm

## Installation & Running

### Prerequisites
- Go 1.19+ or Docker

### Running with Go
```bash
# Clone the repository
git clone https://github.com/ycChu711/receipt-processor
cd receipt-processor

# Install dependencies
go mod download

# Run the application
go run main.go