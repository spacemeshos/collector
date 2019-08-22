# collector
collects events via node pubsub and writes them to db

## setup
docker run --name some-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres

docker run --name pgadmin -p 8080:80 -e "PGADMIN_DEFAULT_EMAIL=anton@mail.com" -e "PGADMIN_DEFAULT_PASSWORD=12341234" -d dpage/pgadmin4

to compile protobuf files run

`./scripts/genproto.sh `