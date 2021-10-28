# ddns

[![CI](https://github.com/milgradesec/ddns/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/milgradesec/ddns/actions/workflows/golang-ci.yml)
[![Docker CI](https://github.com/milgradesec/ddns/actions/workflows/docker-ci.yml/badge.svg)](https://github.com/milgradesec/ddns/actions/workflows/docker-ci.yml)
[![Codecov](https://codecov.io/gh/milgradesec/ddns/branch/main/graph/badge.svg)](https://codecov.io/gh/milgradesec/ddns)
[![Go Report Card](https://goreportcard.com/badge/milgradesec/ddns)](https://goreportcard.com/badge/github.com/milgradesec/ddns)
![Latest Release](https://img.shields.io/github/v/release/milgradesec/ddns)
[![Go Reference](https://pkg.go.dev/badge/github.com/milgradesec/ddns.svg)](https://pkg.go.dev/github.com/milgradesec/ddns)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/milgradesec/ddns/blob/master/LICENSE)

## Dynamic DNS for Cloudflare

`ddns` is a dynamic DNS client for domains managed by Cloudflare.

## Usage

### 📜 CLI Reference

```shell
Usage: ddns [options]

Options:
  -config string
        Set configuration file. (default "config.json")
  -help
        Show help.
  -service string
        Manage DDNS as a system service
  -version
        Show version information.
```

### Example

`config.json` example:

```json
{
  "provider": "Cloudflare",
  "zone": "domain.com",
  "email": "email@domain.com",
  "apikey": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "apitoken": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
  "exclude": ["example.domain.com"],
  "interval": 3
}
```

Start `ddns` especifiying the configuration file:

```shell
ddns -config config.json
```

Run `ddns` as a system service:

```shell
ddns -service install
ddns -service start
```

### 🐋 Docker

`docker-compose.yaml` example:

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
      - CLOUDFLARE_API_TOKEN=XXXXXXXXXXXXXXXXXXXXXXXXXXX
      # Or use a docker secret
      - CLOUDFLARE_API_TOKEN_FILE=/run/secrets/api_token
    secrets:
      - api_token
    deploy:
      restart_policy:
        delay: 5s
```

### ☸️ Kubernetes

`deployment.yaml` example:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ddns
spec:
  selector:
    matchLabels:
      app: ddns
  template:
    metadata:
      labels:
        app: ddns
    spec:
      containers:
        - name: ddns
          image: ghcr.io/milgradesec/ddns:latest
          args:
            - "-config"
            - "/config/config.json"
          env:
            - name: CLOUDFLARE_API_TOKEN
              value: "XXXXXXXXXXXXXXXXXXXXXXXXXXX"
          volumeMounts:
            - name: config
              mountPath: /config
              readOnly: true
          resources:
            limits:
              cpu: "500m"
              memory: "64Mi"
      volumes:
        - name: config
          configMap:
            name: ddns-config
            items:
              - path: config.json
                key: config.json
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: ddns-config
data:
  config.json: |
    {
      "provider": "Cloudflare",
      "zone": "domain.com",
      "email": "email@domain.com",
      "apikey": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
      "apitoken": "XXXXXXXXXXXXXXXXXXXXXXXXXXX",
      "exclude": ["example.domain.com"],
      "interval": 5
    }
```
