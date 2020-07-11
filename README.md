![Image of Engarde Gopher](https://storage.googleapis.com/gopherizeme.appspot.com/gophers/fef90366a965fc1a12ede6225907f2007ecdf5a1.png)

# Engarde - Stay on guard with Envoy Access Logs
Parse [default envoy access logs](https://www.envoyproxy.io/docs/envoy/v1.8.0/configuration/access_log#default-format)  like a champ with engarde and [jq](https://github.com/stedolan/jq)

[![GoDoc](https://godoc.org/github.com/nitishm/engarde?status.svg)](https://godoc.org/github.com/nitishm/engarde)

# Motivation
Envoy access log messages are packed with a lot of useful information but in an unstructured log format. Without prior context, or even with context it can get cumbersome to visually inspect these log messages to extract useful information.

I was inspired by this [tweet](https://twitter.com/askmeegs/status/1157029140693995521?ref_src=twsrc%5Etfw%7Ctwcamp%5Etweetembed&ref_url=https%3A%2F%2Fcdn.embedly.com%2Fwidgets%2Fmedia.html%3Ftype%3Dtext%252Fhtml%26key%3Da19fcc184b9711e1b4764040d3dc5c07%26schema%3Dtwitter%26url%3Dhttps%253A%2F%2Ftwitter.com%2Faskmeegs%2Fstatus%2F1157029140693995521%26image%3Dhttps%253A%2F%2Fi.embed.ly%2F1%2Fimage%253Furl%253Dhttps%25253A%25252F%25252Fpbs.twimg.com%25252Fmedia%25252FEA6X3jiX4AYh5X_.jpg%25253Alarge%2526key%253Da19fcc184b9711e1b4764040d3dc5c07) from [Megan O'Keefe](https://twitter.com/askmeegs) on twitter to create this tool for better readability of the envoy/istio-proxy access logs.

In addition, a special shout out to Richard Li from ambassador.io for this excellent [article](https://blog.getambassador.io/understanding-envoy-proxy-and-ambassador-http-access-logs-fee7802a2ec5) that provides more details on each of the subcomponents of the log message.

# Installing

## Homebrew

```console
brew tap nitishm/homebrew-engarde
brew install engarde
```

## Scoop

```console
scoop bucket add engarde REPO_URL
scoop install engarde
```

## Source

```console
go get github.com/nitishm/engarde
```
This should install the compiled binary to your `$GOBIN` (or `$GOPATH/bin`).

Otherwise, clone the repository and build the binary manually using,

*This package uses gomodules. Ensure that GO111MODULE=on is set if building outside `$GOPATH` in go version 1.11+*

```console
git clone https://github.com/nitishm/engarde.git
cd engarde/
go build -o engarde .
mv engarde /usr/local/bin/
export PATH=$PATH:/usr/local/bin/
``` 

# Example (default format only)
**Prerequisites**

[jq](https://github.com/stedolan/jq) must be in your `$PATH`. If you do not have `jq` installed please download the binary [here](https://stedolan.github.io/jq/)

## Envoy
```console
echo '[2016-04-15T20:17:00.310Z] "POST /api/v1/locations HTTP/2" 204 - 154 0 226 100 "10.0.35.28" "nsq2http" "cc21d9b0-cf5c-432b-8c7e-98aeb7988cd2" "locations" "tcp://10.0.2.1:80"' | engarde | jq
```
```json
2019/09/03 22:31:33 Reading input from STDIN. Use the pipe "|" operator to redirect traffic to engarde
{
  "authority": "locations",
  "bytes_received": "154",
  "bytes_sent": "0",
  "duration": "226",
  "method": "POST",
  "protocol": "HTTP/2",
  "request_id": "cc21d9b0-cf5c-432b-8c7e-98aeb7988cd2",
  "response_flags": "-",
  "status_code": "204",
  "timestamp": "2016-04-15T20:17:00.310Z",
  "upstream_service": "tcp://10.0.2.1:80",
  "upstream_service_time": "100",
  "uri_path": "/api/v1/locations",
  "user_agent": "nsq2http",
  "original_message": "[2016-04-15T20:17:00.310Z] \"POST /api/v1/locations HTTP/2\" 204 - 154 0 226 100 \"10.0.35.28\" \"nsq2http\" \"cc21d9b0-cf5c-432b-8c7e-98aeb7988cd2\" \"locations\" \"tcp://10.0.2.1:80\""
}
```

## Istio Proxy
```console
kubectl logs -f foo-app-1 -c istio-proxy | engarde --use-istio | jq
```
```json
{
  "authority": "hello-world",
  "bytes_received": "148",
  "bytes_sent": "171",
  "duration": "4",
  "method": "GET",
  "protocol": "HTTP/1.1",
  "request_id": "c0ce81db-4f5a-9134-8a5c-f8c076c91652",
  "response_flags": "-",
  "status_code": "200",
  "timestamp": "2019-09-03T05:37:41.341Z",
  "upstream_service": "192.168.89.50:9001",
  "upstream_service_time": "3",
  "upstream_cluster": "outbound|80||hello-world.default.svc.cluster.local",
  "upstream_local": "-",
  "downstream_local": "10.97.86.53:80",
  "downstream_remote": "192.168.167.113:39953",
  "uri_path": "/index",
  "user_agent": "-",
  "mixer_status": "-",
  "original_message": "[2019-09-03T05:37:41.341Z] \"GET /index HTTP/1.1\" 200 - \"-\" 148 171 4 3 \"-\" \"-\" \"c0ce81db-4f5a-9134-8a5c-f8c076c91652\" \"hello-world\" \"192.168.89.50:9001\" outbound|80||hello-world.default.svc.cluster.local - 10.97.86.53:80 192.168.167.113:39953 -"
}
```
