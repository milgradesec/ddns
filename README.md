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

Configuration example:

~~~ json
{
    "provider": "Cloudflare",
    "zone": "domain.com",
    "email": "email@domain.com",
    "apikey": "difiowehfhsahsdshndjqwh",
    "exclude": [
        "example.domain.com"
    ]
}
~~~

Start `ddns`

~~~ cmd
$ ddns -config config.json
~~~
