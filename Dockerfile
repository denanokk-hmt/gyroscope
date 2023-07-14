 
##############[Build stage]##############
FROM golang:1.19.4-alpine AS builder

# env value
ARG COMMITID
ENV COMMITID ${COMMITID}
ARG SHA_COMMIT_ID
ENV SHA_COMMIT_ID ${SHA_COMMIT_ID}

# Set golang environment
ENV GO111MODULE=on \
    GOPATH=  \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

#install SSL certification
RUN apk add --update --no-cache ca-certificates git

#worker dir
WORKDIR /app

#go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

#appli building
COPY . .
RUN CGO_ENABLED=0 go build -o bin/gyroscope cmd/main.go

##############[Run stage]##############
FROM alpine

WORKDIR /app

COPY --from=builder /app/bin/gyroscope bin/gyroscope
COPY --from=builder /app/authorization authorization

EXPOSE 8080

#CMD "bin/gyroscope" $GCP_PROJECT_ID $SERVER_CODE $APPLI_NAME $SERIES $ENV $GCS_BUCKET_NAME $GCS_BUCKET_PATH
CMD ["/app/bin/gyroscope"]