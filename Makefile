APP_NAME := gourd

build: generate
	@echo "Shaping gourd..."
	go build -o $(APP_NAME)

generate:
	@echo "Generating templates"
	templ generate

serve: build
	@echo "Serving the application..."
	./$(APP_NAME) serve

clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)

test:
	@echo "Running tests..."
	ginkgo ./...

.PHONY: build generate serve clean