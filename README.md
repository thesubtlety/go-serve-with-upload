Super simple Go file server with upload and basic authentication.

This is a pet project to play with go, use at your own risk...

1. `./gen-cert`
2. `./build.sh server.go`
3. `./build/server`

```
Usage: ./server
  -d string
        directory to serve up (default ".")
  -h string
        interface to serve on (default "0.0.0.0")
  -k    don't use TLS
  -kc string
        path to cert file (default "server.pem")
  -kf string
        path to key file (default "server.key")
  -p string
        port to serve on (default "8000")
  -u string
        user:pass for basic auth
```

