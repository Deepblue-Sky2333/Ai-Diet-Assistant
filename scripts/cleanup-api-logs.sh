#!/bin/bash

# ============================================
# API Logs Cleanup Script
# ============================================
# This script deletes API log records older than 90 days
# to prevent database bloat and maintain performance.
#
# Optional: Archive logs to object storage before deletion
#
# Usage: ./cleanup-api-logs.sh [--archive]
# Options:
#   --archive    Archive logs to file before deletion
#   --help       Show this help message
#
# Recommended: Run weekly via cron
# ============================================

set -e

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIG_FILE="${PROJECT_ROOT}/configs/config.yaml"
LOG_FILE="${PROJECT_ROOT}/logs/cleanup.log"
ARCHIVE_DIR="${PROJECT_ROOT}/logs/archive"
RETENTION_DAYS=90
ARCHIVE_ENABLED=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --archive)
            ARCHIVE_ENABLED=true
            shift
            ;;
        --help)
            echo "Usage: $0 [--archive]"
            echo ""
            echo "Options:"
            echo "  --archive    Archive logs to file before deletion"
            echo "  --help       Show this help message"
            echo ""
            echo "This script deletes API log records older than ${RETENTION_DAYS} days."
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

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

log_debug() {
    log "DEBUG" "${BLUE}$@${NC}"
}

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    log_error "Configuration file not found: $CONFIG_FILE"
    exit 1
fi

# Extract database configuration from YAML
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
log_info "API Logs Cleanup Started"
log_info "=========================================="
log_info "Retention period: ${RETENTION_DAYS} days"
log_info "Database: ${DB_NAME}@${DB_HOST}:${DB_PORT}"
log_info "Archive enabled: ${ARCHIVE_ENABLED}"

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
COUNT_QUERY="SELECT COUNT(*) FROM api_logs WHERE created_at < '${CUTOFF_DATE}';"
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

# Archive logs if enabled
if [ "$ARCHIVE_ENABLED" = true ]; then
    log_info "Archiving logs before deletion..."
    
    # Create archive directory if it doesn't exist
    mkdir -p "$ARCHIVE_DIR"
    
    # Generate archive filename with timestamp
    ARCHIVE_FILE="${ARCHIVE_DIR}/api_logs_$(date '+%Y%m%d_%H%M%S').sql.gz"
    
    log_info "Archive file: ${ARCHIVE_FILE}"
    
    # Export logs to archive
    EXPORT_QUERY="SELECT * FROM api_logs WHERE created_at < '${CUTOFF_DATE}';"
    
    # Use mysqldump for better performance with large datasets
    MYSQLDUMP_CMD="mysqldump -h${DB_HOST} -P${DB_PORT} -u${DB_USER}"
    if [ -n "$DB_PASS" ]; then
        MYSQLDUMP_CMD="${MYSQLDUMP_CMD} -p${DB_PASS}"
    fi
    MYSQLDUMP_CMD="${MYSQLDUMP_CMD} ${DB_NAME} api_logs --where=\"created_at < '${CUTOFF_DATE}'\""
    
    # Execute dump and compress
    eval "$MYSQLDUMP_CMD" | gzip > "$ARCHIVE_FILE" 2>&1
    
    if [ $? -ne 0 ]; then
        log_error "Failed to archive logs"
        exit 1
    fi
    
    # Get archive file size
    ARCHIVE_SIZE=$(du -h "$ARCHIVE_FILE" | cut -f1)
    log_info "Archive created successfully (size: ${ARCHIVE_SIZE})"
    
    # Optional: Upload to object storage (S3, MinIO, etc.)
    # Uncomment and configure if you want to upload to cloud storage
    # if command -v aws &> /dev/null; then
    #     log_info "Uploading archive to S3..."
    #     aws s3 cp "$ARCHIVE_FILE" "s3://your-bucket/api-logs-archive/"
    #     if [ $? -eq 0 ]; then
    #         log_info "Archive uploaded to S3 successfully"
    #     else
    #         log_warn "Failed to upload archive to S3"
    #     fi
    # fi
fi

# Delete old records
log_info "Deleting records older than ${RETENTION_DAYS} days..."
DELETE_QUERY="DELETE FROM api_logs WHERE created_at < '${CUTOFF_DATE}';"
DELETE_RESULT=$(echo "$DELETE_QUERY" | $MYSQL_CMD 2>&1)

if [ $? -ne 0 ]; then
    log_error "Failed to delete records: $DELETE_RESULT"
    exit 1
fi

log_info "Successfully deleted ${RECORD_COUNT} API log records"

# Optimize table to reclaim space
log_info "Optimizing table to reclaim disk space..."
OPTIMIZE_QUERY="OPTIMIZE TABLE api_logs;"
OPTIMIZE_RESULT=$(echo "$OPTIMIZE_QUERY" | $MYSQL_CMD 2>&1)

if [ $? -ne 0 ]; then
    log_warn "Table optimization failed (non-critical): $OPTIMIZE_RESULT"
else
    log_info "Table optimization completed"
fi

# Get current table statistics
log_info "Fetching table statistics..."
STATS_QUERY="SELECT COUNT(*) as total_records, MIN(created_at) as oldest_record, MAX(created_at) as newest_record FROM api_logs;"
STATS_RESULT=$(echo "$STATS_QUERY" | $MYSQL_CMD 2>&1)

if [ $? -eq 0 ]; then
    log_info "Current table statistics:"
    echo "$STATS_RESULT" | tee -a "$LOG_FILE"
fi

# Clean up old archives (keep last 12 months)
if [ "$ARCHIVE_ENABLED" = true ]; then
    log_info "Cleaning up old archives (keeping last 12 months)..."
    find "$ARCHIVE_DIR" -name "api_logs_*.sql.gz" -type f -mtime +365 -delete 2>/dev/null || true
    ARCHIVE_COUNT=$(find "$ARCHIVE_DIR" -name "api_logs_*.sql.gz" -type f | wc -l)
    log_info "Current archive count: ${ARCHIVE_COUNT}"
fi

log_info "=========================================="
log_info "API Logs Cleanup Completed"
log_info "=========================================="

exit 0
