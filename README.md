

[![CodeQL](https://github.com/ihulsbus/cookbook/actions/workflows/codeql-analysis.yml/badge.svg?branch=master)](https://github.com/ihulsbus/cookbook/actions/workflows/codeql-analysis.yml)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=reliability_rating)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=ihulsbus_cookbook&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=ihulsbus_cookbook)
# cookbook

Web app to create, store and manage recipes, shopping lists and meal planner

# Getting Started
The backend server has the following dependencies:
* PostgreSQL database
* S3 compatible storage for images
* OIDC IDP (Auth0 or Dex)
* ENV variables for configuration


The following ENV variables are required:

| env key               | type   | example                       |
|-----------------------|--------|-------------------------------|
| cbb_debug             | bool   | true                          |
| cbb_oidc_url          | string | https://example.com/          |
| cbb_oidc_clientid     | string | example                       |
| cbb_database_host     | string | cookbook-database.example.com |
| cbb_database_port     | int    | 5432                          |
| cbb_database_database | string | cookbook                      |
| cbb_database_username | string |                               |
| cbb_database_password | string |                               |
| cbb_database_sslmode  | string | require                       |
| cbb_database_timezone | string | Europe/Amsterdam              |
| cbb_s3_endpoint       | string | https://s3.provider.com       |
| cbb_s3_key            | string |                               |
| cbb_s3_secret         | string |                               |
| cbb_s3_bucket         | string | cookbook                      |

# Development
The code can be run locally by providing environment variables locally. There is no requirement for working S3 or OIDC credentials. A local database can be created with Docker. 
```console
foo@bar:~/cookbook$ docker run -d --name cookbook_db \
                    -e POSTGRES_DB=<databasename> \
                    -e POSTGRES_USER=<username> \
                    -e POSTGRES_PASSWORD=<password> \
                    -p 5432:5432 \
                    postgres
```

API docs are provided as Swagger docs. These are generated from endpoint annotations.  To update the docs, make sure swag is installed:
```console
foo@bar:~/cookbook$ go install github.com/swaggo/swag/cmd/swag@latest
```

update the docs with the swag init command:
```console
foo@bar:~/cookbook$ swag init --parseDependency -g main.go
```