ifneq (,$(wildcard ./.env))
    include .env
    export
endif

ifeq ($(OS),Windows_NT)
    GET_DATE := $(shell powershell -Command "(Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')")
    RM := if exist $(BUILD_DIR) rmdir /s /q
    MKDIR := if not exist $(BUILD_DIR) mkdir
    TARGET := app.exe
	COPY := xcopy
else
    GET_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
    RM := rm -rf
    MKDIR := mkdir -p
    TARGET := app
	COPY := cp
endif

BACKEND_DIR=backend
BUILD_DIR=$(BACKEND_DIR)/build
GO_MAIN=./cmd/server/main.go
GO_DOCS=./cmd/docs/main.go
LDFLAGS=-ldflags "-X main.BuildTime=$(GET_DATE)"

FRONTEND_DIR=frontend
FRONT_IMAGE?=maywalk-frontend
FRONT_PORT?=8082

.PHONY: all back-build back-run back-clean docs back-test back-generate back-migrate back-migrate-down front-dev front-run front-build front-lint front-test front-generate help

# =================================
# |          make help            |
# | Выводит список команд          |
# =================================

# ===========
# | Команды |
# ===========
all: back-clean back-build ## Очистка и сборка проекта

help: ## Вывод списка команд
	@echo "Доступные команды:"
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

back-run: ## Запуск приложения без сборки
	cd $(BACKEND_DIR) && go run $(GO_MAIN)

back-build: ## Сборка проекта
	@echo "[$(GET_DATE)] Сборка..."
	@$(MKDIR) $(BUILD_DIR)
	cd $(BACKEND_DIR) && go build $(LDFLAGS) -o build/$(TARGET) $(GO_MAIN)
	@echo "Сборка завершена: $(BUILD_DIR)/$(TARGET)"

back-clean: ## Удаление скомпилированных файлов и очистка кэша
	@echo "Очистка..."
	@$(RM) $(BUILD_DIR)
	cd $(BACKEND_DIR) && go clean

docs: ## Запускает Swagger-UI
	cd $(BACKEND_DIR) && go run $(GO_DOCS)

back-test: ## Запуск тестов
	@echo "Запуск тестов..."
	cd $(BACKEND_DIR) && go test -v -race ./...

back-generate: ## Генерация типов
	$(COPY) api/openapi.yml $(BACKEND_DIR)/docs
	cd $(BACKEND_DIR) && oapi-codegen -config oapi-codegen.yml docs/openapi.yml

back-migrate: ## Применение миграции
	cd $(BACKEND_DIR) && migrate -path ./migrations -database $(DB_URL) up

back-migrate-down: ## Откат миграции
	cd $(BACKEND_DIR) && migrate -path ./migrations -database $(DB_URL) down

# ================
# | Фронтенд (frontend/)
# ================

front-dev: ## Запуск dev-сервера Vite (локально, :5173)
	cd $(FRONTEND_DIR) && npm run dev

front-build: ## Сборка docker-образа фронтенда
	docker build -t $(FRONT_IMAGE) -f $(FRONTEND_DIR)/Dockerfile .

front-run: ## Сборка и запуск фронтенда в docker (:$(FRONT_PORT))
	docker rm -f $(FRONT_IMAGE) 2>/dev/null || true
	docker build -t $(FRONT_IMAGE) -f $(FRONTEND_DIR)/Dockerfile .
	docker run -d --name $(FRONT_IMAGE) -p $(FRONT_PORT):80 $(FRONT_IMAGE)
	@echo "Фронтенд: http://localhost:$(FRONT_PORT)"

front-lint: ## Линтер фронтенда
	cd $(FRONTEND_DIR) && npm run lint

front-test: ## Юнит-тесты фронтенда
	cd $(FRONTEND_DIR) && npm run test

front-generate: ## Генерация TS-типов из api/openapi.yml
	cd $(FRONTEND_DIR) && npm run generate:api
