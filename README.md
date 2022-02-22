# Go Random Chat
Fast and scalable real-time random chat written in go.

Features:
- Real-time communication and efficient websocket handling using [Melody](https://github.com/olahol/melody).
- At-least-once delivery for message fan-out with the help of [Kafka](https://kafka.apache.org).
- High performance and linear scalability using Kafka as message broker.
- User Matching with idempotency.
- Responsive web design.
## Usage
```bash
./run.sh run
```
This will spin up all services declared in `docker-compose.yaml`. Visit `localhost` and you will see the application home page.

Environment variables:
- `HTTP_PORT`: Opened port of HTTP server
- `KAFKA_ADDRS`: Kafka broker addresses
- `REDIS_PASSWORD`: Redis password
- `REDIS_ADDRS`: Redis node addresses
- `REDIS_EXPIRATION_HOURS`: The expiration of all Redis keys (in hour). Default: `24`
- `MAX_ALLOWED_CONNS`: Max allowed connections to the server. Default: `200`
- `MAX_MSGS`: Max number of messages in a channel. Default: `500`
- `JWT_SECRET`: JWT secret key
- `JWT_EXPIRATION_SECONDS`: JWT expiration seconds. Default: `86400` (24 hours)
## Screenshots
<img src="https://i.imgur.com/4ctofQv.png" alt="" data-canonical-src="https://i.imgur.com/4ctofQv.png" width="60%" height="60%" />

<img src="https://i.imgur.com/NL60zFN.png" alt="" data-canonical-src="https://i.imgur.com/NL60zFN.png" width="60%" height="60%" />
