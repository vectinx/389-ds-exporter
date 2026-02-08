# 389-ds-exporter
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?color=blue)](https://opensource.org/licenses/MIT)
[![CI](https://img.shields.io/github/actions/workflow/status/vectinx/389-ds-exporter/ci.yml?label=ci&?branch=master)](https://github.com/vectinx/389-ds-exporter/actions/workflows/ci.yml)


A Prometheus exporter for 389-ds that collects metrics over the LDAP protocol

![Dashboard](./.res/dashboard.png)

## Features
- More than 80 metrics from the 389-ds server
- Support for 389-ds version 2.3 and higher
- Support for Berkeley DB and LMDB backends of the 389-ds server
- Minimal load on the LDAP directory thanks to connection reuse via LDAP pool
- Configuration via YAML
- Ready-to-use dashboard included

## Quick Start

### Build from Source

Requirements:

- Go >= 1.25.7
- make

```bash
# Clone project repository
git clone git@github.com:vectinx/389-ds-exporter.git
cd 389-ds-exporter

# Build the 389-ds-exporter binary
make build

# Build the docker image
make docker
```

### Run with Docker

Pull the Docker image
```bash
docker pull vectinx/389-ds-exporter
```

Prepare the configuration file according to the [documentation](docs/en/config.md). Then run the container and pass it the generated config:
```bash
docker run -d --name 389-ds-exporter \
    -v $PWD/config.yml:/etc/config.yml:ro \
    -p 9389:9389 vectinx/389-ds-exporter \
    --config /etc/config.yml
```

To test the exporter:
```bash
curl localhost:9389/metrics
```

If something goes wrong, check the logs:
```bash
docker logs 389-ds-exporter
```


## Command-Line Interface

The CLI is self-documented and available via the `-h` (`--help`) option:
```bash
usage: 389-ds-exporter [<flags>]
389 Directory Server Prometheus exporter

Flags:
  -h, --[no-]help                Show context-sensitive help (also try --help-long and --help-man).
      --[no-]config.check        Validate the current configuration and print it to stdout
      --web.metrics.path="/metrics"
                                 Path under which to expose metrics.
      --[no-]web.systemd-socket  Use systemd socket activation listeners instead of port listeners (Linux only).
      --web.listen-address=:9389 ...
                                 Addresses on which to expose metrics and web interface. Repeatable for multiple addresses. Examples: `:9100` or `[::1]:9100` for http,
                                 `vsock://:9100` for vsock
      --web.config.file=""       Path to configuration file that can enable TLS or authentication. See:
                                 https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md
      --config.file="config.yml"
                                 Path to configuration file
      --log.level=info           Only log messages with the given severity or above. One of: [debug, info, warn, error]
      --log.format=logfmt        Output format of log messages. One of: [logfmt, json]
      --[no-]version             Show application version.
```

## Example

To see the 389-ds-exporter in action, you can refer to the examples:
```bash
cd examples
docker-compose up -d
```

Then open `http://localhost:3000` in your browser and wait for the infrastructure to finish initializing.

##  Based on

This project is inspired by and **partially based on** the open-source project **[389DS‑exporter](https://github.com/ozgurcd/389DS-exporter)** by **[ozgurcd](https://github.com/ozgurcd)** (MIT Licensed). Although most of the codebase has been significantly rewritten or replaced, the original project served as an architectural and conceptual starting point. The original code remains available here:

https://github.com/ozgurcd/389DS-exporter

Please see the `LICENSE` file for details.

##  License

This project is licensed under the [MIT License](./LICENSE).