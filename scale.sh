#!/bin/bash 

rm -rf server2
cp -rf server server2 
rm -f server2/.env 
touch server2/.env 
echo "PORT=4445" >> server2/.env 
echo "REDIS_URL=127.0.0.1:6379" >> server2/.env 

