# Engarde - Stay on your guard with Envoy
Parse envoy access logs like a champ with engarde and [jq](https://github.com/stedolan/jq)

# Example

```
kubectl logs -f hello-world-1-0-0-155-dbg-85dc46fbb6-7qct6 -c istio-proxy | engarde | jq
{
  "authority": "10.10.12.12:8500",
  "bytes_received": "2",
  "bytes_sent": "0",
  "duration": "10165",
  "method": "GET",
  "protocol": "HTTP/1.1",
  "request_id": "78490242-53c8-949b-91de-6df3a1e4b09b",
  "response_flags": "- \"-\"",
  "status_code": "200",
  "timestamp": "2019-08-31T15:46:32.828Z",
  "upstream_service": "10.15.12.33:8500",
  "upstream_service_time": "10164",
  "uri_param": "?wait=10s",
  "uri_path": "/v1/foo/test",
  "user_agent": "Python-urllib/2.7",
  "original_message": "[2019-08-31T15:46:32.828Z] \"GET /v1/foo/test?wait=10s HTTP/1.1\" 200 - \"-\" 0 2 10165 10164 \"-\" \"Python-urllib/2.7\" \"78490242-53c8-949b-91de-6df3a1e4b09b\" \"10.15.12.33:8500\" \"10.15.12.33:8500\" PassthroughCluster - 10.15.12.33:8500 192.168.89.50:34106 -"
}
```
