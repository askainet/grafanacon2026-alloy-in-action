.PHONY: help start stop restart logs build clean test

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Build the Docker image
	docker compose build

start: ## Start the mission control server
	docker compose up -d --build
	@echo ""
	@echo "✅ Mission Control is online"
	@echo "   Health check: curl http://localhost:8080/health"

stop: ## Stop the mission control server
	docker compose down -v

restart: stop start ## Restart the server

alloy-reload: ## Reload the alloy configuration
	curl http://localhost:12347/-/reload

logs: ## Tail server logs
	docker compose logs -f mission-control

alloy-logs: ## Tail alloy logs
	docker compose logs -f alloy

clean: ## Clean up containers and volumes
	docker compose down -v
	rm -rf logs/

test: ## Run Go tests
	go test -v ./...

# Mission control commands
status: ## Check mission status
	@curl -s http://localhost:8080/admin/status | jq .

metrics: ## View metrics endpoint
	@curl -s http://localhost:8080/metrics

mission1: ## Start Mission 1 (High Cardinality)
	@curl -s -X POST http://localhost:8080/admin/missions/mission1/start | jq .
	@echo "Mission 1 activated: High Cardinality Explosion"

mission1-verify: ## Verify Mission 1 (check if rogue paths are filtered)
	@curl -s http://localhost:8080/admin/mission1/verify | jq .

mission2: ## Start Mission 2 (Log Routing)
	@curl -s -X POST http://localhost:8080/admin/missions/mission2/start | jq .
	@echo "Mission 2 activated: Log Routing with DEBUG logs"

mission2-verify: ## Verify Mission 2 (S3 has all logs + Loki has no DEBUG logs)
	@curl -s http://localhost:8080/admin/s3/audit-logs/verify | jq .

mission3: ## Start Mission 3 (Probabilistic Sampling)
	@curl -s -X POST http://localhost:8080/admin/missions/mission3/start | jq .
	@echo "Mission 3 activated: Trace Sampling - Probabilistic"

mission3-verify: ## Verify Mission 3 (check if probabilistic sampling is applied)
	@curl -s http://localhost:8080/admin/mission3/verify | jq .

mission4: ## Start Mission 4 (Tail Sampling)
	@curl -s -X POST http://localhost:8080/admin/missions/mission4/start | jq .
	@echo "Mission 4 activated: Tail Sampling - Access Token Recovery"

mission4-verify: ## Verify Mission 4 (check tail sampling + key fragment recovery)
	@curl -s http://localhost:8080/admin/mission4/verify | jq .

reset: ## Reset all missions
	@curl -s -X POST http://localhost:8080/admin/reset | jq .
	@echo "All missions reset"

KEY ?= ""

deaddrop: ## Attempt to view deaddrop (requires access token from Mission 4)
	@echo "Attempting to access classified deaddrop..."
	@curl -s "http://localhost:8080/api/comms/dead-drop/view?key=${KEY}" | jq .
	@echo ""
	@echo "💡 Hint: Use 'make access-token' to check recovery status, then 'make deaddrop KEY=YOUR_TOKEN_HERE'"

# S3 validation commands
access-token: ## Check Mission 4 access token recovery status
	@curl -s http://localhost:8080/admin/mission4/access-token | jq .

s3-list: ## List recent log objects from the audit-logs S3 bucket
	@curl -s http://localhost:8080/admin/s3/audit-logs | jq .
