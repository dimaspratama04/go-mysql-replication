#!/bin/bash
# =============================================================================
# setup-replica.sh
# Auto-configure MySQL Replica to connect to Primary (GTID-based)
# Runs inside mysql-replica container after MySQL is ready
# =============================================================================

set -e

MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:-rootpassword}"
PRIMARY_HOST="${MYSQL_PRIMARY_HOST:-mysql-primary}"
PRIMARY_PORT="${MYSQL_PRIMARY_PORT:-3306}"
REPL_USER="${MYSQL_REPL_USER:-replicator}"
REPL_PASSWORD="${MYSQL_REPL_PASSWORD:-replicatorpass}"
MAX_RETRY=5
RETRY_INTERVAL=5

log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] [REPLICA-SETUP] $1"
}

# ─── Wait for replica MySQL to be ready ───────────────────────────────────────
log "Waiting for local MySQL (replica) to be ready..."
until mysqladmin ping -h 127.0.0.1 -u root -p"${MYSQL_ROOT_PASSWORD}" --silent 2>/dev/null; do
  log "  Replica MySQL not ready yet, retrying in ${RETRY_INTERVAL}s..."
  sleep $RETRY_INTERVAL
done
log "✅ Replica MySQL is ready"

# ─── Wait for primary MySQL to be reachable ───────────────────────────────────
log "Waiting for Primary MySQL at ${PRIMARY_HOST}:${PRIMARY_PORT}..."
RETRY=0
until mysqladmin ping -h "${PRIMARY_HOST}" -P "${PRIMARY_PORT}" -u root -p"${MYSQL_ROOT_PASSWORD}" --silent 2>/dev/null; do
  RETRY=$((RETRY + 1))
  if [ $RETRY -ge $MAX_RETRY ]; then
    log "❌ Primary MySQL not reachable after ${MAX_RETRY} attempts. Exiting."
    exit 1
  fi
  log "  Primary not ready (attempt $RETRY/$MAX_RETRY), retrying in ${RETRY_INTERVAL}s..."
  sleep $RETRY_INTERVAL
done
log "✅ Primary MySQL is reachable"

# ─── Check if replication already configured ──────────────────────────────────
REPLICA_STATUS=$(mysql -h 127.0.0.1 -u root -p"${MYSQL_ROOT_PASSWORD}" \
  -e "SHOW REPLICA STATUS\G" 2>/dev/null | grep "Replica_IO_Running" | awk '{print $2}')

if [ "$REPLICA_STATUS" = "Yes" ]; then
  log "✅ Replication already running. Skipping setup."
  exit 0
fi

log "Configuring replication..."

# ─── Setup replication ────────────────────────────────────────────────────────
mysql -h 127.0.0.1 -u root -p"${MYSQL_ROOT_PASSWORD}" <<EOF
-- Stop any existing replica threads
STOP REPLICA;
RESET REPLICA ALL;

-- Configure source (primary)
CHANGE REPLICATION SOURCE TO
  SOURCE_HOST='${PRIMARY_HOST}',
  SOURCE_PORT=${PRIMARY_PORT},
  SOURCE_USER='${REPL_USER}',
  SOURCE_PASSWORD='${REPL_PASSWORD}',
  SOURCE_AUTO_POSITION=1,
  GET_SOURCE_PUBLIC_KEY=1;

-- Start replication
START REPLICA;
EOF

log "✅ Replication configured. Waiting 5s for threads to start..."
sleep 5

# ─── Verify replication status ────────────────────────────────────────────────
REPLICA_IO=$(mysql -h 127.0.0.1 -u root -p"${MYSQL_ROOT_PASSWORD}" \
  -e "SHOW REPLICA STATUS\G" 2>/dev/null | grep "Replica_IO_Running" | awk '{print $2}')

REPLICA_SQL=$(mysql -h 127.0.0.1 -u root -p"${MYSQL_ROOT_PASSWORD}" \
  -e "SHOW REPLICA STATUS\G" 2>/dev/null | grep "Replica_SQL_Running:" | awk '{print $2}')

REPLICA_ERROR=$(mysql -h 127.0.0.1 -u root -p"${MYSQL_ROOT_PASSWORD}" \
  -e "SHOW REPLICA STATUS\G" 2>/dev/null | grep "Last_Error:" | head -1 | cut -d':' -f2- | xargs)

log "──────────────────────────────────────────"
log "Replica IO  Running : ${REPLICA_IO}"
log "Replica SQL Running : ${REPLICA_SQL}"
log "Last Error          : ${REPLICA_ERROR:-none}"
log "──────────────────────────────────────────"

if [ "$REPLICA_IO" = "Yes" ] && [ "$REPLICA_SQL" = "Yes" ]; then
  log "🎉 Replication is ACTIVE and healthy!"
else
  log "⚠️  Replication may have issues. Check: docker exec -it mysql-replica mysql -uroot -prootpassword -e 'SHOW REPLICA STATUS\G'"
fi
