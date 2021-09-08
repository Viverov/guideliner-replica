#!/bin/sh

echo >&2 "Wait for server..."
result=$(curl -s -o /dev/null -w "%{http_code}" "$GUIDELINER_SERVER_HOST":"$GUIDELINER_SERVER_PORT"/api/v1/ping)
iterations=0
until [ "$result" -eq 200 ]
do
   echo >&2 "Server is unavailable."
   sleep 1 # highly recommended - if it's in your local network, it can try an awful lot pretty quick...
   result=$(curl -s -o /dev/null -w "%{http_code}" "$GUIDELINER_SERVER_HOST":"$GUIDELINER_SERVER_PORT"/api/v1/ping)
   iterations=$((iterations+1))
   if [ $iterations -ge 10 ]; then
     echo "Timeout: exit from script"
     exit 1
   fi
done
echo >&2 "Server is ready"
exit 0
