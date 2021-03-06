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
```
