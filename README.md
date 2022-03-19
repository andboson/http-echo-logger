# http-cli-echo-logger

A simple http echo server for logging incoming requests

* echo server with multiple endpoints support
* mock server with custom response with request matching (as http://wiremock.org/ kinda)
* log requests to stdout
* see the history of requests in the browser 
* get the history of requests via API
* one binary file only

### Run locally

```shell
go run ./cmd/main.go
```

* Default app port - `80`
* Default echo endpoint - any, except `/`, `/ipa`
* Default api endpoint - `/ipa`

### Run with docker

Run a container:

```shell
docker run -it -p8088:80 andboson/http-cli-echo-logger 
```

or with a custom echo endpoint (`/api/v1/`):

```shell
docker run -it -p8081:80 -eCUSTOM_ENDPOINTS="[{"path":"/auth","request":"","mock":"{\"key\":\"auth_key\"}]" andboson/http-cli-echo-logger 
```

(see [docker-compose.yaml](docker-compose.yaml) to full example)

Exec `curl` request to the `/echo` endpoint:

```shell
curl -X 'POST' -i \
  'http://localhost:80/echo?new=1' \       
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

### History

You can see the history of requests in your browser http://localhost/

![Screenshot from 2021-11-20 19-05-33](https://user-images.githubusercontent.com/2089327/142736723-9031ae8a-45a2-4f21-9b04-57e48955bfd4.png)

### API

To get an array of history items

```shell
curl 'http://localhost:80/ipa' 
```

returns:
```json
[{"Body":"{\"foo2\":\"bar\"}","Header":{"Accept":["*/*"],"Content-Length":["14"],"Content-Type":["application/x-www-form-urlencoded"],"User-Agent":["curl/7.74.0"]},"Method":"POST","RemoteAddr":"172.90.20.1:58468","RequestURI":"/graphQl","URL":{"Scheme":"","Opaque":"","User":null,"Host":"","Path":"/graphQl","RawPath":"","ForceQuery":false,"RawQuery":"","Fragment":"","RawFragment":""}}]
```

