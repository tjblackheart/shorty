# Shorty

A small link shortener app written in Go using [Gorilla Mux](https://github.com/gorilla/mux), [Gorilla CSRF](https://github.com/gorilla/csrf), [SCS](https://github.com/alexedwards/scs) and [Pongo2](https://github.com/flosch/pongo2).

## Usage

* Copy .env to .env.local
* Set the desired variables in `.env.local` or just export them at runtime

## Dev setup

Run `docker-compose up` to have an auto reload dev base.

The frontend then is available at `https://localhost:3000` or whereever you set APP_HOST to.
The backend can be reached with the route `/_a/`.

## LICENSE

Copyright (c) 2019 Thomas Gensicke

[GPLv3](LICENSE)
