# ddns

[![Build Status](https://img.shields.io/travis/milgradesec/ddns/master.svg?label=build)](https://travis-ci.org/milgradesec/ddns)
[![Go Report Card](https://goreportcard.com/badge/milgradesec/ddns)](https://goreportcard.com/badge/github.com/milgradesec/ddns)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/milgradesec/ddns/blob/master/LICENSE)

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
        "example1.domain.com",
        "example2.domain.com"
    ]
}
~~~

Start `ddns` especifiying the configuration file:

~~~ cmd
$ ddns -config config.json
~~~

Use `ddns` as a system service:

~~~ cmd
$ ddns -service install
$ ddns -service start
~~~

Update `ddns` to the latest version:

~~~ cmd
$ ddns -update
~~~
