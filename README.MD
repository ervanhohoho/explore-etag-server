# ETag Validation Sample Project

## Method Applications
### GET
When getting a resource from the server we can 
provide `If-None-Match` header to validate if 
any change has been made since the last time we 
get the resource.

If there's no change the server will return `304 Status Not Modified`
### POST/PUT/PATCH
When updating a server's resource we have to 
provide `If-Match` header taken from the GET
Request. 

If the token is different, the request
will be rejected with status `412 Precondition Failed`

## How to run
1. `go mod tidy`
2. `go run main.go`
