# S3-at-Home

S3-at-Home is a simplified, local version of Amazon S3, implemented in Go. It provides a basic object storage system with a RESTful API, allowing you to store and retrieve objects using a bucket/key structure.

## Features

- NoSQL database backend (Redis) for storing objects and metadata
- Bucket management (creation, deletion, listing)
- Object operations (upload, download, delete)
- RESTful API for interacting with the storage system
- CORS support for web applications

## Prerequisites

- Go 1.16 or higher
- Redis server

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/s3-at-home.git
   cd s3-at-home
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Build the application:
   ```
   go build ./cmd/server
   ```

## Configuration

The application uses environment variables for configuration:

- `SERVER_ADDR`: The address and port the server will listen on (default: ":8080")
- `REDIS_ADDR`: The address of your Redis server (default: "localhost:6379")

You can set these environment variables before running the server, for example:

```
export SERVER_ADDR=":3000"
export REDIS_ADDR="localhost:6379"
```

## Usage

1. Start the Redis server.

2. Run the S3-at-Home server:
   ```
   ./server
   ```

3. The server will start and listen for requests on the specified port (default: 8080).

## API Endpoints

- `GET /`: List all buckets
- `PUT /:bucket`: Create a new bucket
- `DELETE /:bucket`: Remove a bucket
- `GET /:bucket`: List objects in a bucket
- `PUT /:bucket/:object`: Upload an object
- `GET /:bucket/:object`: Download an object
- `DELETE /:bucket/:object`: Remove an object

## Example Usage

Here are some example curl commands to interact with the API:

1. Create a bucket:
   ```
   curl -X PUT http://localhost:8080/my-bucket
   ```

2. Upload an object:
   ```
   curl -X PUT -H "Content-Type: text/plain" --data-binary "Hello, S3!" http://localhost:8080/my-bucket/hello.txt
   ```

3. Download an object:
   ```
   curl http://localhost:8080/my-bucket/hello.txt
   ```

4. List buckets:
   ```
   curl http://localhost:8080/
   ```

5. List objects in a bucket:
   ```
   curl http://localhost:8080/my-bucket
   ```

6. Delete an object:
   ```
   curl -X DELETE http://localhost:8080/my-bucket/hello.txt
   ```

7. Delete a bucket:
   ```
   curl -X DELETE http://localhost:8080/my-bucket
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.