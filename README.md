# kyc-service

## Description

The service has a partial integration with several identity providers. It provides an ability to check user's uniqueness
and if so then issue a corresponding [Iden3](https://docs.iden3.io/) Verifiable Credentials via an interface to the [Issuer
service](gitlab.com/rarimo/identity/issuer)

Integrated identity provider:
- [Unstoppable domains](https://unstoppabledomains.com/blog/categories/web3-domains/article/introducing-login-with-unstoppable)
- [Civic](https://www.civic.com/)
- [WorldCoin](https://docs.worldcoin.org/)
- [GitCoin Passport](https://docs.passport.gitcoin.co/overview/introducing-gitcoin-passport)

## Install

  ```
  git clone github.com/rarimo/kyc-service-legacy
  cd kyc-service
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main migrate up
  ./main run service
  ```

## Documentation

We do use openapi:json standard for API. We use swagger for documenting our API.

To open online documentation, go to [swagger editor](http://localhost:8080/swagger-editor/) here is how you can start it
```
  cd docs
  npm install
  npm start
```
To build documentation use `npm run build` command,
that will create open-api documentation in `web_deploy` folder.

To generate resources for Go models run `./generate.sh` script in root folder.
use `./generate.sh --help` to see all available options.


## Running from docker 
  
Make sure that docker installed.

use `docker run ` with `-p 8080:80` to expose port 80 to 8080

  ```
  docker build -t github.com/rarimo/kyc-service-legacy .
  docker run -e KV_VIPER_FILE=/config.yaml github.com/rarimo/kyc-service-legacy
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command


### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).
