

 curl \
-X POST \
-d '{"event":{"title":"deployment started on myServer","description":"my description","category":"deployments","source":"chef","properties":[{"key":"service","value":"myService"}],"startDate":"1583836150"}}' \
-H "Content-Type:application/json" \
'http://54.173.224.180/api/v1/user-events?token=ffce170d694c831c82ee6747ace5f588'
