FROM    golang:1.6

# Download dependencies
RUN     go get -d github.com/dgrijalva/jwt-go

# Generate public/private key pair
RUN     mkdir -p $GOPATH/src/meetup/hello && openssl genrsa -out private.key 2048 && openssl rsa -in private.key -pubout -out public.key

ADD     . $GOPATH/src/meetup/hello

RUN     go install meetup/hello

EXPOSE  8080

CMD     hello
