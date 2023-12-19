FROM golang:alpine AS build
ADD . /go/src/DuckPolice/
ARG GOARCH=amd64
ENV GOARCH ${GOARCH}
ENV CGO_ENABLED 0
WORKDIR /go/src/DuckPolice
RUN go build .

FROM alpine
COPY --from=build /go/src/DuckPolice/DuckPolice /bin/DuckPolice
WORKDIR /data
CMD DuckPolice