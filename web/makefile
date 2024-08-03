# Define the default ports for the application and proxy
APP_PORT ?= 8080
PROXY_PORT ?= 3001

# Development command to run the application with Gin
dev-web:
	gin -a $(APP_PORT) -p $(PROXY_PORT) run