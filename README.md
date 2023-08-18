# Go uploader package

## Run server

Default listen port is 9090

```shell
go run server/main.go
```

## Run client

Upload file

```shell
go run client/cmd/main.go -P 9090 -s http -f ~/Downloads/sequel-pro-1.1.2.dmg 
```

## Download file

Support Range Request

```shell

curl -v  http://localhost:9090/file/1694518204508062sequel-pro-1.1.2.dmg --output /dev/null
```

```shell
 % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0*   Trying 127.0.0.1:9090...
* Connected to localhost (127.0.0.1) port 9090 (#0)
> GET /file/1694518204508062sequel-pro-1.1.2.dmg HTTP/1.1
> Host: localhost:9090
> User-Agent: curl/7.87.0
> Accept: */*
>
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Accept-Ranges: bytes
< Content-Disposition: attachment; filename=
< Content-Length: 10628282
< Content-Type: application/octet-stream
< Last-Modified: Tue, 12 Sep 2023 11:30:04 GMT
< Date: Tue, 12 Sep 2023 11:31:10 GMT
<
{ [20194 bytes data]
100 10.1M  100 10.1M    0     0   435M      0 --:--:-- --:--:-- --:--:--  563M
* Connection #0 to host localhost left intact
```
