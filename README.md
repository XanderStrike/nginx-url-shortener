# nginx-url-shortener

Self-hosted URL shortener. One container.

nginx handles all the redirects, Go binary configures nginx.


```
podman compose up -d
```

Visit `http://localhost:8080`, paste a URL, get a short link back.

Config: set `ID_LENGTH` env var (default 5, ~916M combinations).
