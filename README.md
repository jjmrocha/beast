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
* help
* config
* template
* run

### help
Displays the help information.

```sh
$ beast help
The Beast - Stress testing for RESTful APIs
Usage:
   beast [help]
   beast config <fileName>
   beast template [-m <http method>] [url] <template>
   beast run [-n <number of requests>] [-c <number of concurrent requests>] <template>
Where:
   config   Creates a file with the default parameters to setup HTTP connections
            fileName string Name of the file to be created
			 			
   template Creates a request template file, using user provided parameters
            -m       string HTTP method (default "GET")
            url      string Endpoint to be tested
            template string JSON file with details about the request to test

   run      Executes a script and presents a report with execution results
            -c       int    Number of concurrent requests (default 1)
            -n       int    Number of requests (default 1)
            -config  string Config file to use
            template string JSON file with details about the request to test 
```

### config
The config command creates a JSOn file with parameters used to setup HTTP connections, with
default values.

* disable-compression 
  > If true, prevents this client from requesting compression  with an "Accept-Encoding: gzip"

* disable-keep-alives 
  > If true, disables HTTP keep-alives and will only use the  connection to the server for a single request 
  
* max-connections 
  > Limits the total number of connections per host,  zero means no limit

* request-timeout
  > Specifies a time limit in seconds for requests made by this Client,  zero means no timeout

* disable-certificate-check
  > If true, disables TLS certificate checking, allowing the use of expired or invalid certificates

* disable-redirects
  > If true, the http client will not follow an HTTP redirect

```sh
$ beast config config.json
File config.json was created with default configuration

$ cat config.json
{
	"disable-compression": true,
	"disable-keep-alives": false,
	"max-connections": 0,
	"request-timeout": 30,
	"disable-certificate-check": false,
	"disable-redirects": true
}
```

### template
The template command functions as a utility to generate request files.

Can be used to generate an "empty" request:
```sh
$ beast template test.json                                                        
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
}%
```

Or, to generate the base request, that you can customize, adding headers or payload:
```sh
$ beast template https://github.com/jjmrocha/beast test.json
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
$ beast run -n 100 -c 5 -config config.json test.json
=== Request ===
Request template: test.json
Configuration: config.json
Number of requests: 100
Number of concurrent requests: 5
=== Test ===
2019/10/05 01:08:30 [#...................] 5%
2019/10/05 01:08:32 [##..................] 10%
2019/10/05 01:08:33 [###.................] 15%
2019/10/05 01:08:35 [####................] 20%
2019/10/05 01:08:37 [#####...............] 25%
2019/10/05 01:08:38 [######..............] 30%
2019/10/05 01:08:39 [#######.............] 35%
2019/10/05 01:08:40 [########............] 40%
2019/10/05 01:08:40 [#########...........] 45%
2019/10/05 01:08:41 [##########..........] 50%
2019/10/05 01:08:41 [###########.........] 55%
2019/10/05 01:08:41 [############........] 60%
2019/10/05 01:08:42 [#############.......] 65%
2019/10/05 01:08:42 [##############......] 70%
2019/10/05 01:08:42 [###############.....] 75%
2019/10/05 01:08:42 [################....] 80%
2019/10/05 01:08:43 [#################...] 85%
2019/10/05 01:08:43 [##################..] 90%
2019/10/05 01:08:43 [###################.] 95%
2019/10/05 01:08:44 [####################] 100%
=== Result Stats ===
Executed requests: 100
Time taken to complete: 1m17.630595087s
Requests per second: 6.4408
Avg response time: 776.30595ms
=== Status 200 ===
40 requests, with avg response time of 1.48163871s
And the following distribution:
  The fastest request took 849.443863ms
  20% of requests under 1.05389514s
  40% of requests under 1.192986887s
  60% of requests under 1.590280578s
  80% of requests under 1.95515696s
  The slowest request took 2.43441918s
=== Non Success Status ===
Status 429: 60 requests
```

## License
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.md).
