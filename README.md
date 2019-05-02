# NMap Trace API

## Architechture Summary
This API uses a very light weight custom framework that supports any standard Go router. It is based of a design I leanred From Bill Kenedy and have adopted with much success in our productions systems where I currently work. It also lends for a very small footprint and less relying upon 3rd party packages. This API uses (https://github.com/julienschmidt/httprouter) well known for its strait forward use and exceptional performance. The mini framework supports middleware chaining and has logging built in. This project uses `go mod` for vendoring. This API will consume Nmap trace result files in XML format.

## Folder structure
```
trace-api
└───cmd
    └───bishopfox-api
└───internal
    └───middleware
    └───platform
    └───traces
...
```

### /cmd
The cmd folder houses all standalone binaries for the trace-api project.

### /internal
The internal folder houses all application code specific to the bishopfox-api.

### /platform
The platform folder houses application code that is used for basic functions of the api. Helper functions, middlewares, validation wrappers etc.

### /traces
The traces folder contains the buisness logic. Code inside this folder is specific to the function of the API and what the user interacts with / receives.. 

## Running the API
The trace-api folder should reside in your Go path in the folder `github.com/dbubel/bchat-api/internal/platform/web`. The API can be started and stopped with the following commands. These commands will run the API inside a docker container with port 3000 exposed. Running in containers will not persist the DB once the container is stopped using `make stop`. You also must have docker-compose installed.

```
make start
make stop
```

To test that the API is running execute the command:
```
curl http://localhost:3000/health
```
You should recieve the following back:
```
{
  "status": "ok"
}
```

Or you can run the API locally with the following command from the cmd/bishopfox-api folder.
```
make db
go run main.go
```

## Tests
From the trace-api folder run:

```
make test
```
