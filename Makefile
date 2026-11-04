ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifeq ($(OS),Windows_NT)
    GET_DATE := $(shell powershell -Command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')")
    RM := if exist $(BUILD_DIR) rmdir /s /q
    MKDIR := if not exist $(BUILD_DIR) mkdir
    TARGET := app.exe
else
    GET_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
    RM := rm -rf
    MKDIR := mkdir -p
    TARGET := app
endif

BUILD_DIR=build
GO_MAIN=./cmd/server/main.go
LDFLAGS=-ldflags "-X main.BuildTime=$(GET_DATE)"
GO_BUILD=go build $(LDFLAGS) -o $(BUILD_DIR)/$(TARGET) $(GO_MAIN)
GO_TEST=-v -race ./...

.PHONY: all build run clean test migrate migrate-down help

# =================================
# |          make help            |
# | Выводит список команд          |
# =================================

# ===========
# | Команды |
# ===========
all: clean build ## Очистка и сборка проекта

help: ## Вывод списка команд
	@echo "Доступные команды:"
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

run: ## Запуск приложения без сборки
	-go run $(GO_MAIN)

build: ## Сборка проекта
	@echo "[$(GET_DATE)] Сборка..."
	@$(MKDIR) $(BUILD_DIR)
	$(GO_BUILD)
	@echo "Сборка завершена: $(BUILD_DIR)/$(TARGET)"

clean: ## Удаление скомпилированных файлов и очистка кэша
	@echo "Очистка..."
	@$(RM) $(BUILD_DIR)
	@go clean

test: ## Запуск тестов
	@echo "Запуск тестов..."
	go test $(GO_TEST)

migrate: ## Применение миграции
	migrate -path $(DB_MIGRATION_PATH) -database $(DB_URL) up

migrate-down: ## Откат миграции
	migrate -path $(DB_MIGRATION_PATH) -database $(DB_URL) down
