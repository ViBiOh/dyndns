# dyndns

[![Build](https://github.com/ViBiOh/dyndns/workflows/Build/badge.svg)](https://github.com/ViBiOh/dyndns/actions)
[![codecov](https://codecov.io/gh/ViBiOh/dyndns/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/dyndns)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_dyndns&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_dyndns)

## CI

Following variables are required for CI:

|            Name            |           Purpose           |
| :------------------------: | :-------------------------: |
|      **DOCKER_USER**       | for publishing Docker image |
|      **DOCKER_PASS**       | for publishing Docker image |
| **SCRIPTS_NO_INTERACTIVE** |  for disabling bash prompt  |

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of dyndns:
  -domain string
        [dyndns] Domain to configure {DYNDNS_DOMAIN}
  -entry string
        [dyndns] DNS Entry CNAME {DYNDNS_ENTRY} (default "dyndns")
  -network string
        [ip] Network {DYNDNS_NETWORK} (default "tcp4")
  -proxied
        [dyndns] Proxied {DYNDNS_PROXIED}
  -token string
        [dyndns] Cloudflare token {DYNDNS_TOKEN}
  -uRL string
        [ip] URL for getting IPv4 or v6 {DYNDNS_URL} (default "https://api64.ipify.org")
```
