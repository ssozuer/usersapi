JWT_SECRET="VERY_SECRET_KEY" \
MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin&readPreference=primary&ssl=false" \
MONGO_DATABASE="users_db" \
USERS_COLLECTION="users" \
REDIS_URI="localhost:6379" \
REDIS_PASSWORD="" \
go test