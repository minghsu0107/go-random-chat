# Go Random Chat
Modern real-time random chat written in go.

Features:
- Real-time communication and efficient websocket handling using [Melody](https://github.com/olahol/melody).
- Stateless chat servers with the help of [Redis Pub/Sub](https://redis.io/topics/pubsub).
  - Redis Pub/Sub provides only at-most-once delivery, so there is chance of data loss during matching and chatting stages.
  - If you prefer at-least-once delivery for message Pub/Sub, please refer to [this branch](https://github.com/minghsu0107/go-random-chat/tree/kafka) where Kafka is used as the message broker.
- High performance and linear scalability.
- User matching with idempotency.
- Message seen feature.
- Auto-scroll to the first unseen message.
- Automatic websocket reconnection.
- File uploads using object storage.
- Responsive web design.
## Usage
To run locally, execute the following command:
```bash
cd deployments
sudo ./run.sh run
```
`run.sh` needs root permission to alias `minio` to `localhost` in `/etc/hosts`.

This will spin up all services declared in `docker-compose.yaml`. Visit `http://localhost` and you will see the application home page.

Environment variables for the chat server:
- `HTTP_PORT`: Opened port of HTTP server
- `REDIS_PASSWORD`: Redis password
- `REDIS_ADDRS`: Redis node addresses
- `REDIS_EXPIRATION_HOURS`: The expiration of all Redis keys (in hour). Default: `24`
- `MAX_ALLOWED_CONNS`: Max allowed connections to a chat server. Default: `200`
- `MAX_MSG_SIZE_BYTE`: Max message size in byte. Default: `4096`
- `MAX_MSGS`: Max number of messages in a channel. Default: `500`
- `JWT_SECRET`: JWT secret key
- `JWT_EXPIRATION_SECONDS`: JWT expiration seconds. Default: `86400` (24 hours)
## Deploy with SSL
A common scenario is that one deploys the application behind a reverse proxy with SSL termination. If that is your case, then you should change the websocket connection from `ws` to `wss` in [chat.js](./web/assets/js/chat.js) and [home.js](./web/assets/js/home.js) respectively and rebuild Docker image. 

Also, remember to correctly configure your proxy for websocket. For example, in Google Cloud Platform, for websocket traffic sent through a Google Cloud external HTTP(S) load balancer, the backend service timeout is interpreted as the maximum amount of time that a WebSocket connection can remain open, whether idle or not. Therefore, you may want to use a `timeoutSec` value larger than the default 30 seconds in your `BackendConfig`.
## Docker Tagging Rules
| Event          | Ref                        | Docker Tags                |
| -------------- | -------------------------- | -------------------------- |
| `pull_request` | `refs/pull/2/merge`        | `pr-2`                     |
| `push`         | `refs/heads/master`        | `master`                   |
| `push`         | `refs/heads/releases/v1`   | `releases-v1`              |
| `push tag`     | `refs/tags/v1.2.3`         | `v1.2.3`, `latest`         |
| `push tag`     | `refs/tags/v2.0.8-beta.67` | `v2.0.8-beta.67`, `latest` |
## Screenshots
<img src="https://i.imgur.com/4ctofQv.png" alt="" data-canonical-src="https://i.imgur.com/4ctofQv.png" width="60%" height="60%" />

<img src="https://user-images.githubusercontent.com/50090692/156455665-1944f5b3-ce52-4465-b465-5a7e3d6f1c2a.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/156455665-1944f5b3-ce52-4465-b465-5a7e3d6f1c2a.png" width="60%" height="60%" />

<img src="https://user-images.githubusercontent.com/50090692/156457027-cd957d71-344a-46f1-8b1e-29b5eefb558f.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/156457027-cd957d71-344a-46f1-8b1e-29b5eefb558f.png" width="60%" height="60%" />
