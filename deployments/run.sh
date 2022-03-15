#!/bin/sh
export REDIS_CLUSTER_IP=$(ifconfig | grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" \
    | grep -v 127.0.0.1 | awk '{ print $2 }' | cut -f2 -d: | head -n1)
export REDIS_PASSWORD=pass.123
export JWT_SECRET=mysecret

sudo -- sh -c -e "echo '127.0.0.1   minio' >> /etc/hosts"

case "$1" in
    "run")
        docker-compose up --scale random-chat=3;;
    "stop")
        docker-compose stop;;
    *)
        echo "command should be 'run' or 'stop'"
        exit 1;;
esac