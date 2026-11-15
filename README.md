<div align="center">

# the Beast

Stress testing for HTTP APIs, from the command line.

[![CI](https://github.com/jjmrocha/beast/actions/workflows/ci.yml/badge.svg)](https://github.com/jjmrocha/beast/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/jjmrocha/beast)](https://github.com/jjmrocha/beast/releases)
[![Go](https://img.shields.io/github/go-mod/go-version/jjmrocha/beast)](https://go.dev/dl/)
[![License](https://img.shields.io/github/license/jjmrocha/beast)](LICENSE)

<img src="docs/demo.gif" alt="beast template writes get.yaml for GET http://127.0.0.1:8080/apps/101, then beast run -n 20000 -c 10 sends 20,000 requests with a progress bar and prints the report: requests per second, average response time and the response time distribution" width="720">

</div>

beast sends many HTTP requests at once to an API and reports how it held up: requests per second, response time percentiles, and every status code and error. You describe the request once, in a small YAML or JSON template. A CSV file can fill in the template, so each request is different.

## Features

- Runs for a number of requests (`-n`) or a number of seconds (`-t`), with `-c` requests in progress at the same time.
- Requests come from a YAML or JSON template, and a CSV file can fill in the URL, the headers and the body.
- The report has requests per second, response time percentiles from 20% to 99%, a count for each status, and every timeout and connection error.
- `-output` writes one CSV line for each request, with its time, status and duration, for your own analysis.
- Ctrl-C stops a run early. beast waits for the requests in progress and still prints the report.
- A config file sets the timeout, keep-alives, connection limits, redirects and the TLS certificate check.
- One binary for macOS, Linux and Windows, on Intel and ARM.

## Installation

### Release binaries

Download the archive for your system from [Releases](https://github.com/jjmrocha/beast/releases), and put `beast` in a folder on your `PATH`. Each release also has `checksums.txt` and a cosign signature, to [check a download](docs/manual.md#check-a-download).

### With Go

```
go install github.com/jjmrocha/beast/cmd/beast@master
```

This needs Go 1.27 or later. Use `@master`, because `@latest` resolves to an old version without `cmd/beast`.

### From source

Clone the repository and run `make build`, which writes `bin/beast`:

```
git clone https://github.com/jjmrocha/beast.git
cd beast
make build
```

## Usage

Create a template for the endpoint you want to test:

```
beast template http://127.0.0.1:8080/apps/101 get.yaml
```

Send 1,000 requests, 10 at a time:

```
beast run -n 1000 -c 10 get.yaml
```

Run for 10 seconds, and take the app ID in the URL from a [data file](docs/manual.md#data-files):

```
beast run -t 10 -c 10 -data apps.csv apps_get.yaml
```

Write one line for each request to an [output file](docs/manual.md#output-file):

```
beast run -n 1000 -c 10 -data apps.csv -output results.csv apps_get.yaml
```

The `apps.csv` and `apps_get.yaml` files are from the manual's [worked example](docs/manual.md#worked-example).

## Documentation

The [manual](docs/manual.md) covers:

- [every command](docs/manual.md#commands), with the flags and exit statuses of `beast run`;
- [templates](docs/manual.md#templates), including a body read from a file;
- [what each line of the report means](docs/manual.md#the-report), and how much memory a long run needs;
- [configuration](docs/manual.md#configuration), the [output file](docs/manual.md#output-file) and [troubleshooting](docs/manual.md#troubleshooting).

## Contributing

Issues and pull requests are welcome. Before you open a pull request, run:

```
make test
make lint
make vuln
```

The [Development](docs/manual.md#development) section of the manual has the other make targets and the release steps.

## License

[Apache License 2.0](LICENSE) © 2019 Joaquim Rocha and Contributors
