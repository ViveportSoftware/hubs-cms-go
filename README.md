# Introduction
The backend serivce for [VIVEPORT VERSE](https://verse.viveport.com/)

## Dependency
Social and authentication provider [Mastodon](https://github.com/ViveportSoftware/mastodon)

Database management service [Directus](https://github.com/directus/directus) 

## Environment Variables
| ENVIRONMENT  VARIABLE   | DESCRIPTION                                                                                                         | EXAMPLE                                                                      |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| GO_HTTP_PORT            | Port used by this service                                                                                           | 9999                                                                         |
| LOG_LEVEL               | Set to INFO to enable more logs                                                                                     | INFO &#124; DEBUG &#124; ERROR                                               |
| ENVIRONMENT             | Set to DEVELOP to enable [Gin](https://github.com/gin-gonic/gin) logs and [Swagger](https://github.com/swaggo/swag) | PRODUCTION &#124; DEVELOP                                                    |
| MASTODON_BASE_URI       | Self hosted Mastodon URL                                                                                            | https://socialverse.viveport.com                                             |
| DIRECTUS_BASE_URI       | Directus service URL                                                                                                | [Installation guide](https://docs.directus.io/getting-started/installation/) |
| DIRECTUS_ADMIN_EMAIL    | Directus admin email                                                                                                | [Installation guide](https://docs.directus.io/getting-started/installation/) |
| DIRECTUS_ADMIN_PASSWORD | Directus admin password                                                                                             | [Installation guide](https://docs.directus.io/getting-started/installation/) |
| HUBS_BASE_URI           | Self hosted Hubs URL                                                                                                | https://verse.viveport.com                                                   |

## API
| PATH                                | METHOD | DESCRIPTION             | HEADER                 |
| ----------------------------------- | ------ | ----------------------- | ---------------------- |
| /health                             | GET    | Health check            |                        |
| /version                            | GET    | Version check           |                        |
| /api/hubs-cms/v1/events             | GET    | Get all events          | Authentication: Bearer |
| /api/hubs-cms/v1/events/:id         | GET    | Get an event            | Authentication: Bearer |
| /api/hubs-cms/v1/events/:id/liked   | POST   | Like an event           | Authentication: Bearer |
| /api/hubs-cms/v1/events/:id/unliked | POST   | Unlike an event         | Authentication: Bearer |
| /api/hubs-cms/v1/events/:id/viewed  | POST   | View an event           | Authentication: Bearer |
| /api/hubs-cms/v1/me                 | GET    | Get user profile        | Authentication: Bearer |
| /api/hubs-cms/v1/accounts/:id       | PATCH  | Update user profile     | Authentication: Bearer |
| /api/hubs-cms/v1/avatars            | GET    | Get public avatars      |                        |
| /api/hubs-cms/v1/my-avatars         | GET    | Get private avatars     | Authentication: Bearer |
| /api/hubs-cms/v1/avatars            | POST   | Create a private avatar | Authentication: Bearer |
| /api/hubs-cms/v1/avatars/:id        | DELETE | Delete a private avatar | Authentication: Bearer |
| /api/hubs-cms/v1/rooms              | GET    | Get public rooms        | Authentication: Bearer |
| /api/hubs-cms/v1/my-rooms           | GET    | Get private rooms       | Authentication: Bearer |
| /api/hubs-cms/v1/rooms/:id          | GET    | Get a room              | Authentication: Bearer |
| /api/hubs-cms/v1/rooms/:id/liked    | POST   | Like a room             | Authentication: Bearer |
| /api/hubs-cms/v1/rooms/:id/unliked  | POST   | Unlike a room           | Authentication: Bearer |
| /api/hubs-cms/v1/rooms/:id/viewed   | POST   | View a room             | Authentication: Bearer |
| /api/hubs-cms/v1/passcode/:hubsid   | POST   | Check a room's passcode |                        |

## swag
Please install swag on your build machine
https://github.com/swaggo/gin-swagger
```bash
go get -u github.com/swaggo/swag/cmd/swag
```
swagger page: http://localhost:<GO_HTTP_PORT>/swagger/index.html

## go-junit-report
Please install go-junit-report on your build machine to create test report in JUnit format
https://github.com/jstemmer/go-junit-report
```bash
go get -u github.com/jstemmer/go-junit-report
```

## gocover-cobertura
Please install gocover-cobertura on your build machine to create coverage report
https://github.com/t-yuki/gocover-cobertura
```bash
go get -u github.com/t-yuki/gocover-cobertura
```

## Install dependency
```bash
make install
```

## Format coding style
```bash
make fmt
```

## Clean-up project
```bash
make clean
```

## Build binary
```bash
make build
```

## Testing Your Application
```bash
make test
```
