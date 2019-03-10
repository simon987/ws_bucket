[![Build Status](https://ci.simon987.net/buildStatus/icon?job=ws_bucket)](https://ci.simon987.net/job/ws_bucket/) [![CodeFactor](https://www.codefactor.io/repository/github/simon987/ws_bucket/badge)](https://www.codefactor.io/repository/github/simon987/ws_bucket)

### Environment variables

| Name | Default |
|:---|:---| 
| `WS_BUCKET_ADDR` | `0.0.0.0:3020` |
| `WS_BUCKET_WORKDIR` | `./data` |
| `WS_BUCKET_LOGLEVEL` | `trace` |
| `WS_BUCKET_CONNSTR` | `host=localhost user=ws_bucket dbname=ws_bucket password=ws_bucket sslmode=disable` |
| `WS_BUCKET_DIALECT` | `postgres` |
| `WS_BUCKET_SECRET` | `default_secret`* |

\* You should change this value!

### Running tests
```bash
export WS_BUCKET_ADDR=0.0.0.0:3021
export WS_BUCKET_WORKDIR=.

cd test/
go test
```

### Auth
Administration endpoints require HMAC_SHA256 authentication.
Request header:
```
{
    "Timestamp": <Current time (RFC1123)>
    "X-Signature": <HMAC_SHA256(BODY + TIMESTAMP)>
}
```

Upload endpoint requires a valid upload token:
```
{
    "X-Upload-Token": <token>
}
```
