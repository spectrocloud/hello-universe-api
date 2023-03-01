
[![semantic-release: angular](https://img.shields.io/badge/semantic--release-angular-e10079?logo=semantic-release)](https://github.com/semantic-release/semantic-release)
# Hello Universe API

A Spectro Cloud demo application. This is the API server for the [Hello Universe](https://github.com/spectrocloud/hello-universe) app.

<p align="center">
<img src="./static/img/spectronaut.png" alt="drawing" width="400"/>
</p>

# Overview
The [Hello Universe](https://github.com/spectrocloud/hello-universe) app includes an API server that expands the capabilities of the application. The API server requires a Postgres database to store and retrieve data. Use the [Hello Universe DB](https://github.com/spectrocloud/hello-universe-db) container for a simple integration with a Postgres database.

# Endpoints

A Postman collection is available to help you explore the API. Review the [Postman collection](./tests/postman_collection.json) to get started.

# Usage

The quickest method to start the API server locally is by using the Docker image. 

```shell
docker pull ghcr.io/spectrocloud/hello-universe-api:1.0.8
docker run -p 3000:3000 ghcr.io/spectrocloud/hello-universe-api:1.0.8
```

To start the API server you must have connectivity to a postgres instance. Use [environment variables](#environment-variables) to customize the API server start parameters.

## Environment Variables

The API server accepts the following environment variables.

| Variable    | Description                                        | Default   |
|-------------|----------------------------------------------------|-----------|
| `PORT`        | The port number the application will listen on.    | `3000`      |
| `HOST`        | The host value name the API server will listen on. | `0.0.0.0` |
| `DB_NAME`     | The database name.                                 | `counter`   |
| `DB_USER`     | The database user name to use for queries.         | `postgres`  |
| `DB_HOST`     | The hostname or url to the database.               | `0.0.0.0` |
| `DB_PASSWORD` | The database password.                             | `password`  |
| `DB_ENCRYPTION`| The Postgres [ssl mode](https://www.postgresql.org/docs/current/libpq-ssl.html) behavior to enable. Allowed values are: `require`, `verify-full`, `verify-ca`, or `disable` |`disable`|
| `DB_INIT`     | Set to `true` if you want the API server to create the required database schema and tables in the target database.| `false` |
| `AUTHORIZATION`     | Set to `true` if you want the API server to require authorization tokens in the request.| `false` |


## Authorization

The API can be enabled with authoriation wich results in all request requiring an authorization header with a token. An anonymous token is available:

```shell
931A3B02-8DCC-543F-A1B2-69423D1A0B94
```

To enable authorization for the API set the environment variable `AUTHORIZATION` to `true`.
Ensure all API requests have an `Authorization` header with the Bearer token.

```shell
curl --location --request POST 'http://localhost:3000/api/v1/counter' \
--header 'Authorization: Bearer 931A3B02-8DCC-543F-A1B2-69423D1A0B94'
```