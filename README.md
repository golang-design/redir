> [!IMPORTANT]
> **This repository is archived and no longer deployed.**
>
> `golang.design/s/` and `golang.design/x/` are now served by
> [changkun/redir](https://github.com/changkun/redir), which serves both
> `golang.design` and `changkun.de` from a single process and a single
> database.
>
> **Why.** These were two separate codebases doing the same job on the same
> host, each with its own store. This one was last updated in November
> 2021; the other kept receiving dependency updates, security fixes and
> features. Running both meant every change had to be made twice, and in
> practice it only ever got made once.
>
> **The links still work.** All 63 short links and their 95,138 recorded
> visits were migrated, with every alias verified to match this
> repository's SQLite database on both page views and unique visitors.
> Nothing was dropped and no URL changed. Each site keeps its own namespace
> of aliases and its own statistics, so `golang.design/s/x` and
> `changkun.de/s/x` remain different links.
>
> **The requests in #7 are all implemented** in the successor: it uses no
> `sqlx`, its charts show page and unique views, and they accept a date
> range.
>
> Issues and contributions belong on
> [changkun/redir](https://github.com/changkun/redir).

# redir

a request redirector, aka URL shortener

## Usage

The current `redir` implmenets the following features:

- Link shortener: shorten links under `/s`
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

If the `-a` is not provided, then redir command will throw an error.

Import from a YAML file is also possible, for instance:

```
$ redir -f import.yml
```

The aliases are either imported as a new alias or updated for an existing alias.

Moreover, it is possible to visit [`/s`](https://golang.design/s) directly listing all exist aliases under [golang.design](https://golang.design/).

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
