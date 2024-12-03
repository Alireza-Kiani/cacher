### Simple In-Memory Cacher With Auto Expiration

#### How to run:
##### Go
```go run ./bin/main.go```
##### Docker
```docker build -t cacher .```
```docker run --name cacher -p 3000:3000 cacher ```

#### Run Tests:
```go test ./internal/service/test```