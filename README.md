# redir

a request redirector

## Usage

The current `redir` implmenets the following features:

- Link shortener: shorten links under `/s` and `/r`
- Go [Vanity Import](https://golang.org/cmd/go/#hdr-Remote_import_paths): redirect domain/x to configured VCS and pkg.go.dev for API documentation
- PV/UV timeline, visitor referer, devices visualization


The [default configuration](./config.yml) is embedded into the binary.

Alternative configuration can be used to replace default config and
specified in environtment variable REDIR_CONF, for example
`REDIR_CONF=/path/to/config.yml redir -s` to run the redir server under
given configuration.

**The served alias can only be allocated by [golang.design](https://golang.design/) members.**
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

examples:
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
$ make        # build native binary
$ make run    # assume your local redis is running
$ make build  # build docker images
$ make up     # run via docker-compose
$ make down   # remove compose stuff
$ make status # view compose status
$ make clean  # cleanup
```

## Troubleshooting

### private golang.design projects `go get` failure

1. make sure you are a member of golang.design
2. add ssh public key to your account
3. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
4. add `export GOPRIVATE=golang.design/x` to your bash profile (e.g. `.zshrc`).

## License

MIT &copy; The [golang.design](https://golang.design) Initiative Authors
