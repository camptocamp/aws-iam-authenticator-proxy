FROM golang:1.11 as builder
WORKDIR /
COPY . .
RUN go get -d -u \
    github.com/Sirupsen/logrus \
	github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags "-linkmode external -extldflags -static" -o /aws-iam-authenticator


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /aws-iam-authenticator /aws-iam-authenticator
ENTRYPOINT ["/aws-iam-authenticator"]
