# Go Random Chat
Fast and scalable real-time random chat written in go.

Features:
- Real-time communication and efficient websocket handling using [Melody](https://github.com/olahol/melody).
- Stateless chat servers with the help of [Redis Pub/Sub](https://redis.io/topics/pubsub).
  - We are not using Redis Stream here since at-most-once delivery is enough for transient messages. Thus, data loss is possible during matching and chatting stages.
  - If you want at-least-once delivery, you can consider using brokers such as Redis Stream, NATS or Kafka. The following examples show how to interact with these brokers using go client.
    - https://github.com/minghsu0107/NATS-PubSub
    - https://github.com/minghsu0107/Kafka-PubSub
    - https://github.com/minghsu0107/watermill-redistream
- High performance and linear scalability with the help of [Redis Cluster](https://redis.io/topics/cluster-spec).
- User Matching with idempotency.
- Responsive web design.
## Usage
```bash
./run.sh run
```
This will spin up all services declared in `docker-compose.yaml`. Visit `localhost` and you will see the application home page.

Environment variables:
- `HTTP_PORT`: Opened port of HTTP server
- `REDIS_PASSWORD`: Redis password
- `REDIS_ADDRS`: Redis node addresses
- `REDIS_EXPIRATION_HOURS`: The expiration of all Redis keys (in hour). Default: `24`
- `MAX_ALLOWED_CONNS`: Max allowed connections to the server. Default: `200`
- `MAX_MSGS`: Max number of messages in a channel. Default: `500`
## Screenshots
<img src="https://i.imgur.com/4ctofQv.png" alt="" data-canonical-src="https://i.imgur.com/4ctofQv.png" width="60%" height="60%" />

<img src="https://i.imgur.com/NL60zFN.png" alt="" data-canonical-src="https://i.imgur.com/NL60zFN.png" width="60%" height="60%" />
