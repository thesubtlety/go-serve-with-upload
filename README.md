Super simple Go file server with upload and basic authentication. And a hardcoded SSL cert.

`./goserve`

```
Usage: ./server -h
  -d string
        directory to serve up (default ".")
  -h string
        interface to serve on (default "0.0.0.0")
  -k    don't use TLS
  -kc string
        path to cert file
  -kf string
        path to key file
  -p string
        port to serve on (default "8000")
  -u string
        user:pass for basic auth
```

