FROM nginx:stable-alpine

#
# Add certificate and key
#
ADD cert.pem /etc/nginx/conf.d/cert.pem
ADD key.pem /etc/nginx/conf.d/key.pem
ADD dhparam.pem /etc/nginx/dhparams.pem

ADD default.conf /etc/nginx/conf.d/default.conf
ADD nginx.conf /etc/nginx/nginx.conf


RUN apk upgrade --update && \
    apk add --no-cache curl && \
    touch /var/log/nginx/access.log && chmod 600 /etc/nginx/conf.d/key.pem


