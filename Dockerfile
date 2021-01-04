FROM golang:alpine AS builder
LABEL stage=builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o indoquran .


# generate clean, final image for end users
FROM alpine AS production
RUN mkdir -p /go/src/config/yaml
WORKDIR /go
COPY --from=builder /build/indoquran /go/.
COPY --from=builder  /build/src/config/yaml/. /go/src/config/yaml/
ENV ENV=development
ENV PORT=8000

# executable
ENTRYPOINT [ "./indoquran" ]
# arguments that can be overridden
# CMD [ "3", "300" ]
