#!/bin/sh
nginx -g 'daemon on;'
exec /url-shortener
