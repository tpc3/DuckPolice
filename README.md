# DuckPolice

[![Go Report Card](https://goreportcard.com/badge/github.com/tpc3/duckpolice)](https://goreportcard.com/report/github.com/tpc3/duckpolice)
[![Docker Image CI](https://github.com/tpc3/DuckPolice/actions/workflows/docker-image.yml/badge.svg)](https://github.com/tpc3/DuckPolice/actions/workflows/docker-image.yml)
<!-- [![Go](https://github.com/tpc3/DuckPolice/actions/workflows/go.yml/badge.svg)](https://github.com/tpc3/DuckPolice/actions/workflows/go.yml) -->

Discord Bot to check duplicate URL.

## Use

### Docker

1. [Download config.yml](https://raw.githubusercontent.com/tpc3/DuckPolice/master/config.yml)
1. Enter your token to config.yml
1. `docker run --rm -it -v $(PWD):/data ghcr.io/tpc3/duckpolice`

#### invite hint

https://discord.com/api/oauth2/authorize?client_id=XXXXXXXXXXXXXXXXXX&permissions=275951766592&scope=bot

## Build

1. Clone this repository
1. `go build`

### required

- git
- golang
- gcc

## Contribute

Any contribute is welcome.
You can use Issue and Pull Requests.
