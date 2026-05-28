FROM golang:alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go .
RUN go build -o url-shortener .

FROM nginx:alpine
RUN mkdir -p /etc/nginx/redirects
COPY nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/url-shortener /url-shortener
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]
