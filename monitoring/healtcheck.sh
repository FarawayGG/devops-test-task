#!/bin/bash

# Script for periodic application accessibility check
# Recommended for use with cronjobs or similar schedulers
# Retrieves output of HTTP request to http://localhost
# Usage:
#   - Set up a cronjob to run periodically
#   - Optionally, consider using an external monitoring service like Site24x7
# Author: [Your Name]
# Date: [Date]



app_health=$(curl -s -o /dev/null -w "%{http_code}" http://localhost)
if [[ "$app_health" == "200" ]]; then
    echo "Application is healthy (Status: $app_health): $(curl -s http://localhost)"
else
    echo "Application is unhealthy (Status: $app_health)"
fi

redis_metrics=$(redis-cli info)
if [ -n "$redis_metrics" ]; then
    echo "Redis Metrics:"
    echo "$redis_metrics"
else
    echo "Failed to retrieve Redis metrics. Check Redis connection or configuration."
fi