# Go Random Chat (Lightweight Redis version)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/minghsu0107/go-random-chat?label=Version&sort=semver)

Modern real-time random chat with high performance and linear scalability, written in go.

Features:
- Real-time communication and efficient websocket handling using [Melody](https://github.com/olahol/melody).
- Microservices architecture. All services can be horizontally scaled on demand.
  - `web`: frontend server
  - `user`: user account server
  - `match`: user matching server
  - `chat`: messaging server
  - `uploader`: file uploader
- Use gRPC for inter-service communication
  - with retry, timeout, rate limiting, and circuit breaker
- Use [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for CLI and configuration management respectively.
- Dependency injection using [wire](https://github.com/google/wire).
- Observability using [client_golang](https://github.com/prometheus/client_golang) and [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go).
- Stateless messaging with the help of [Redis Pub/Sub](https://redis.io/topics/pubsub).
  - Redis Pub/Sub provides only at-most-once delivery, so there is chance of data loss during matching and chatting stages.
  - If you prefer at-least-once delivery for message Pub/Sub, please refer to [this branch](https://github.com/minghsu0107/go-random-chat/tree/kafka) where **Kafka is used as the message broker**.
- Persist messages in Redis.
- User matching with idempotency.
- Message seen feature.
- Auto-scroll to the first unseen message.
- Automatic websocket reconnection.
- File uploads using object storage.
- Responsive web design.

System architecture:

<img width="807" alt="image" src="https://user-images.githubusercontent.com/50090692/160285139-81fc63ad-76ef-41a7-8b33-c67f633f738d.png">

## Usage
Prerequisite:
- Docker-Compose v2
- Root permission

To run locally, execute the following command:
```bash
cd deployments
sudo ./run.sh start
```
`run.sh` needs root permission to alias `minio` to `localhost` in `/etc/hosts`.

This will spin up all services declared in `docker-compose.yaml`. Visit `http://localhost` and you will see the application home page.

Bucket `myfilebucket` will be created automatically on `minio` by `createbucket`. However, if `minio` is still initializing after 5 retries of `createbucket`, the bucket creation will fail. If this happens, please run the following command once `minio` is up and running:
```bash
docker restart deployments-createbucket-1
```
- Visit `http://localhost` for the application home page.
- Visit `http://localhost:8080` for Traefik dashboard.
- VIsit `http://localhost:9000` for Minio dashboard.
- Visit `http://localhost:9090` for Prometheus dashboard.
- Visit `http://localhost:16686` for Jaeger dashboard.

Example configuration: [config.example.yaml](configs/config.example.yaml).
## Deploy with SSL
A common scenario is that one deploys the application behind a reverse proxy with SSL termination. If that is your case, remember to correctly configure your proxy for websocket. For example, in Google Cloud Platform, for websocket traffic sent through a Google Cloud external HTTP(S) load balancer, the backend service timeout is interpreted as the maximum amount of time that a WebSocket connection can remain open, whether idle or not. Therefore, you may want to use a `timeoutSec` value larger than the default 30 seconds in your `BackendConfig`.
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

<img src="https://user-images.githubusercontent.com/50090692/157266585-90082195-0517-47a2-a1ef-20d72fa3a3e6.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/157266585-90082195-0517-47a2-a1ef-20d72fa3a3e6.png" width="60%" height="60%" />

<img src="https://user-images.githubusercontent.com/50090692/156815192-11a251fb-32ee-4888-b79c-aa64c97b407d.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/156815192-11a251fb-32ee-4888-b79c-aa64c97b407d.png" width="60%" height="60%" />
