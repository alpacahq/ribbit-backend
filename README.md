# Ribbit Reference Implementation (Backend)

The reference implementation for the backend of a broker-dealer trading application with the Alpaca [Broker API](https://alpaca.markets/docs/broker/). The backend is implemented using Go. 

To read more about what Ribbit is, it’s use cases, how it works with Broker API, and more check out our [Ribbit documentation](https://alpaca.markets/docs/broker/ribbit/). 

You can also access the [Ribbit Reference Implementation (Android)](https://github.com/alpacahq/ribbit-android) and [Ribbit Reference Implementation (iOS)](https://github.com/alpacahq/ribbit-ios) for a reference implementation of Ribbit’s user interface for both iOS and Android.

## Caveat

This code is provided as open source for the purpose of demonstration of the Broker API usage. It is not designed for the production use and Alpaca does not offer official support of the code available in this repository.

## Development Setup
Ribbit uses golang gin as webserver, and go-pg library for connecting with a PostgreSQL database.

## High-level Architecture

![ribbit](https://user-images.githubusercontent.com/22711718/139060730-a1628b12-cf45-4d6f-ad59-0a36b055b5c5.jpeg)


## Third Party Dependencies

The application uses the following third party apps in order to enable functionality reuse:
1. Twilio Verify(Used for email/mobile verification) 
    - Twilio Verify API (https://www.twilio.com/docs/verify/api)
    - The Application requires the Twilio Account SID under the `TWILIO_ACCOUNT` variable in `env.sample`
    - The Application requires the Twilio Account Auth Token under the `TWILIO_TOKEN` variable in `env.sample`
2. Plaid (Used for quickly transferring funds from your bank to Alpaca)
    - Plaid feature is only currently supported for US customers. Integration with UK and EU customers is coming soon
    - Plaid is used for ACH transfers to Alpaca. Check Alpaca's [funding documentation](https://alpaca.markets/docs/broker/integration/funding/) for more info. Alpaca's [ACH API](https://alpaca.markets/docs/broker/api-references/funding/ach/) might be helpful as well
    - Environment variables (in `env.sample`):
      - `PLAID_CLIENT_ID`: Plaid API account ID
      - `PLAID_SECRET`: Plaid authentication token
      - `PLAID_ENV`: [sandbox|development|production]
      - `PLAID_PRODUCTS`: currently, the app only needs `auth` to authenticate ACH transfers
3. Magic Labs (Used for seemless sign up and login using a single link)
    - Environment variables (in `env.sample`):
      - `MAGIC_API_KEY`: API key provided by Magic labs
      - `MAGIC_API_SECRET`: Secret token provided by Magic labs
4. SendGrid (Used for the traditional sign up flow)
    - Environment variables (in `env.sample`):
      - `SENDGRID_API_KEY`: API key provided by SendGrid

## Get started
----

### Generating private keys
A simple and efficient way of generating private keys is through mkcert. To install it, go over to their [repo](https://github.com/FiloSottile/mkcert#readme).

After successfully installing it in your machine, run `mkcert -install`

``` bash
# allow read write execute for current user
chmod 700 ./generate-ssl.sh
```

The command will do two things:
1. Generate certificates Caddy reverse proxy namely (ribbit.com.pem — private key; ribbit-public.com.pem — public cert)
2. Generate certificates for the client side payload encryption via openssl

After running the script, the give the client side (iOS and Android apps) the public key namely `public_key.pem`

### Initializing and starting the application
-----
#### Run it with docker compose
```
docker compose up
```

#### Run it locally

_Note_: Change `POSTGRES_HOST` to `localhost` when running on local machine and not docker. The same goes for `POSTGRES_SUPERUSER_PASSWORD`, set it to empty.

```bash
# postgresql config
cp .env.sample .env
source .env

# get dependencies and run
go get -v ./...
go run ./entry/ generate_secret
# copy the cli output of the command above and replace {JWT_SECRET} with it
export JWT_SECRET={JWT_SECRET}

# create a new database based on config values in .env
go run ./entry create_db

# create our database schema
go run ./entry create_schema

# create our superadmin user, which is used to administer our API server
go run ./entry create_superadmin

# schema migration and subcommands are available in the migrate subcommand
# go run ./entry migrate [command]

# run the application
go run ./entry/main.go
```

## Tests and coverage

### Run all tests

```bash
go test -coverprofile c.out ./...
go tool cover -html=c.out

# or simply
./test.sh
```

### Run only integration tests

```bash
go test -v -run Integration ./...

./test.sh -i
```

### Run only unit tests

```bash
go test -v -short ./...

# without coverage
./test.sh -s
# with coverage
./test.sh -s -c
```

## Schema migration and cli management commands

```bash
# create a new database based on config values in .env
go run ./entry create_db

# create our database schema
go run ./entry create_schema

# create our superadmin user, which is used to administer our API server
go run ./entry create_superadmin

# schema migration and subcommands are available in the migrate subcommand
# go run ./entry migrate [command]
```
