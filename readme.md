# Featmap

Featmap is a user story mapping tool for product people to build, plan and communicate product backlogs.

![Featmap screenshot](screenshot.png)

- [Introduction](#introduction)
  - [Purpose](#purpose)
  - [Features](#features)
  - [Intended audience](#intended-audience)
  - [Motivation](#motivation)
- [Getting started](#getting-started)
- [Self hosting](#self-hosting)
  - [Requirements](#requirements)
  - [Download](#download)
  - [Configuration](#configuration)
  - [Run](#run)
  - [Upgrade](#upgrade)
  - [Building from source](#building-from-source)
- [Development](#development)
  - [Makefile targets](#makefile-targets)
  - [Testing](#testing)
- [License](#license)

## Introduction

Featmap is an open source user story mapping tool. It is built using React, TypeScript and Go.

### Purpose

Featmap was built for product people to take advantage of a technique called *user story mapping*. User story mapping, or just story mapping, is an effective tool to create, plan and communicate your product backlog. Story mapping was created by Jeff Patton and its primary utility is providing us with an overview of the entire product and how user goals are broken down into a series of tasks. Finally, it helps us to define valuable product slices (releases) and prioritize between them.

### Features

* Personas
* Markdown editing
* Discuss user stories
* Share your user story maps with external stakeholders
* User story annotations
* User story estimates with roll-ups

### Intended audience

Featmap is great for product managers, product owners or just about anyone who is building products. Featmap can also be used as a lightweight work item management system for development teams.

### Motivation

There are many user story mapping tools, however none are really focused on ease-of-use and simplicity. Featmap was built to fill that gap. We hope you will find it as useful as we found building it.

## Getting started

1. Use our hosted service at https://www.featmap.com. This is the simplest way.
2. Host it yourself by running it on your own server. See [self-hosting](#self-hosting) below.

## Self hosting

Featmap can be run on your own server.

### Requirements

* [PostgreSQL](https://www.postgresql.org/) 11+
* One of the [pre-built binaries](https://github.com/amborle/featmap/releases) for your platform

### Download

Download the Featmap binary for your platform from the [releases page](https://github.com/amborle/featmap/releases). Make it executable if needed.

### Configuration

Featmap reads configuration from environment variables. A `conf.json` file in the working directory is also supported as a fallback.

**Environment variables:**

| Variable | Description | Required |
|---|---|---|
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `JWT_SECRET` | Secret key for JWT token signing | Yes |
| `PORT` | HTTP port (default: 5000) | No |
| `APP_SITE_URL` | Public URL of the app | No |
| `ENVIRONMENT` | Set to `development` for non-HTTPS cookie support | No |
| `MODE` | Set to `hosted` if running as a hosted service | No |
| `EMAIL_FROM` | Sender address for invitation/password emails | No |
| `SMTP_SERVER` | SMTP server hostname | No |
| `SMTP_PORT` | SMTP port (default: 587) | No |
| `SMTP_USER` | SMTP username | No |
| `SMTP_PASS` | SMTP password | No |

**conf.json fallback** (optional):

```json
{
  "appSiteURL": "https://localhost:5000",
  "dbConnectionString": "postgresql://postgres:postgres@localhost:5432/featmap?sslmode=disable",
  "jwtSecret": "ChangeMeForProduction",
  "port": "5000",
  "emailFrom": "",
  "smtpServer": "",
  "smtpPort": "587",
  "smtpUser": "",
  "smtpPass": "",
  "environment": "development"
}
```

Environment variables take precedence over conf.json. Values not set in either will use defaults where applicable.

### Run

```bash
./featmap
```

Open http://localhost:5000 in your browser.

### Upgrade

Download the latest release and replace the binary. Back up your database before upgrading.

## Building from source

### Prerequisites

* Go 1.22+
* Node.js 18+
* PostgreSQL

### Quick start

```bash
git clone https://github.com/amborle/featmap.git
cd featmap
cp config/conf.json .

# Build and run
make build
./bin/featmap
```

### With Docker

```bash
cp config/.env .
docker-compose build
docker-compose up -d
```

## Development

The project uses a Makefile for common tasks.

### Makefile targets

| Target | Description |
|---|---|
| `make` | Build the Go binary (includes webapp build via go:embed) |
| `make webapp` | Build the frontend SPA |
| `make dev` | Build webapp and run the server |
| `make test` | Run Go tests |
| `make e2e` | Run Playwright end-to-end tests |
| `make vet` | Run go vet |
| `make clean` | Remove build artifacts |
| `make release` | Cross-compile for darwin/linux/windows |
| `make docker-build` | Build Docker image |
| `make docker-up` | Start docker-compose services |

### Tech stack

**Backend** (Go):
- Chi v5 router with JWT authentication (jwtauth/v5, golang-jwt/jwt/v5)
- PostgreSQL via sqlx
- Structured logging with `log/slog`
- Graceful shutdown with signal handling
- Health check endpoint at `GET /v1/health`
- Embedded migrations and webapp via `go:embed`

**Frontend** (TypeScript + React):
- React 18 with functional components and hooks
- Redux (RTK-based store)
- React Router v5
- Vite 6 for builds
- Tailwind CSS 3
- @hello-pangea/dnd for drag-and-drop
- TypeScript 6

### Testing

**Go tests** -- API integration tests + middleware unit tests + validation tests:

```bash
make test
```

**Playwright e2e tests** -- browser-based tests with mocked API:

```bash
make e2e
```

## License

See [LICENSE](LICENSE).
