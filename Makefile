.PHONY: build run clean

all: build run

build:
	@echo "Building..."
	@go build -o notifier ./

run:
	@echo "Running..."
	@./notifier

clean:
	@echo "Cleaning up..."
	@rm -f notifier
