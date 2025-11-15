#!/bin/bash

# ============================================
# Login Attempts Cleanup Script
# ============================================
# This script deletes login attempt records older than 30 days
# to prevent database bloat and maintain performance.
#
# Usage: ./cleanup-login-attempts.sh
# Recommended: Run weekly via cron
# ============================================

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_FILE="${PROJECT_ROOT}/configs/config.yaml"
LOG_FILE="${PROJECT_ROOT}/logs/cleanup.log"
RETENTION_DAYS=30

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    local level=$1
    shift
    local message="$@"
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    echo -e "${timestamp} [${level}] ${message}" | tee -a "$LOG_FILE"
}

log_info() {
    log "INFO" "${GREEN}$@${NC}"
}

log_warn() {
    log "WARN" "${YELLOW}$@${NC}"
}

log_error() {
    log "ERROR" "${RED}$@${NC}"
}

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    log_error "Configuration file not found: $CONFIG_FILE"
    exit 1
fi

# Extract database configuration from YAML
# Note: This is a simple parser. For production, consider using yq or similar tools
DB_HOST=$(grep -A 10 "^database:" "$CONFIG_FILE" | grep "host:" | awk '{print $2}' | tr -d '"')
DB_PORT=$(grep -A 10 "^database:" "$CONFIG_FILE" | grep "port:" | awk '{print $2}' | tr -d '"')
DB_USER=$(grep -A 10 "^database:" "$CONFIG_FILE" | grep "user:" | awk '{print $2}' | tr -d '"')
DB_NAME=$(grep -A 10 "^database:" "$CONFIG_FILE" | grep "dbname:" | awk '{print $2}' | tr -d '"')

# Check for database password in environment variable or config
if [ -n "$DB_PASSWORD" ]; then
    DB_PASS="$DB_PASSWORD"
else
    DB_PASS=$(grep -A 10 "^database:" "$CONFIG_FILE" | grep "password:" | awk '{print $2}' | tr -d '"' | sed 's/<REPLACE_WITH_STRONG_PASSWORD>//')
fi

# Validate required configuration
if [ -z "$DB_HOST" ] || [ -z "$DB_USER" ] || [ -z "$DB_NAME" ]; then
    log_error "Missing required database configuration"
    exit 1
fi

# Create logs directory if it doesn't exist
mkdir -p "$(dirname "$LOG_FILE")"

log_info "=========================================="
log_info "Login Attempts Cleanup Started"
log_info "=========================================="
log_info "Retention period: ${RETENTION_DAYS} days"
log_info "Database: ${DB_NAME}@${DB_HOST}:${DB_PORT}"

# Calculate cutoff date
CUTOFF_DATE=$(date -u -d "${RETENTION_DAYS} days ago" '+%Y-%m-%d %H:%M:%S' 2>/dev/null || date -u -v-${RETENTION_DAYS}d '+%Y-%m-%d %H:%M:%S')
log_info "Cutoff date: ${CUTOFF_DATE}"

# Build MySQL command
MYSQL_CMD="mysql -h${DB_HOST} -P${DB_PORT} -u${DB_USER}"
if [ -n "$DB_PASS" ]; then
    MYSQL_CMD="${MYSQL_CMD} -p${DB_PASS}"
fi
MYSQL_CMD="${MYSQL_CMD} ${DB_NAME}"

# Count records to be deleted
log_info "Counting records to be deleted..."
COUNT_QUERY="SELECT COUNT(*) FROM login_attempts WHERE attempted_at < '${CUTOFF_DATE}';"
RECORD_COUNT=$(echo "$COUNT_QUERY" | $MYSQL_CMD -N 2>&1)

if [ $? -ne 0 ]; then
    log_error "Failed to count records: $RECORD_COUNT"
    exit 1
fi

log_info "Found ${RECORD_COUNT} records to delete"

if [ "$RECORD_COUNT" -eq 0 ]; then
    log_info "No records to delete. Cleanup complete."
    exit 0
fi

# Delete old records
log_info "Deleting records older than ${RETENTION_DAYS} days..."
DELETE_QUERY="DELETE FROM login_attempts WHERE attempted_at < '${CUTOFF_DATE}';"
DELETE_RESULT=$(echo "$DELETE_QUERY" | $MYSQL_CMD 2>&1)

if [ $? -ne 0 ]; then
    log_error "Failed to delete records: $DELETE_RESULT"
    exit 1
fi

log_info "Successfully deleted ${RECORD_COUNT} login attempt records"

# Optimize table to reclaim space
log_info "Optimizing table to reclaim disk space..."
OPTIMIZE_QUERY="OPTIMIZE TABLE login_attempts;"
OPTIMIZE_RESULT=$(echo "$OPTIMIZE_QUERY" | $MYSQL_CMD 2>&1)

if [ $? -ne 0 ]; then
    log_warn "Table optimization failed (non-critical): $OPTIMIZE_RESULT"
else
    log_info "Table optimization completed"
fi

# Get current table statistics
log_info "Fetching table statistics..."
STATS_QUERY="SELECT COUNT(*) as total_records, MIN(attempted_at) as oldest_record, MAX(attempted_at) as newest_record FROM login_attempts;"
STATS_RESULT=$(echo "$STATS_QUERY" | $MYSQL_CMD 2>&1)

if [ $? -eq 0 ]; then
    log_info "Current table statistics:"
    echo "$STATS_RESULT" | tee -a "$LOG_FILE"
fi

log_info "=========================================="
log_info "Login Attempts Cleanup Completed"
log_info "=========================================="

exit 0
