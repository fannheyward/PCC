#!/bin/sh

echo 'stop redis!'
redis-cli -p 6079 shutdown

rm -rf /tmp/redis_data
mkdir -p /tmp/redis_data

redis-server conf/redis.conf
echo "start redis!"
