# nginx-url-shortener

> [!WARNING]
> Turns out hosting a url shortener publicly is a really great way to get your infrastructure abused by malicious actors and then get a bunch of abuse reports on your domain so probably learn from my mistakes and don't do that.

Self-hosted URL shortener. One container.

nginx handles all the redirects, Go binary configures nginx.


```
podman compose up -d
```

Visit `http://localhost:8080`, paste a URL, get a short link back.

Config: set `ID_LENGTH` env var (default 5, ~916M combinations).
