FROM golang:1.15-alpine AS builder
LABEL stage=builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -o indoquran .


# generate clean, final image for end users
FROM alpine
RUN mkdir -p /go/src/config/yaml
WORKDIR /go
COPY --from=builder /build/indoquran /go/.
COPY --from=builder  /build/src/config/yaml/. /go/src/config/yaml/
ENV ENV=staging
ENV PORT 5000
EXPOSE 5000

# executable
ENTRYPOINT [ "./indoquran" ]
# arguments that can be overridden
# CMD [ "3", "300" ]
