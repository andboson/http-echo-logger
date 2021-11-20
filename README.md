# http-cli-echo-logger

A simple http-web server logging incoming requests to stdout with simple http-interface.

### Run locally

```shell
go run ./cmd/main.go
```

* Default app port - `8080`
* Default echo endpoint - `/echo`

### Run with docker

Run a container:

```shell
docker run -it -p8088:8080 andboson/http-cli-echo-logger 
```

Exec `curl` request to the `/echo` endpoint:

```shell
curl -X 'POST' -i \
  'http://localhost:8088/echo?new=1' \       
  -H 'accept: application/json' \  
  -H 'Content-Type: application/json' \
  -d '{
      "foo":"bar"   
    }'
```

This response will be returned as an answer.
Also, this request will be logged in docker console:

```shell
POST
RemoteAddr: 172.90.20.1:47694
RequestURI: /echo?new=1
 Content-Type: application/json
 Content-Length: 17
 User-Agent: curl/7.74.0
 Accept: application/json
Body:
 {
  "foo":"bar"
}
```

You can see the history of requests in your browser http://localhost:8088/

![Screenshot from 2021-11-20 19-05-33](https://user-images.githubusercontent.com/2089327/142736723-9031ae8a-45a2-4f21-9b04-57e48955bfd4.png)


