# dyndns

[![Build](https://github.com/ViBiOh/dyndns/workflows/Build/badge.svg)](https://github.com/ViBiOh/dyndns/actions)
[![codecov](https://codecov.io/gh/ViBiOh/dyndns/branch/main/graph/badge.svg)](https://codecov.io/gh/ViBiOh/dyndns)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ViBiOh_dyndns&metric=alert_status)](https://sonarcloud.io/dashboard?id=ViBiOh_dyndns)

Create or update a root A record for your DNS zone on Cloudflare based on your public IP.

## Usage

The application can be configured by passing CLI args described below or their equivalent as environment variable. CLI values take precedence over environments variables.

Be careful when using the CLI values, if someone list the processes on the system, they will appear in plain-text. Pass secrets by environment variables: it's less easily visible.

```bash
Usage of dyndns:
  --domain            string slice  [dyndns] Domain to configure ${DYNDNS_DOMAIN}, as a string slice, environment variable separated by ","
  --entry             string        [dyndns] DNS Entry CNAME ${DYNDNS_ENTRY} (default "dyndns")
  --loggerJson                      [logger] Log format as JSON ${DYNDNS_LOGGER_JSON} (default false)
  --loggerLevel       string        [logger] Logger level ${DYNDNS_LOGGER_LEVEL} (default "INFO")
  --loggerLevelKey    string        [logger] Key for level in JSON ${DYNDNS_LOGGER_LEVEL_KEY} (default "level")
  --loggerMessageKey  string        [logger] Key for message in JSON ${DYNDNS_LOGGER_MESSAGE_KEY} (default "msg")
  --loggerTimeKey     string        [logger] Key for timestamp in JSON ${DYNDNS_LOGGER_TIME_KEY} (default "time")
  --network           string        [ip] Network ${DYNDNS_NETWORK} (default "tcp4")
  --proxied                         [dyndns] Proxied ${DYNDNS_PROXIED} (default false)
  --token             string        [dyndns] Cloudflare token ${DYNDNS_TOKEN}
  --uRL               string        [ip] URL for getting IPv4 or v6 ${DYNDNS_URL} (default "https://api64.ipify.org")
```
