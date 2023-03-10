FROM golang:alpine AS build
RUN apk add git gcc musl-dev
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 1
ADD . /go/src/DuckPolice/
WORKDIR /go/src/DuckPolice
RUN go build .

FROM alpine
COPY --from=build /go/src/DuckPolice/DuckPolice /bin/DuckPolice
WORKDIR /data
CMD DuckPolice