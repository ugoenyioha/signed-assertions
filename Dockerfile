# Stage 1 - Builder Img
# Define the building base image
FROM golang:alpine AS builder

# Set environmet variables
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    GCCGO=gccgo \
    CGO_CFLAGS="-g -O2" \
    CGO_CXXFLAGS="-g -O2" \
    CGO_FFLAGS="-g -O2" \
    CGO_LDFLAGS="-g -O2" \
    CC="gcc"

# Create and move to working directory
RUN mkdir /build
WORKDIR /build

# Copy in files to Img
COPY . .

# Download dependencies
RUN apk upgrade --update-cache --available && \
    apk add openssl
RUN apk add git zip curl wget ca-certificates
RUN apk add xxd
RUN apk add sed
RUN apk add jq
RUN apk add openssl-dev
RUN apk add build-base
RUN apk add pkgconfig

RUN go mod download

RUN go build -o main .

# Stage 2 - Application Img
# Define the running base image 
FROM alpine 

LABEL "type"="assertingwl"

### Set working directory  
RUN mkdir /build
WORKDIR /build

### Copy in built application and other files
COPY --from=builder /build /build
RUN apk add --no-cache bash

# Export necessary port
EXPOSE 8443

# Command to run when starting the container
CMD ["/build/main"]
# - or instead for debuging ... ENTRYPOINT ["tail", "-f", "/dev/null"]

