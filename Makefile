SHELL=cmd.exe
APP_BINARY=luckyWheelGame.exe

up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

clean:
	@echo "Clean and Build"
	docker-compose down -v
	docker-compose build --no-cache
	docker-compose up --force-recreate
	@echo "Clean Build Done!"

down:
	@echo "Stopping docker compose..."
	docker-compose down -v
	@echo "Done!"

build:
	@echo "Building Bank Service..."
	chdir ..\bank-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${APP_BINARY} .
	@echo "Done!"

