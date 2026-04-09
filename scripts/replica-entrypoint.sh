#!/bin/bash
# =============================================================================
# replica-entrypoint.sh
# Wraps the official MySQL entrypoint, then runs setup-replica.sh in background
# =============================================================================

set -e

# Run setup-replica.sh in background after a delay
# (MySQL needs time to fully initialize before we configure replication)
(
  sleep 10
  /usr/local/bin/setup-replica.sh
) &

# Hand off to official MySQL Docker entrypoint
exec docker-entrypoint.sh "$@"
