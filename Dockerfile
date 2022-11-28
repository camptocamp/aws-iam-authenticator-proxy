FROM golang:1.19 as builder
WORKDIR /
COPY . .
RUN make build


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /aws-iam-authenticator-proxy /aws-iam-authenticator-proxy
ENTRYPOINT ["/aws-iam-authenticator-proxy"]
