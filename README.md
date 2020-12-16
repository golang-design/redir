# redir

a request redirector that is dedicated for golang.design

## Design Purpose

The current `redir` implementation talks to a redis data store for PV/UV counting,
as well as short alias storage. In the booting phase, it will read `REDIR_CONF`
from environment variable to identify configuration file (default: [`./config.yml`](./config.yml)).

`redir` is designed for the following purpose: serve two major
redirectors `/s` and `/x` (at the moment).

### 1. Redirect `golang.design/x/pkg` to the `pkg`'s actual VCS.

This is based on the `go get` vanity import path convention. With this
feature, all packages issued by [golang.design](https://golang.design) 
requires to use `golang.design/x/` import path.
That is saying, any `pkg` will be redirected to `github.com/golang-design/pkg`
if exist. The website itself will redirect the request to [pkg.go.dev](https://pkg.go.dev).

There is a reserved ping router for debugging purpose `/x/.ping` which will
give you a pong.

### 2. Redirect `golang.design/s/alias` and `golang.design/r/randstr`

The served alias can be allocated by [golang.design](https://golang.design/) members.
The current approach is to use `redir` command on the [golang.design](https://golang.design/)
server. Here is the overview of its usage:

```
$ redir
usage: redir [-s] [-f <file>] [-op <operator> -a <alias> -l <link>]
options:
  -a string
        alias for a new link
  -f string
        import aliases from a YAML file
  -l string
        actual link for the alias, optional for delete/fetch
  -op string
        operators, create/update/delete/fetch (default "create")
  -s    run redir service
example:
redir -s                  run the redir service
redir -f ./import.yml     import aliases from a file
redir -a alias -l link    allocate new short link if possible
redir -l link             allocate a random alias for the given link if possible
redir -op fetch -a alias  fetch alias information
```

For the command line usage, one only needs to use `-a`, `-l`, and `-op` if needed.
The command will talk to the Redis data store and issue a new allocated alias.
For instance, the following command:

```
$ redir -a changkun -l https://changkun.de
https://golang.design/s/changkun
```

creates a new alias under [golang.design/s/changkun](https://golang.design/s/changkun).

If the `-a` is not provided, then redir command will generate a random string as an alias, but the link can only be accessed under `/r/alias`. For instance:

```
$ redir -l https://changkun.de
https://golang.design/r/qFlKSP
```

creates a new alias under [golang.design/r/qFlKSP](https://golang.design/r/qFlKSP).

Import from a YAML file is also possible, for instance:

```
$ redir -f import.yml
```

The aliases are either imported as a new alias or updated for an existing alias.

Moreover, it is possible to visit [`/s`](https://golang.design/s) or [`/r`](https://golang.design/r) directly listing all exist aliases under [golang.design](https://golang.design/).

## Build

`Makefile` defines different ways to build the service:

```bash
$ make       # build native binary
$ make run   # assume your local redis is running
$ make build # build docker images
$ make up    # run via docker-compose
$ make down  # remove compose stuff
$ make clean # cleanup
```

## Troubleshooting

### private golang.design projects `go get` failure

1. make sure you are a member of golang.design
2. add ssh public key to your account
3. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
4. add `export GOPRIVATE=golang.design/x` to your bash profile (e.g. `.zshrc`).

## License

MIT &copy; The [golang.design](https://golang.design) Initiative Authors
