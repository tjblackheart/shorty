# Shorty

A small link shortener app written in Go.

**Build:**

```
git clone github.com/tjblackheart/shorty
cd shorty
go get ./...
make
```

**Usage:**

* Copy .env.dist to .env: `cp .env.dist .env`
* Set the desired variables in `.env` or just export them at runtime - at least an `APP_SECRET` is needed and should be a 32 byte long string (sha256sum will do the trick nicely).
* Create an admin user: `bin/create_user`
* Create self signed certificates (the path depends on your system): `cd tls && go run /usr/lib/go/src/crypto/tls/generate_cert.go --host MYHOSTNAME` - or use letsencrypt for that and move the files there.
* Run the app: `bin/shorty`

If you want to run without TLS, use `bin/shorty -disableTLS`

## LICENSE

[GPLv3](LICENSE)
