# Backend Services & API Internship Assignment

## Available endpoints

### /transaction
### Method GET

Lists all transactions.

### /balance
### Method GET

Shows current balance.

### /transfer
### Method POST

Makes a transfer.

Request body:
```json
{
  "amount": "<float or int>"
}
```

### /add
### Method POST
Adds a transaction.

Request body:
```json
{
  "amount": "<float or int>"
}
```

## Usage

To build and run the API locally:
```bash
make; build/api
```

Or with Docker:

```bash
sudo docker build -t api .

sudo docker run -p 8080:8080 api
```

Make sure to replace <container_name> with the actual name of the container.

```bash
sudo docker stop <container_name>
```