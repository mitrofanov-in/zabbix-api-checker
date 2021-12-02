FROM golang:alpine

WORKDIR /app
COPY src .
ADD .env.enc .
ARG PASS_CHECKER
RUN apk add openssl
RUN openssl enc -d -aes-128-cbc -nosalt -k $PASS_CHECKER -in .env.enc -out .env
ENV GOROOT /usr/local/go
RUN ln -s /app/getenv ${GOROOT}/src/
RUN ln -s /app/auth ${GOROOT}/src/
RUN go get  github.com/go-sql-driver/mysql
RUN go get  github.com/joho/godotenv
RUN cd api-checker && go build -o check_health_zeus_api.go
EXPOSE 8080
ENTRYPOINT ["/app/api-checker/check_health_zeus_api.go"]
CMD exec /bin/sh -c "trap : TERM INT; sleep infinity & wait"
