# Use the offical Golang image to create a build artifact.
# https://hub.docker.com/_/golang
FROM golang:1.12-alpine3.10 as builder

RUN apk --no-cache add ca-certificates git

# make the 'build' folder the current working directory
WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

# copy project files and folders to the current working directory of container (i.e. 'build' folder)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o pbill-chat

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:latest

RUN apk update
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add wget curl
RUN apk --no-cache add tzdata

# removing apk cache
RUN rm -rf /var/cache/apk/*

# make the 'app' folder the current working directory
WORKDIR /app

# setup environment variables
ENV TZ=Asia/Ho_Chi_Minh

# Copy the binary to the production image from the builder stage.
COPY --from=builder /build/pbill-chat /app/pbill-chat
COPY --from=builder /build/player.html /app/
COPY --from=builder /build/supporter.html /app/

# Run the web service on container startup.
CMD ["./pbill-chat"]
