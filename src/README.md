# Simple API with Golang

## Running

```bash
docker image build --tag juliocesarmidia/http-simple-api:v1.0.0 .

docker container run --rm --name http-simple-api -d -p 9000:9000 --env MESSAGE="Hello World Golang" juliocesarmidia/http-simple-api:v1.0.0

curl -X GET --url http://localhost:9000

docker container logs -f http-simple-api

docker container rm -f http-simple-api
```
