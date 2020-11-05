FROM golang:1.15 as builder
WORKDIR /
COPY . .
RUN make build


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /aws-iam-authenticator /aws-iam-authenticator
ENTRYPOINT ["/aws-iam-authenticator"]
