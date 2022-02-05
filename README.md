# ddns

![CI](https://github.com/milgradesec/ddns/workflows/CI/badge.svg)
[![Docker CI](https://github.com/milgradesec/ddns/actions/workflows/docker-ci.yml/badge.svg)](https://github.com/milgradesec/ddns/actions/workflows/docker-ci.yml)
[![CodeQL](https://github.com/milgradesec/ddns/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/milgradesec/ddns/actions/workflows/codeql-analysis.yml)
[![codecov](https://codecov.io/gh/milgradesec/ddns/branch/master/graph/badge.svg)](https://codecov.io/gh/milgradesec/ddns)
[![Go Report Card](https://goreportcard.com/badge/milgradesec/ddns)](https://goreportcard.com/badge/github.com/milgradesec/ddns)
![Latest Release](https://img.shields.io/github/v/release/milgradesec/ddns)
[![Go Reference](https://pkg.go.dev/badge/github.com/milgradesec/ddns.svg)](https://pkg.go.dev/github.com/milgradesec/ddns)
![GitHub](https://img.shields.io/github/license/milgradesec/ddns)

## Dynamic DNS for Cloudflare

`ddns` is a dynamic DNS client for domains managed by Cloudflare.

## Usage

### 📜 CLI Reference

```shell
Usage: ddns [options]
Options:
  -service string
        Manage DDNS as a system service
  -version
        Show version information.
  -help
        Show help.
```

<!-- Configuration example:

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
``` -->

### 🐋 Docker

`docker-compose.yaml` example:

```yaml
version: "3.8"

secrets:
  api_token:
    file: API_TOKEN

services:
  ddns:
    image: ghcr.io/milgradesec/ddns:latest
    environment:
      # Set DDNS_PROVIDER
      - DDNS_PROVIDER=Cloudflare
      # Set DDNS_ZONE
      - DDNS_ZONE=example.com
      # Set the API Token
      - CLOUDFLARE_API_TOKEN=XXXXXXXXXXXXXXXXXXXXXXXXXXX
      # Or use a docker secret
      - CLOUDFLARE_API_TOKEN_FILE=/run/secrets/api_token
    secrets:
      - api_token
    deploy:
      restart_policy:
        delay: 5s
```

## 📜 License

MIT License
