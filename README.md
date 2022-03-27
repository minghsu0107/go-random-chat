# Go Random Chat (Kafka version)
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
- Use [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for CLI and configuration management respectively.
- Dependency injection using [wire](https://github.com/google/wire).
- Observability using [client_golang](https://github.com/prometheus/client_golang) and [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go).
- At-least-once delivery for message Pub/Sub using [Kafka](https://kafka.apache.org).
- Message retention for a certain period.
- User matching with idempotency.
- Message seen feature.
- Auto-scroll to the first unseen message.
- Automatic websocket reconnection.
- File uploads using object storage.
- Responsive web design.

System architecture:

<img width="723" alt="image" src="https://user-images.githubusercontent.com/50090692/159188893-6036a8a7-4318-48d2-a281-3567a12071bd.png">

## Usage
Prerequisite:
- Docker-Compose v2
- root permission

To run locally, execute the following command:
```bash
cd deployments
sudo ./run.sh run
```
`run.sh` needs root permission to alias `minio` to `localhost` in `/etc/hosts`.

This will spin up all services declared in `docker-compose.yaml`. Visit `http://localhost` and you will see the application home page.

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
