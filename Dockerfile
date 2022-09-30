FROM golang:alpine as builder
ENV USER=pow APP_NAME=pow USER_ID=1000

RUN adduser -D -H -u ${USER_ID} ${USER}

ADD go.mod /build/
RUN cd /build && go mod download

ADD . /build/
RUN cd /build && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -o ${APP_NAME} ./cmd

FROM scratch
ENV USER=pow APP_NAME=pow APP_DIR=/app
COPY --from=builder /build/${APP_NAME} ${APP_DIR}/
COPY --from=builder /etc/passwd /etc/passwd
WORKDIR ${APP_DIR}
USER ${USER}
ENTRYPOINT ["/app/pow"]
