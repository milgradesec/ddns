# ddns

![CI](https://github.com/milgradesec/ddns/workflows/CI/badge.svg)
![Docker](https://github.com/milgradesec/ddns/workflows/Docker/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/milgradesec/ddns)](https://goreportcard.com/badge/github.com/milgradesec/ddns)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/milgradesec/ddns/blob/master/LICENSE)

## Dynamic DNS for Cloudflare

`ddns` is a dynamic DNS client for domains managed by Cloudflare.

## How to Use

Configuration example:

```json
{
  "provider": "Cloudflare",
  "zone": "domain.com",
  "email": "email@domain.com",
  "apikey": "difiowehfhsahsdshndjqwh",
  "exclude": ["example1.domain.com"]
}
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
