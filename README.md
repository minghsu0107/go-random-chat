# Go Random Chat (Kafka version)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/minghsu0107/go-random-chat?label=Version&sort=semver)

Modern real-time random chat with high performance and linear scalability, written in go.

## Screenshots
<img src="https://user-images.githubusercontent.com/50090692/202243227-022dfe85-c36c-49d0-a46d-7db1d2bae16f.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/202243227-022dfe85-c36c-49d0-a46d-7db1d2bae16f.png" width="50%" height="50%" />

<img src="https://i.imgur.com/4ctofQv.png" alt="" data-canonical-src="https://i.imgur.com/4ctofQv.png" width="40%" height="40%" />

<img src="https://user-images.githubusercontent.com/50090692/157266585-90082195-0517-47a2-a1ef-20d72fa3a3e6.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/157266585-90082195-0517-47a2-a1ef-20d72fa3a3e6.png" width="40%" height="40%" />

<img src="https://user-images.githubusercontent.com/50090692/156815192-11a251fb-32ee-4888-b79c-aa64c97b407d.png" alt="" data-canonical-src="https://user-images.githubusercontent.com/50090692/156815192-11a251fb-32ee-4888-b79c-aa64c97b407d.png" width="40%" height="40%" />

## Overview

### Features
- Real-time communication and efficient websocket handling using [Melody](https://github.com/olahol/melody).
- Microservices architecture. All services **are stateless** and can be horizontally scaled on demand.
  - `web`: frontend server
  - `user`: user account server
  - `match`: user matching server
  - `chat`: messaging server
  - `uploader`: file uploader
- Use gRPC for inter-service communication
  - with retry, timeout, and circuit breaker
- Use [cobra](https://github.com/spf13/cobra) and [viper](https://github.com/spf13/viper) for CLI and configuration management respectively.
- Dependency injection using [wire](https://github.com/google/wire).
- Observability using [Golang Prometheus client](https://github.com/prometheus/client_golang) for monitoring and [opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go) for tracing.
- At-least-once delivery for message Pub/Sub using [Kafka](https://kafka.apache.org).
- Persist messages and chat channel metadata in [Cassandra](https://cassandra.apache.org), an open source NoSQL distributed database trusted by thousands of companies for scalability and high availability.
- Automatically generate RESTful API documentation with Swagger 2.0.
- User login session management using http-only cookie.
- Support Google OAuth2 login.
  - Display the name and picture of the logged in user's google account.
  - [OAuth2 userinfo spec](https://any-api.com/googleapis_com/oauth2/docs/userinfo/oauth2_userinfo_get).
- User matching with idempotency.
- Chat channel authentication using JWT.
- Store uploaded files in S3-compatible object storage.
- Support uploading image from clipboard.
- Use [Traefik FowardAuth](https://doc.traefik.io/traefik/middlewares/http/forwardauth/) for file upload authentication.
- Protect file upload api with distributed rate limiting (token bucket algorithm).
- Message seen feature.
- Auto-scroll to the first unseen message.
- Persist chat history on browser close or page refresh.
- Automatic websocket reconnection.
- Responsive web design.

### System architecture
In a microservice architecture, each service holds its own data within a dedicated database and will not share its database with others. Even though we used the same Redis cluster for multiple services, we did not violate this pattern as all services acess only their own Redis keys.

<img width="807" alt="image" src="https://user-images.githubusercontent.com/50090692/160285139-81fc63ad-76ef-41a7-8b33-c67f633f738d.png">

## Getting Started

Prerequisite:
- Docker-Compose v2
- Root permission

First, [create OAuth client ID credentials](https://developers.google.com/workspace/guides/create-credentials#web-application) and replace `USER_OAUTH_GOOGLE_CLIENTID` and `USER_OAUTH_GOOGLE_CLIENTSECRET` with your credentials in `run.sh`.

To run locally, execute the following command:
```bash
cd deployments
sudo ./run.sh start
```
`run.sh` needs root permission to alias `minio` to `localhost` in `/etc/hosts`.

Check cassandra connection:
```
docker exec deployments-cassandra-1 bash -c "cqlsh -u ming -p cassandrapass"
```
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
- Visit `http://localhost/api/<svc>/swagger/index.html` for API documentation, where `<svc>` could be `user`, `match`, `chat`, or `uploader`.

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
