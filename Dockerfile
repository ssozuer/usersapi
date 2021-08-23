FROM golang:latest AS buildContainer
WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -mod mod -ldflags "-s -w" -o restapi .


FROM alpine:latest
WORKDIR /app
COPY --from=buildContainer /go/src/app/restapi .

ENV GIN_MODE release
ENV JWT_SECRET "VERY_SECRET_KEY"
ENV MONGO_URI "mongodb://admin:password@localhost:27017/test?authSource=admin&readPreference=primary&ssl=false"
ENV MONGO_DATABASE "users_db"
ENV USERS_COLLECTION "users"
ENV REDIS_URI "localhost:6379"
ENV REDIS_PASSWORD ""
ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080

CMD ["./restapi"]