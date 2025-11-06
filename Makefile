# Makefile — Bets API
# Запусти: make dev / make test / make build / make docker-up

.PHONY: dev test build docker-up clean

BUILD_DIR := ./build
BINARY    := bets-api
TARGET    := $(BUILD_DIR)/$(BINARY)

MODULE    := ./cmd/bets-api

dev:
	@clear
	@echo "Running in dev mode"
	@go run $(MODULE)/main.go

build:
	@echo "Build..."
	@mkdir -p $(BUILD_DIR)
	@go build -ldflags="-s -w" -o $(TARGET) $(MODULE)
	@echo "Build complete!"
	@echo "File: $(TARGET)$(RESET)"
	@echo "Size: $$(du -h $(TARGET) | cut -f1)"

docker-up:
	@echo "Bringing up the container..."
	@docker compose up -d
	@echo "API run: http://localhost:8080"
	@echo "Health: curl http://localhost:8080/health"

clean:
	@echo "Clean build..."
	@rm -rf $(BUILD_DIR)
	@echo "Ready!"

help:
	@echo "Bets API in Makefile"
	@echo ""
	@echo "Commands:"
	@echo "make dev	- Running in dev mode"
	@echo "make build	- Build app"
	@echo "docker-up	- Bringing up the container"
	@echo "make clean	- Clean build"