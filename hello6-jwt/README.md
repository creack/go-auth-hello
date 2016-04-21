# Hello6 - JWT

## Generate public / private key pair

```bash
openssl genrsa -out private.key 2048
openssl rsa -in private.key -pubout -out public.key
```

## Dependencies

```bash
go get -d go get -d github.com/dgrijalva/jwt-go
```
