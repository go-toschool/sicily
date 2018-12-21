# Sicily

Gateway micro service to route graphql queries to GRPC services.

## Graphql query

```
{
  user(id: "user_id") {
    email,
    fullname
  }
}
```

## Request server with cURL

```
$ curl -X POST \
  http://localhost:3000/users \
  -H 'Content-Type: application/graphql' \
  -d 'query {
  user(id: "user_id") {
    email,
    fullname
  }
}'
```