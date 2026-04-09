.PHONY: up down build logs status replication-status check-replica \
        test-create test-list test-update test-delete clean help

# ─── Docker Commands ──────────────────────────────────────────────
up:
	@echo "🚀 Starting all services..."
	docker compose up --build -d
	@echo ""
	@echo "⏳ Waiting for services to be healthy (this may take ~40s)..."
	@sleep 40
	@make status

down:
	@echo "🛑 Stopping all services..."
	docker compose down

clean:
	@echo "🧹 Removing all containers, volumes, and images..."
	docker compose down -v --rmi local

build:
	@echo "🔨 Building images..."
	docker compose build

logs:
	docker compose logs -f

logs-app:
	docker compose logs -f app

logs-primary:
	docker compose logs -f mysql-primary

logs-replica:
	docker compose logs -f mysql-replica

# ─── Status & Health ──────────────────────────────────────────────
status:
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Container Status"
	@echo "═══════════════════════════════════════════"
	@docker compose ps
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  API Health Check"
	@echo "═══════════════════════════════════════════"
	@curl -s http://localhost:8080/health | python3 -m json.tool 2>/dev/null || echo "API not ready yet"

replication-status:
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Primary - Binary Log Status"
	@echo "═══════════════════════════════════════════"
	@docker exec mysql-primary mysql -uroot -prootpassword -e "SHOW MASTER STATUS\G" 2>/dev/null
	@echo ""
	@echo "═══════════════════════════════════════════"
	@echo "  Replica - Replication Status"
	@echo "═══════════════════════════════════════════"
	@docker exec mysql-replica mysql -uroot -prootpassword \
		-e "SHOW REPLICA STATUS\G" 2>/dev/null | \
		grep -E "Replica_IO_Running|Replica_SQL_Running|Seconds_Behind|Last_Error|Source_Host|Source_Port"

check-replica:
	@echo "🔍 Checking data sync between Primary and Replica..."
	@echo ""
	@echo "── Primary products count:"
	@docker exec mysql-primary mysql -uroot -prootpassword products_db \
		-e "SELECT COUNT(*) as total_products FROM products WHERE deleted_at IS NULL;" 2>/dev/null
	@echo ""
	@echo "── Replica products count:"
	@docker exec mysql-replica mysql -uroot -prootpassword products_db \
		-e "SELECT COUNT(*) as total_products FROM products WHERE deleted_at IS NULL;" 2>/dev/null

# ─── Quick API Tests ──────────────────────────────────────────────
test-list:
	@echo "📋 GET /api/products"
	@curl -s http://localhost:8080/api/products | python3 -m json.tool

test-create:
	@echo "➕ POST /api/products"
	@curl -s -X POST http://localhost:8080/api/products \
		-H "Content-Type: application/json" \
		-d '{"name":"Test Product RnD","description":"Testing replication","price":99000,"stock":10}' \
		| python3 -m json.tool

test-update:
	@echo "✏️  PUT /api/products/1"
	@curl -s -X PUT http://localhost:8080/api/products/1 \
		-H "Content-Type: application/json" \
		-d '{"name":"Laptop Dell XPS 13 (Updated)","description":"Updated via RnD test","price":17500000,"stock":20}' \
		| python3 -m json.tool

test-delete:
	@echo "🗑️  DELETE /api/products/5"
	@curl -s -X DELETE http://localhost:8080/api/products/5 | python3 -m json.tool

# ─── DB Shell Access ─────────────────────────────────────────────
shell-primary:
	docker exec -it mysql-primary mysql -uroot -prootpassword products_db

shell-replica:
	docker exec -it mysql-replica mysql -uroot -prootpassword products_db

help:
	@echo ""
	@echo "  MySQL Replication RnD - Available Commands"
	@echo "  ─────────────────────────────────────────"
	@echo "  make up                 Start all services (build + up)"
	@echo "  make down               Stop all services"
	@echo "  make clean              Remove containers + volumes + images"
	@echo "  make logs               Tail all logs"
	@echo "  make logs-app           Tail app logs only"
	@echo "  make status             Show container + API health status"
	@echo "  make replication-status Show MySQL replication thread status"
	@echo "  make check-replica      Compare row count Primary vs Replica"
	@echo "  make test-list          GET all products"
	@echo "  make test-create        POST create product"
	@echo "  make test-update        PUT update product ID 1"
	@echo "  make test-delete        DELETE product ID 5"
	@echo "  make shell-primary      MySQL shell on Primary"
	@echo "  make shell-replica      MySQL shell on Replica"
	@echo ""
EOF