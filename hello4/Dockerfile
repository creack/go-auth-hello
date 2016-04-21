FROM    golang:1.6-alpine

ADD     . $GOPATH/src/meetup/hello

RUN     go install meetup/hello

EXPOSE  8080

CMD     hello
