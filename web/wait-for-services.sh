#!/bin/sh

./wait-for-it.sh mysql:3306 --timeout=15 --strict -- echo "MySQL is up"
./wait-for-it.sh redis-server:6379 --timeout=15 --strict -- echo "Redis is up"

exec "$@"