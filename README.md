# ddns

![CI](https://github.com/milgradesec/ddns/workflows/CI/badge.svg)
![Docker](https://github.com/milgradesec/ddns/workflows/Docker/badge.svg)
[![CodeQL](https://github.com/milgradesec/ddns/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/milgradesec/ddns/actions/workflows/codeql-analysis.yml)
[![codecov](https://codecov.io/gh/milgradesec/ddns/branch/master/graph/badge.svg)](https://codecov.io/gh/milgradesec/ddns)
[![Go Report Card](https://goreportcard.com/badge/milgradesec/ddns)](https://goreportcard.com/badge/github.com/milgradesec/ddns)
![Latest Release](https://img.shields.io/github/v/release/milgradesec/ddns)
[![Go Reference](https://pkg.go.dev/badge/github.com/milgradesec/ddns.svg)](https://pkg.go.dev/github.com/milgradesec/ddns)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/milgradesec/ddns/blob/master/LICENSE)

## Dynamic DNS for Cloudflare

`ddns` is a dynamic DNS client for domains managed by Cloudflare.

## Usage

Configuration example:

```json
{
  "provider": "Cloudflare",
  "zone": "domain.com",
  "email": "email@domain.com",
  "apikey": "API_KEY",
  "apitoken": "API_TOKEN",
  "exclude": ["example.domain.com"],
  "interval": 5
}
```

Docker Compose example:

```yaml
version: "3.8"

configs:
  config.json:
    file: config.json

secrets:
  api_token:
    file: API_TOKEN

services:
  ddns:
    image: ghcr.io/milgradesec/ddns:latest
    configs:
      - source: config.json
        target: /config.json
    environment:
      # Set the API Token/Key in env
      - CLOUDFLARE_API_TOKEN=API_TOKEN
      # Or use a docker secret
      - CLOUDFLARE_API_TOKEN_FILE=/run/secrets/api_token
    secrets:
      - api_token
    deploy:
      restart_policy:
        delay: 5s
```

Start `ddns` especifiying the configuration file:

```cmd
ddns -config config.json
```

Use `ddns` as a system service:

```cmd
ddns -service install
ddns -service start
```
