Draft project for central survey repository for Influenzanet platforms

## Configuration

app.toml
```toml

# Path to store 
survey_path = "data"

[db]

# Database DSN, actually only sqlite is implemented, expect file name
dsn = "sqlite://data/repository.db"

debug = true

[server]

# Host to listen to hostname:port
host = "localhost:8080"

# Rate limiter for import
# Accept {limiter_max} request by {limiter_window} seconds
limiter_max=3
limiter_window=30

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
| POST | /import/{namespace}              | Import a survey in the namespace {namespace} | platform, survey, version  |
| GET  | /survey/{id}                     | Survey info about for record {id}            |             |
| GET  | /survey/{id}/data                | Survey json data for {id}                    |             |
| GET  | namespace/{namespace}/surveys    | List surveys in a namespace                  | See details |

### Import a survey

To import a survey you will need
- The json file of the survey
- The version of the survey (please provide the version of the survey in your influenzanet platform with the form 'yy-mm-xx')

Provided by Influenzanet team (contact us to have yours)
- Your platform code
- Credentials (user/password)

Using Curl :

```bash
SURVEY=/path/to/you/survey
VERSION=version
PLATFORM=code
curl -X POST -u "user:password" -F "platform=$PLATFORM" -F "survey=@$SURVEY" http://localhost:8080/import/influenzanet
```

Fields:

Mandatory:
- `platform`:  code of the platform in Influenzanet
- `survey`: The survey json data (Full json of survey (study 1.2, 1.3+), or survey preview)

Optional (can be mandatory if information is not provided inside the provided survey data):

- `version` : Version of the survey if not provided inside the json (if the survey is not published)
- `name` : Name of the survey, if the survey data is a preview

### namespace/{namespace}/surveys

This route accept several query parameters


- `platforms`: list of platform code (comma separated list) e.g. 'fr,it,de'
- `names`: list of surveys names
- `published_from`: timestamp, minimal published time
- `published_to`: timestamp, maximum published time

Limit the results
- `limit`  : Number of results to return
- `offset` : Start index of the first result to show (only accepted with `limit`)