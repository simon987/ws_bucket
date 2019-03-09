### Environment variables

| Name | Default |
|:---|:---| 
| `WS_BUCKET_ADDR` | `0.0.0.0:3020` |
| `WS_BUCKET_WORKDIR` | `./data` |
| `WS_BUCKET_LOGLEVEL` | `trace` |
| `WS_BUCKET_CONNSTR` | `host=localhost user=ws_bucket dbname=ws_bucket password=ws_bucket sslmode=disable` |
| `WS_BUCKET_DIALECT` | `postgres` |

### Running tests
```bash
cd test/
go test
```