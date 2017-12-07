# NGinx Reverse Proxy with self-signed certificates

In this folder you can find very easy demo Nginx that could be used 
to serve certificates and redirect all the traffic to your bot backend.

In order to begin, firstly you need to generate self-signed sertificates.
Plrase follow [the tutorial](https://github.com/w32blaster/bot-tfl-next-departure/wiki) to do that.
As result, you will have three files in this directory:

* cert.pem
* key.pem
* dhparam.pem

Then, open file `default.conf` and adjust the URL of your bot backend. If you redirect to 
http tunnel, then the URL could be something like `http://xxxxxx.ngrok.io`, if you redirect
to linked docker container, use url with the link hostname, if you try to call the backend that 
is running not in a docker container (for example, if you run it using `go run main.go`), then 
you should call the Docker bridge IP which is `http://172.17.0.1:8444`. Up to you.

Then, just build and run:

```
sudo docker build -t bot-nginx .
sudo docker run -d --name bot-nginx -p 8443:8443 bot-nginx
```

Have fun! :)