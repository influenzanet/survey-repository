Draft project for central survey repository for Influenzanet platforms

## Configuration

app.toml
```toml

# Host to listen to hostname:port
host = "localhost:8080"

# Database DSN, actually only sqlite is implemented, expect file name
dsn = "sqlite://data/repostory.db"

# Path to store 
survey_path = "data"

# Users, list of available users
# username=argon hash
# You can generate hash using the 'password' command of the binary
[users] 

dummy="$argon2id$v=19$m=65536,t=4,p=1$dA80h2Dl1u20SmdOqptPHw$eWwdQ0Qn+b8goFJ8xUB0P2RKPRTNu5+pzMxFLQ5rxvA"
```

## Server

The server exposes several routes (routes are expressed as URI without the base URL)

| verb |route                             |  description                                 | Parameters  |
|------|----------------------------------|----------------------------------------------| ----------- |
| GET  | /namespaces                      | List of available namespaces                 |             |
| POST | /import/{namespace}              | Import a survey in the namespace {namespace} | platform    |
| GET  | /survey/{id}                     | Survey info about for record {id}            |             |
| GET  | /survey/{id}/data                | Survey json data for {id}                    |             |
| GET  | namespace/{namespace}/surveys    | List surveys in a namespace                  | See details |

### namespace/{namespace}/surveys

This route accept several query parameters

- `platforms` list of platform code (coma separated list) e.g. 'fr,it,de'
- `limit`  : Number of results to return
- `offset` : Nomber of results