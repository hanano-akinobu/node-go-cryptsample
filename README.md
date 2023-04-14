
## sample

```shell
openssl genrsa -out private-key.pem 2048
openssl rsa -in private-key.pem -pubout -out public-key.pem

node ./client/index.js # start up node server on 3000
go run ./server/server.go # start up go server on 1323
curl -X POST http://localhost:3000 -H "Content-Type: application/json" -d '{"message": "john"}'

# then you can see the log of encrypted / decrypted message on go server

```
