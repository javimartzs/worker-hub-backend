# Import environment variables 
include .env
export $(shell sed 's/=.*//' .env)

# Wake up postgres container from Docker
start-db: 
	@docker run --name $(DB_CONTAINER_NAME) \
	-e POSTGRES_USER=$(DB_USER) \
	-e POSTGRES_PASSWORD=$(DB_PASS) \
	-e POSTGRES_DB=$(DB_NAME) \
	-p $(DB_PORT):5432 \
	-d $(DB_IMAGE) 

# Stop and remove postgres container 
stop-db: 
	@docker stop $(DB_CONTAINER_NAME) || true 
	@docker rm $(DB_CONTAINER_NAME) || true

# Build golang application 
build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(APP_NAME) $(GO_APP)

# Run built application 
run: 
	@$(BIN_DIR)/$(APP_NAME)


# Clean binary files 
clean-bin:
	@rm -rf $(BIN_DIR)