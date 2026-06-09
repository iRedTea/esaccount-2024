swagger:
	swag init -g ./pkg/handler/handler.go -o docs
release:
	go run cmd/release/main.go