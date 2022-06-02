
PROJECT_NAME := backend
BASE_DIR := $(shell pwd)
BIN_DIR := $(BASE_DIR)/bin

TARGET_FILE := ./main.go

run-dev: build
	@echo "run-dev $(BIN_DIR)/$(PROJECT_NAME)"
	@$(BIN_DIR)/$(PROJECT_NAME)

run: build
	@echo "run $(BIN_DIR)/$(PROJECT_NAME)"
	nohup $(BIN_DIR)/$(PROJECT_NAME) > log/run.log 2>&1 &

build: clean
	@echo "build $(TARGET_FILE)"
	go build -v -o $(BIN_DIR)/$(PROJECT_NAME) $(TARGET_FILE)
	@chmod +x $(BIN_DIR)/$(PROJECTNAME)

clean:
	@rm -rf $(BIN_DIR)/*
	@go mod tidy
