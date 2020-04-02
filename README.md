# aws-secretsmanager-proxy

Fetches secret from AWS Secrets Manager by given name and exposes it via an HTTP API without authentication.

The tool can be used as GitLab CI service for getting secrets in job containers which haven't installed the AWS SDK:

example secret value:
```json
{
    "MYSECRET": "secure"
}
```

example GitLab CI job:
```yaml
sampleJob:
  image: alpine
  variables:
    SECRET_NAME: mysecret
    AWS_REGION: eu-west-1
    AWS_ACCESS_KEY_ID: hidden
    AWS_SECRET_ACCESS_KEY: hidden
  services:
    - name: siticom/secret-proxy
      alias: secret-proxy
  script:
    # get secret as json
    - curl http://secret-proxy:8080/json > secret.json
    # get secret as environment file and load it
    - eval $(curl http://secret-proxy:8080/env)
    # get specific secret value for json key
    - MYSECRET=$(curl http://secret-proxy:8080/get?key=MYSECRET)
```
