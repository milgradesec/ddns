# ddns

[![Build Status](https://img.shields.io/travis/milgradesec/ddns/master.svg?label=build)](https://travis-ci.org/milgradesec/ddns)

## Dynamic DNS for Cloudflare

`ddns` is a dynamic DNS client for domains managed by Cloudflare.

## Building

### Build from Source

~~~ cmd
$ git clone github.com/milgradesec/ddns
$ cd ddns
$ make
~~~

### Build with Docker

~~~ cmd
$ git clone github.com/milgradesec/ddns
$ cd ddns
$ make docker
~~~

## How to Use

Deploy with Docker Compose:

~~~ yaml
version: "3.7"

services:
  ddns:
    image: ddns:tag
    environment:
      - PROVIDER=Cloudflare
      - CF_API_EMAIL=your_email
      - CF_API_KEY=api_key
      - CF_ZONE_NAME=domain_name
    deploy:
      restart_policy:
        condition: on-failure
~~~
