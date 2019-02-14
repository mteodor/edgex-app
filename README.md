# edgex-app
Edgex-app service provides means of catching the events incoming from the edgex gateway.
When Mainflux is running and Edgex gateway is subscribed to channel messages being sent
are not normalized since the unadequate format is used.
Normalizer will send unormalized messages to NATS on topic ```out.unknown```


## Configuration

The service is configured using the environment variables presented in the
following table. Note that any unset variables will be replaced with their
default values.

| Variable                   | Description                                                            | Default        |
|----------------------------|------------------------------------------------------------------------|----------------|
| MF_EDGEX_APP_LOG_LEVEL     | Log level for edgex-app (debug, info, warn, error)                     | error          |
| MF_EDGEX_APP_PORT          | Default application port for status querrying                          | 8000           |
| MF_EDGEX_DB_HOST           | Database host address                                                  | localhost      |
| MF_EDGEX_DB_PORT           | Database host port                                                     | 5432           |
| MF_EDGEX_DB_USER           | Database user                                                          | mainflux       |
| MF_EDGEX_DB_PASS           | Database password                                                      | mainflux       |
| MF_EDGEX_DB                | Name of the database used by the service                               | things         |
| MF_EDGEX_DB_SSL_MODE       | Database connection SSL mode (disable, require, verify-ca, verify-full)| disable        |
| MF_EDGEX_DB_SSL_CERT       | Path to the PEM encoded certificate file                               |                |
| MF_EDGEX_DB_SSL_KEY        | Path to the PEM encoded key file                                       |                |
| MF_EDGEX_DB_SSL_ROOT_CERT  | Path to the PEM encoded root certificate file                          |                |


## Starting
 
To start the service outside of the container, execute the following shell script:

```bash
# download the latest version of the service
go get github.com/mteodor/edgex-app

cd $GOPATH/src/github.com/mteodor/edgex-app

# compile the app
make 

# copy binary to bin
make install

# set the environment variables and run the service
MF_EDGEX_APP_LOG_LEVEL=[edgex-app log level] MF_EDGEX_DB_HOST=[Database host address] MF_EDGEX_DB_PORT=[Database host port] MF_EDGEX_DB_USER=[Database user] MF_EDGEX_DB_PASS=[Database password] MF_EDGEX_DB=[Name of the database used by the service] MF_EDGEX_DB_SSL_MODE=[SSL mode to connect to the database with] MF_EDGEX_DB_SSL_CERT=[Path to the PEM encoded certificate file] MF_THINGS_DB_SSL_KEY=[Path to the PEM encoded key file] MF_THINGS_DB_SSL_ROOT_CERT=[Path to the PEM encoded root certificate file]  MF_EDGEX_APP_PORT=[Running port] $GOBIN/edgex-app
```



## Usage
Application can be used to catch edgex events

```
curl  http://localhost:8000/version

{"service":"exapp","version":"0.0.1"}

curl -X POST http://localhost:8000/status -d '{"name": "John"}'

{"greeting":"Hello, John, I'm working fine"}

```
