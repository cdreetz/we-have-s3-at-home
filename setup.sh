go mod download

go build ./cmd/serve

redis-server --daemonize yes

python3 -m http.server &

echo "Redis-Server running on Port: 6379"
echo "S3 Server running on Port: 8080"
echo "Please visit localhost:8000/site.html"
