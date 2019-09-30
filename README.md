the Beast
=========
*Stress testing for RESTful APIs*

## Installation
Binaries are available on [releases](https://github.com/jjmrocha/beast/releases).

Or you can download the code and compile yourself:
```sh
$ go get -u github.com/jjmrocha/beast
```

## Usage
Beast currently supports the following commands:
.help
.generate
.run

### help
Displays the help information.

```sh
$ beast help
The Beast - Stress testing for RESTful APIs
Usage:
	beast [help]
	beast generate [-m <http method>] [url] <file>
	beast run [-n <number of requests>] [-c <number of concurrent requests>] <file>
Where:
	generate Creates a request file, using user provided parameters
	         -m     string HTTP method (default "GET")
	         url    string Endpoint to be tested
	         file   string JSON file with details about the request to test

	run      Executes a script and presents a report with execution results
	         -c     int    Number of concurrent requests (default 1)
	         -n     int    Number of requests (default 1)
	         file   string JSON file with details about the request to test
```

### generate
The generate command functions as a utility to generate request files.

Can be used to generate an "empty" request:
```sh
$ beast generate test.json                                                        
File test.json was created, please edit before use

$ cat test.json
{
  "method": "Use Http method: GET/POST/PUT/DELETE",
  "url": "Http URL to be invoked",
  "headers": [
    {
      "key": "User-Agent",
      "value": "Beast/1"
    }
  ],
  "body": "Optional, enter body to send with POST or PUT"
}
```

Or, to generate the base request, that you can customize, adding headers or payload:
```sh
$ beast generate https://github.com/jjmrocha/beast test.json
File test.json was created for 'GET https://github.com/jjmrocha/beast'

$ cat test.json
{
  "method": "GET",
  "url": "https://github.com/jjmrocha/beast",
  "headers": [
    {
      "key": "User-Agent",
      "value": "Beast/1"
    }
  ]
}
```

### run
The run command loads a request from a file and executes the request concurrently multiple times.

```sh
$ beast run -n 100 -c 5 test.json
=== Test ===
Script to execeute: test.json
Number of requests: 100
Number of concurrent requests: 5
2019/10/01 00:45:20 [#...................] 5%
2019/10/01 00:45:21 [##..................] 10%
2019/10/01 00:45:22 [###.................] 15%
2019/10/01 00:45:24 [####................] 20%
2019/10/01 00:45:25 [#####...............] 25%
2019/10/01 00:45:26 [######..............] 30%
2019/10/01 00:45:28 [#######.............] 35%
2019/10/01 00:45:29 [########............] 40%
2019/10/01 00:45:30 [#########...........] 45%
2019/10/01 00:45:30 [##########..........] 50%
2019/10/01 00:45:31 [###########.........] 55%
2019/10/01 00:45:31 [############........] 60%
2019/10/01 00:45:31 [#############.......] 65%
2019/10/01 00:45:32 [##############......] 70%
2019/10/01 00:45:32 [###############.....] 75%
2019/10/01 00:45:32 [################....] 80%
2019/10/01 00:45:33 [#################...] 85%
2019/10/01 00:45:33 [##################..] 90%
2019/10/01 00:45:33 [###################.] 95%
2019/10/01 00:45:34 [####################] 100%
=== Results ===
Executed requests: 100
Time taken to complete: 1m22.767370021s
=== Stats ===
Min response time: 309.645195ms
Max response time: 2.635138135s
Avg response time: 827.6737ms
Requests per second: 6.0410
=== Status Code ===
Status 200: 40 requests
Status 429: 60 requests
```

## License
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.md).