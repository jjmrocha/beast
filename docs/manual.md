# the Beast manual

This manual is for engineers who load test HTTP APIs, and the [README](../README.md) has the overview.

- [Install](#install)
- [Commands](#commands)
- [Configuration](#configuration)
- [Templates](#templates)
- [Data files](#data-files)
- [The report](#the-report)
- [Output file](#output-file)
- [Worked example](#worked-example)
- [Troubleshooting](#troubleshooting)
- [Development](#development)

---

## Install

### Release binaries

Download the archive for your system from [Releases](https://github.com/jjmrocha/beast/releases). Each release has builds for macOS, Linux and Windows, on 64-bit Intel (`amd64`) and ARM (`arm64`):

| System | Archive |
|---|---|
| macOS | `beast_<version>_darwin_amd64.tar.gz`, `beast_<version>_darwin_arm64.tar.gz` |
| Linux | `beast_<version>_linux_amd64.tar.gz`, `beast_<version>_linux_arm64.tar.gz` |
| Windows | `beast_<version>_windows_amd64.zip`, `beast_<version>_windows_arm64.zip` |

`<version>` is the release tag without the `v`, for example `2.7.0`. Each archive holds the `beast` binary (`beast.exe` on Windows), `LICENSE` and `README.md`. Extract it and put `beast` in a folder on your `PATH`.

### Check a download

Every release also has `checksums.txt`, with the SHA-256 checksum of each archive. Put it in the same folder as the archive and run:

```
sha256sum --check --ignore-missing checksums.txt
```

On older macOS versions without `sha256sum`, use `shasum -a 256 --check --ignore-missing checksums.txt`. On Windows, compare the output of `Get-FileHash <archive>` in PowerShell with the archive's line in `checksums.txt`.

The release workflow also signs `checksums.txt` with [cosign](https://github.com/sigstore/cosign), and attaches the signature as `checksums.txt.sigstore.json`. To check that the file came from beast's own release workflow, install cosign and run this, with the tag of the release you downloaded:

```
cosign verify-blob \
  --certificate-identity 'https://github.com/jjmrocha/beast/.github/workflows/release.yml@refs/tags/v2.7.0' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  --bundle checksums.txt.sigstore.json \
  checksums.txt
```

Run the checksum check after the signature check passes.

### With Go

With Go 1.27 or later, install the latest code from the `master` branch:

```
go install github.com/jjmrocha/beast/cmd/beast@master
```

`@latest` doesn't work for this module. It resolves to v2.6.1, which has no `cmd/beast` package. A binary installed with `go install` prints `dev` as its version.

### From source

Clone the repository and run `make build`. It writes `bin/beast`, with the last git tag as its version.

---

## Commands

Run `beast` or `beast help` to print the usage text:

```
$ beast help
the Beast v2.7.0 - Stress testing for RESTful APIs

Usage:
   beast [help]
   beast version
   beast config <configFile>
   beast template [-m <http method>] [url] <templateFile>
   beast run (-n <number of requests> | -t <test duration>)
             [-c <number of concurrent requests>]
             [-config <configFile>] [-data <dataFile>]
             [-output <outputFile>] <templateFile>

Where:
   version  Prints the version of the Beast

   config   Creates a file with the default parameters to setup HTTP connections
            configFile   string Name of the file to be created

   template Creates a request template file, using user-provided parameters
            -m           string HTTP method (default "GET")
            url          string Endpoint to be tested
            templateFile string JSON/YAML file with details about the request to test

   run      Executes a script and presents a report with execution results
            -c           int    Number of concurrent requests (default 1)
            -n           int    Number of requests (can't be used with "-t")
            -t           int    Duration of the test in seconds (can't be used with "-n")
            -config      string Config file to setup HTTP client
            -data        string CSV file with data for request generation
            -output      string CSV file with detailed execution results
            templateFile string JSON/YAML file with details about the request to test

```

An unknown command prints the same text and exits with status 2.

### beast version

```
$ beast version
the Beast v2.7.0 - Stress testing for RESTful APIs
```

Release binaries print their release tag. A binary built with `make build` prints the last git tag, and any other build prints `dev`.

### beast config

`beast config <configFile>` writes the default HTTP client settings to a JSON file, which you can edit and pass to `beast run -config`:

```
$ beast config config.json
File config.json was created with default configuration
```

The keys are described in [Configuration](#configuration).

### beast template

`beast template [-m <http method>] [url] <templateFile>` writes a request template that you can edit before a run. The file is JSON when its name ends in `.json`, and YAML otherwise.

Without a URL, the file shows where each value goes:

```
$ beast template example.yaml
File example.yaml was created, please edit before use

$ cat example.yaml
method: 'Use Http method: GET/POST/PUT/DELETE'
endpoint: Http URL to be invoked
headers:
    User-Agent: Beast/1
request-body: Optional, enter body to send with POST or PUT
```

With a URL, the template is ready to run. `-m` sets the method and defaults to `GET`:

```
$ beast template http://127.0.0.1:8080/apps/101 get.yaml
File get.yaml was created for 'GET http://127.0.0.1:8080/apps/101'

$ cat get.yaml
method: GET
endpoint: http://127.0.0.1:8080/apps/101
headers:
    User-Agent: Beast/1
```

For `-m POST` or `-m PUT`, the file also gets a body field with placeholder text to replace. [Templates](#templates) has the full format.

The files that `beast config` and `beast template` create can be read only by their owner (mode 0600). An existing file is overwritten.

### beast run

`beast run` sends the request in `<templateFile>` again and again, from several workers at once, and then prints a [report](#the-report).

| Flag | Default | Meaning |
|---|---|---|
| `-n <requests>` | none | Send this number of requests. |
| `-t <seconds>` | none | Send requests for this number of seconds. |
| `-c <workers>` | 1 | Number of requests in progress at the same time. |
| `-config <file>` | defaults | HTTP client settings, see [Configuration](#configuration). |
| `-data <file>` | none | CSV file with values for the template, see [Data files](#data-files). |
| `-output <file>` | none | CSV file with one line for each request, see [Output file](#output-file). |

Use exactly one of `-n` and `-t`. Put every flag before the template file. beast reads a flag after the file name as a second file name, and prints the usage text.

```
beast run -n 1000 -c 10 get.yaml
beast run -t 60 -c 50 -config config.json -data apps.csv -output results.csv apps_get.yaml
```

With `-t`, beast stops sending at the deadline and waits for the requests in progress, so a run can last a little longer than `-t`.

#### Stopping a run

Press Ctrl-C once to stop a run early. beast sends no new request, waits for the requests in progress, prints the report and closes the output file:

```
2026/11/14 23:52:37 Stopping: waiting for requests in flight. Press Ctrl-C again to exit now.
===== Stats =====
Executed requests: 8335
```

The wait is at most the request timeout, 30 seconds by default. Press Ctrl-C a second time to exit at once, with no report. The output file can then miss its last lines.

#### Exit status

| Status | When |
|---|---|
| 0 | The run finished, or `help` or `version` ran. |
| 1 | A file could not be read or written. The error is printed, for example `Error reading template file missing.yaml: open missing.yaml: no such file or directory`. |
| 2 | The command line was wrong, for example `-n` and `-t` together. The usage text is printed. |
| 130 | The run was stopped with Ctrl-C and printed its report. |

A file error stops the run before the first request.

---

## Configuration

A config file is the JSON written by `beast config`:

```json
{
	"disable-compression": true,
	"disable-keep-alives": false,
	"max-connections": 0,
	"max-idle-connections": 0,
	"request-timeout": 30,
	"disable-certificate-check": false,
	"disable-redirects": true
}
```

| Key | Default | Meaning |
|---|---|---|
| `disable-compression` | `true` | When `true`, beast doesn't ask for gzip responses. |
| `disable-keep-alives` | `false` | When `true`, every request opens a new connection. |
| `max-connections` | `0` | Most connections to one host at a time. `0` means no limit. |
| `max-idle-connections` | `0` | Most idle connections kept open for reuse. `0` uses the `-c` value. |
| `request-timeout` | `30` | Seconds to wait for each response. `0` means no timeout. |
| `disable-certificate-check` | `false` | When `true`, beast accepts expired, self-signed or otherwise invalid TLS certificates. |
| `disable-redirects` | `true` | When `true`, a redirect is reported with its own status, such as 302, and not followed. |

A key that is missing from the file keeps its default. A negative `max-connections`, `max-idle-connections` or `request-timeout` stops the run:

```
Error reading config file negative.json: invalid config: 'request-timeout' must be zero or positive
```

A file that isn't valid JSON is only reported, and the run goes on with the defaults:

```
Invalid JSON config file broken.json: invalid character '}' looking for beginning of object key string
```

---

## Templates

A template describes the request that `beast run` sends. It can be YAML or JSON. beast reads a file whose name ends in `.json` (in any case) as JSON, and any other file as YAML.

### YAML

```yaml
method: PUT
endpoint: "http://127.0.0.1:8080/apps/{{ .Data.APP_ID }}"
headers:
    User-Agent: Beast/1
    Content-Type: application/json
request-body: >
    {
       "app-id": {{ .Data.APP_ID }},
       "name": "{{ .Data.NAME }}"
    }
```

| Field | Meaning |
|---|---|
| `method` | HTTP method, such as `GET` or `PUT`. |
| `endpoint` | The URL. |
| `headers` | Header names and values. Optional. |
| `request-body` | The body. Optional. |

Put quotes around a value that starts with `{{`, such as `"{{ .Data.X_KEY }}"`. Without them, YAML reads the braces as its own syntax. beast then prints `Invalid YAML template file`, drops that value and sends the request without it.

### JSON

```json
{
	"method": "PUT",
	"url": "http://127.0.0.1:8080/apps/{{ .Data.APP_ID }}",
	"headers": [
		{
			"key": "Content-Type",
			"value": "application/json"
		},
		{
			"key": "X-Request-Id",
			"value": "{{ .RequestID }}"
		}
	],
	"body": "@app.json"
}
```

The JSON fields are `method`, `url`, `headers` (a list of `key` and `value` pairs) and `body`. JSON allows no comma after the last item of a list or object, and a comma there leaves the whole template empty (see [Broken templates](#broken-templates)).

### Values in a template

beast fills in a template with Go's [text/template](https://pkg.go.dev/text/template) package before each request. Two values are available:

| Value | What it holds |
|---|---|
| `{{ .RequestID }}` | The number of the request, starting at 1. |
| `{{ .Data.<column> }}` | A value from the current row of the data file, see [Data files](#data-files). |

You can use them in the URL, in header values and in the body. Header names are not filled in.

### Body from a file

A body that starts with `@` names a file whose content is the body. In the JSON example above, `"body": "@app.json"` reads `app.json`:

```
{
   "app-id": {{ .Data.APP_ID }},
   "name": "{{ .Data.NAME }}"
}
```

beast reads a relative path from the template's own folder, not from the folder you run beast in. An absolute path is used as it is. This works for `body` in JSON and `request-body` in YAML, and the file's content can use template values. A missing body file stops the run:

```
Error reading template file nobody.yaml: external body: open missing_body.json: no such file or directory
```

### Broken templates

A template file that can't be read, or a body file that is missing, stops the run with status 1. So does a template with a syntax error in a value, such as `{{ .Data.APP_ID }` with a brace missing, which prints `Error compiling template`.

A file that isn't valid JSON or YAML is only reported, and the run goes on with the fields that could be read. Check the first lines of the output for `Invalid JSON template file` or `Invalid YAML template file`. A YAML error drops only the value that is wrong. A JSON syntax error leaves every field empty, so every request fails as `Unexpected error` with `unsupported protocol scheme ""`.

---

## Data files

A data file is a CSV file with values for the template. The first line names the columns:

```
APP_ID,NAME
101,App 101
102,App 102
103,App 103
```

Each request uses the next row. After the last row, beast starts again from the first, so with 20 rows and `-n 1000`, each row is used 50 times. Request 1 gets the first row, but with `-c` above 1 the requests run at the same time, so they can finish in any order.

A template reads a column as `{{ .Data.<column> }}`, for example `{{ .Data.APP_ID }}`. That form needs a column name made of letters, digits and underscores, not starting with a digit. For any other name, use `index`, and put single quotes around the YAML value:

```yaml
endpoint: 'http://127.0.0.1:8080/apps/{{ index .Data "app id" }}'
```

A value that contains a comma needs double quotes, such as `"Smith, John"`.

Two problems in a data file stop the run with status 1:

| Problem | Message |
|---|---|
| No rows after the header | `Error reading data file empty.csv: data file has no rows after the header` |
| A row with a different number of values than the header | `Error reading data file short.csv: record on line 3: wrong number of fields` |

---

## The report

`beast run` prints the report to standard output in sections:

| Section | What it shows |
|---|---|
| `System` | Operating system, processor architecture and number of CPUs. |
| `Test` | The template, data file, config file, the `-n` or `-t` limit and `-c`. |
| `Preparing` | The files beast reads before the first request. |
| `Executing` | A progress bar at every 5%. It goes to standard error, with a timestamp. |
| `Stats` | Totals for the whole run, see below. |
| `Status <code>` | One section for each 2xx status, with the count, the average time and the distribution. |
| `Non Success Status` | The count for each status that isn't 2xx, such as `Status 404: 10 requests`. |
| `Errors` | The count for each request that got no HTTP answer, see below. |
| `Output File` | Whether the output file was written, when `-output` is set. |

The `Stats` section holds four totals:

| Line | Meaning |
|---|---|
| `Executed requests` | Every request, including the ones that failed. |
| `Time taken to complete` | Time from the start of the requests to the end of the run. |
| `Requests per second` | Executed requests divided by the time taken. |
| `Avg response time` | The average time of every request, including failures. |

The `Errors` section has three kinds:

| Error | Cause |
|---|---|
| `Request timeout` | No response within `request-timeout`. |
| `Unexpected error` | The request failed in another way, for example a refused connection. beast also prints the error. |
| `Request generation error` | The template could not be filled in for that request, or the result was not a valid request. |

A response time is measured from the moment the request is sent until the status line and headers arrive. Reading the body isn't included.

### Distribution

Each `Status <code>` section lists the fastest and slowest request and a set of percentiles. A line such as `- 90% of requests under 7.072708ms` means that 90% of the requests with that status took less than that. The percentiles printed depend on the number of requests with that status:

| Requests | Lines |
|---|---|
| Fewer than 5 | Count and average only |
| 5 or more | Fastest, 20%, 40%, 60%, 80%, slowest |
| 10 or more | Adds 90% |
| 20 or more | Adds 95% |
| 100 or more | Adds 99% |

The `Status`, `Non Success Status` and `Errors` lines can come out in a different order from one run to the next.

### Memory use

beast keeps the time of every successful request in memory until the run ends, so that the percentiles are exact. Each request takes 8 bytes, so an hour at 2,000 successful requests per second needs about 58 MB. The peak is about twice that while the list grows. Plan for this before a long `-t` run at a high rate.

---

## Output file

`beast run -output results.csv` writes one line for each request to a CSV file:

```
$ head -4 results.csv
Timestamp,Request,Result,StatusCode,IsSuccess,Duration
2026-11-14 23:52:06.25,GET http://127.0.0.1:8080/apps/109,Executed,200,true,2
2026-11-14 23:52:06.25,GET http://127.0.0.1:8080/apps/114,Executed,200,true,2
2026-11-14 23:52:06.252,GET http://127.0.0.1:8080/apps/119,Executed,200,true,0
```

| Column | Meaning |
|---|---|
| `Timestamp` | When the request was sent, in local time, as `yyyy-MM-dd HH:mm:ss.fff`. Trailing zeros of the milliseconds are left out, as in `.25`. |
| `Request` | The method and the URL that were sent. |
| `Result` | `Executed` when a response arrived. Otherwise `Request timeout`, `Unexpected error` or `Request generation error`. |
| `StatusCode` | The HTTP status. Empty when no response arrived. |
| `IsSuccess` | `true` for a 2xx status, otherwise `false`. |
| `Duration` | The response time in whole milliseconds, so a request under 1 ms shows `0`. Empty for a request generation error. |

Lines are in the order the responses arrived. beast creates the file before the first request, and replaces a file that already exists. The file can be read only by its owner (mode 0600). If writing fails during the run, for example on a full disk, the run goes on and the `Output File` section reports the error.

For analysis, [jjmrocha/beast-analytics](https://github.com/jjmrocha/beast-analytics) has a starting point.

---

## Worked example

This example tests the app in [jjmrocha/load-testing-playground](https://github.com/jjmrocha/load-testing-playground). It stores apps with `PUT /apps/{app_id}` and reads them with `GET /apps/{app_id}`, on `127.0.0.1:8080`. Its README explains how to start one of its versions.

### The files

Save the three files below in one folder. The playground has the same templates in its `beast` folder, and the same data as `test_data.csv`.

`apps.csv` has 20 apps, from 101 to 120:

```
APP_ID,NAME
101,App 101
102,App 102
...
120,App 120
```

`apps_put.yaml` stores one app:

```yaml
method: PUT
endpoint: "http://127.0.0.1:8080/apps/{{ .Data.APP_ID }}"
headers:
    User-Agent: Beast/1
    Content-Type: application/json
request-body: >
    {
       "app-id": {{ .Data.APP_ID }},
       "name": "{{ .Data.NAME }}"
    }
```

`apps_get.yaml` reads one app:

```yaml
method: GET
endpoint: "http://127.0.0.1:8080/apps/{{ .Data.APP_ID }}"
headers:
    User-Agent: Beast/1
```

### Store the apps

Send one `PUT` for each row:

```
$ beast run -n 20 -data apps.csv apps_put.yaml
...
===== Status 201 =====
20 requests, with avg response time of 900.945µs
```

### Read them under load

Send 10,000 `GET` requests, 10 at a time:

```
$ beast run -n 10000 -c 10 -data apps.csv apps_get.yaml
===== System =====
Operating System: darwin
System Architecture: arm64
Logical CPUs: 8
===== Test =====
Request template: apps_get.yaml
Sample Data: apps.csv
Number of requests: 10000
Number of concurrent requests: 10
===== Preparing =====
- Loading request template
- Loading data file
===== Executing =====
2026/11/14 23:51:46 [#...................] 5%
2026/11/14 23:51:46 [##..................] 10%
...
2026/11/14 23:51:48 [###################.] 95%
2026/11/14 23:51:48 [####################] 100%
===== Stats =====
Executed requests: 10000
Time taken to complete: 2.665779209s
Requests per second: 3751.2361
Avg response time: 2.653864ms
===== Status 200 =====
10000 requests, with avg response time of 2.653864ms
And the following distribution:
- The fastest request took 284.542µs
- 20% of requests under 666.584µs
- 40% of requests under 986.125µs
- 60% of requests under 1.626541ms
- 80% of requests under 5.469958ms
- 90% of requests under 7.072708ms
- 95% of requests under 8.214042ms
- 99% of requests under 10.575458ms
- The slowest request took 18.014458ms
```

The app answered every request with 200, at about 3,750 requests per second. 60% of the requests took under 1.7 ms, but the slowest 20% took between 5.5 and 18 ms. To find where the app stops keeping up, run it again with `-c 50` or `-c 100`. Add `-output results.csv` to keep the time of each request.

---

## Troubleshooting

### Many "Unexpected error" results with a high `-c`

Every worker holds a connection, and your system limits how many files and processes one user can open. Check the limits:

```
ulimit -n
ulimit -u
```

Raise them for the current shell with `ulimit -n <limit>` and `ulimit -u <limit>`. beast prints the cause of each unexpected error to standard error.

### Requests time out

A request that gets no response within `request-timeout` seconds (30 by default) counts as `Request timeout`. Raise the value in a [config file](#configuration) if the API is slow on purpose, or set it to `0` to wait without a limit.

### A URL contains `<no value>`

A template that uses `{{ .Data.<column> }}` needs `-data`. Without a data file, or with a column name that isn't in the file, the value becomes `<no value>`. beast then sends a URL such as `/apps/%3Cno%20value%3E`. Check the column names in the header line of the data file.

### macOS doesn't open the binary

macOS blocks a binary that was downloaded in a browser and isn't signed by Apple. beast's release binaries aren't. Remove the block once, in the folder that holds `beast`:

```
xattr -d com.apple.quarantine beast
```

---

## Development

You need Go 1.27 or later. The Makefile has these targets:

| Target | What it does |
|---|---|
| `make build` | Builds `bin/beast` with the last git tag as its version. |
| `make test` | Runs every test with the race detector. |
| `make lint` | Runs golangci-lint. |
| `make vuln` | Runs govulncheck. |
| `make snapshot` | Builds every release archive into `dist/` with [GoReleaser](https://goreleaser.com), and publishes nothing. |
| `make clean` | Removes `bin/` and `dist/`. |

CI runs the tests, golangci-lint, govulncheck and a gitleaks secret scan on every push and pull request.

### Releasing

Create a release on GitHub with a new tag, such as `v2.7.0`, and write its notes. The Release workflow then builds the archives with GoReleaser, writes and signs `checksums.txt`, and attaches all of them to the release. Run `make snapshot` first to check the build on your machine.

### The README demo

`docs/demo.gif` is recorded with [VHS](https://github.com/charmbracelet/vhs) from `docs/demo.tape`. Start the playground app from the [worked example](#worked-example), then run `vhs docs/demo.tape` from the repository root.
