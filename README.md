## A Simple User Service Rest API in Golang

#

A simple Rest API written in Go by using Gin Framework. The app runs on `localhost` at port `8080`.

## Api Endpoints

The exposed API endpoints are shown in below. Please refer [`Api Doc`](./swagger.yml) for more details.

#### Endpoints:

| Endpoint     | Method | Authorized Access | Description                                          |
| :----------- | :----- | :---------------- | :--------------------------------------------------- |
| `/users`     | POST   | No                | Creates new user.                                    |
| `/users`     | GET    | No                | Returns users.                                       |
| `/users/:id` | DELETE | Yes               | Deletes user for given user id.                      |
| `/login`     | POST   | No                | Logins user with `email` and `password` credentials. |

### Configuration

Before starting the api, MongoDB and Redis must be run on your local machine.

## Install MongoDB Docker Container

(Please note that the following command will not mount a local `volume` for the MongoDB container. So, the data will not be persistent. If you want to persist your data you can run the following command with -v `YOUR_LOCAL_DATA_DIR`:/data/db)

```sh
docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password  -p 27017:27017 mongo:latest
```

## Install Redis Docker Container

```sh
docker run -d --name redis -p 6379:6379 redis:latest
```

## Setting Environment Variables

If you want to start app via `go run main.go` on your local machine, the following environment variables must be set and exported in the terminal.
(Note `runApp.sh` and `runTests.sh` will be set those environment variables before running `go` command)

```sh
JWT_SECRET="VERY_SECRET_KEY"
MONGO_URI="mongodb://admin:password&localhost:27017/test?authSource=admin&readPreference=primary&ssl=false"
MONGO_DATABASE="users_db"
USERS_COLLECTION="users"
REDIS_URI="localhost:6379"
REDIS_PASSWORD=""
```

## To Run App

```sh
./runApp.sh
```

## To Run Tests

```sh
./runTests.sh
```

## To Run with Docker

```sh
docker build -t userapi .
docker run -p 8080:8080 userapi
```
