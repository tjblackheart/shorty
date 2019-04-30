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
* Create self signed certificates (the path depends on your system): `cd tls && go run /usr/lib/go/src/crypto/tls/generate_cert.go --host MYHOSTNAME`
* or: create certificate via Let's Encrypt and place them in `/tls`
* Run the app: `bin/shorty`

The frontend then is available at `https://localhost:3000` or whereever you set APP_HOST to.

The backend can be reached with the route `/_a/`.

If you want to run without TLS, use `bin/shorty -disableTLS`.

## LICENSE

Copyright (c) 2019 Thomas Gensicke
[GPLv3](LICENSE)
