#!/bin/sh
mkdir -p /etc/nginx/redirects
if [ -z "$(ls -A /etc/nginx/redirects 2>/dev/null)" ]; then
    echo "# placeholder" > /etc/nginx/redirects/default.conf
fi
nginx -g 'daemon on;'
exec /url-shortener
