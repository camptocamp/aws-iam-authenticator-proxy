aws-iam-authenticator HTTP Proxy
================================


Amazon Services require valid accounts to be used. This proxy allows external
users to access an AWS EKS cluster without requiring access to AWS credentials.

**Disclaimer**: the proxy does not implement any form of authentication. You are
responsible for implementing whatever security measure you wish to enforce in
front of it.


### Example usage

In order to give access to an AWS EKS cluster without distribution credentials,
you can start the proxy with the necessary credentials as well as the cluster ID. For example, using Docker:

```bash
$ docker run --rm -p 8080:8080 \
             -e AWS_ACCESS_KEY_ID=<AWS_ACCESS_KEY_ID> \
             -e AWS_SECRET_ACCESS_KEY=<AWS_SECRET_ACCESS_KEY> \
             -e EKS_CLUSTER_ID=<EKS_CLUSTER_ID> \
    camptocamp/aws-iam-authenticator-proxy:latest
```

You should then be able to retrieve authentication tokens for your user at
http://localhost:8080.

You can set up your `~/.kube/config` to use the `exec` authentication mechanism
using `curl`:

```yaml
users:
- name: <cluster_name>
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: curl
      args:
        - -s
        - "http://<your_ip>:8080/"
```
