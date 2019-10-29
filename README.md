the Beast
=========
*Stress testing for RESTful APIs*

Installation
------------
Binaries are available on [releases](https://github.com/jjmrocha/beast/releases).

Or you can download the code and compile yourself:
```
$ go get -u github.com/jjmrocha/beast
```

Usage
--------
Beast currently supports the following commands:
* [help](#beast-help)
* [config](#beast-config)
* [template](#beast-template)
* [run](#beast-run)

### beast help
Displays the help information.

__Usage__
```
beast [help]
```

Example:
```
$ beast help
the Beast v2.x.x - Stress testing for RESTful APIs

Usage:
   beast [help]
   beast config <configFile>
   beast template [-m <http method>] [url] <templateFile>
   beast run [-n <number of requests>] [-c <number of concurrent requests>] 
             [-config <configFile>] [-data <dataFile>] <templateFile>

Where:
   config   Creates a file with the default parameters to setup HTTP connections
            configFile   string Name of the file to be created
			 			
   template Creates a request template file, using user-provided parameters
            -m           string HTTP method (default "GET")
            url          string Endpoint to be tested
            templateFile string JSON file with details about the request to test

   run      Executes a script and presents a report with execution results
            -c           int    Number of concurrent requests (default 1)
            -n           int    Number of requests (default 1)
            -config      string Config file to setup HTTP client
            -data        string CSV file with data for request generation
            templateFile string JSON file with details about the request to test 
```

### beast config
The config command creates a JSON file with parameters used to setup HTTP connections, with
default values.

__Usage__
```
beast config <configFile>
```
* configFile
  > Name of the file to be created with the default configuration

__Configuration file__
```
$ beast config config.json
File config.json was created with the default configuration

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
  > If true, the HTTP client will not follow an HTTP redirect

### beast template
The template command functions as a utility to generate request files.

__Usage__
```
beast template [-m <http method>] [url] <templateFile>
```
* http method
  > HTTP method to use on the request, defaults to GET

* url
  > Endpoint to test

* templateFile
  > Name of the file to be created with the request template

__Examples__
Can be used to generate an "empty" request template:
```
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
}
```

Or, to generate the base request, that you can customize, adding headers or payload:
```
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

### beast run
The run command loads a request from a file and executes the request concurrently multiple times.

__Usage__
```
beast run [-n <number of requests>] [-c <number of concurrent requests>] 
          [-config <configFile>] [-data <dataFile>] <templateFile>
```
* number of requests
  > Number the request to be performed, defaults to 1

* number of concurrent requests
  > Number of requests to be performed simultaneously, defaults to 1

* configFile
  > Name of the file with the configuration to set up HTTP connections

* dataFile
  > Name of the CSV files to be used on the generation of dynamic requests

* templateFile
  > Name of the file with request template to test

__Data files__
Data files are CSV files with data that can be replaced on the request template, to generate
dynamic requests.

The first line should contain the name of the columns.

__Template language__
The template language used on the template files is the [GO template](https://golang.org/pkg/text/template/).

The fields that may contain dynamic expressions are the following:
* url
* headers/value
* body

Special features implemented:
1. To include the request id you can use ```{{ .RequestID }}```
2. To include a value from the data file, use ```{{ .Data.<column name> }}```

```
{
	"method": "POST",
	"url": "http://someendpoint.pt/{{ .RequestID }}",
	"headers": [
		{
			"key": "Content-Type",
			"value": "application/json"
		}
	],
	"body": "{\"id\": {{ .RequestID }}, \"value\": \"{{ .Data.A }}\"}"
}
```

__Example__
```
$ beast run -n 100 -c 5 -config config.json test.json
===== System =====
Operating System: darwin
System Architecture: amd64
Logical CPUs: 4
===== Test =====
Request template: test.json
Configuration: config.json
Number of requests: 100
Number of concurrent requests: 5
===== Preparing =====
- Reading configuration
- Loading request template
- Generating requests
===== Executing =====
2019/10/29 21:39:22 [#...................] 5%
2019/10/29 21:39:24 [##..................] 10%
2019/10/29 21:39:26 [###.................] 15%
2019/10/29 21:39:28 [####................] 20%
2019/10/29 21:39:30 [#####...............] 25%
2019/10/29 21:39:32 [######..............] 30%
2019/10/29 21:39:35 [#######.............] 35%
2019/10/29 21:39:37 [########............] 40%
2019/10/29 21:39:39 [#########...........] 45%
2019/10/29 21:39:41 [##########..........] 50%
2019/10/29 21:39:43 [###########.........] 55%
2019/10/29 21:39:45 [############........] 60%
2019/10/29 21:39:47 [#############.......] 65%
2019/10/29 21:39:50 [##############......] 70%
2019/10/29 21:39:52 [###############.....] 75%
2019/10/29 21:39:54 [################....] 80%
2019/10/29 21:39:56 [#################...] 85%
2019/10/29 21:39:58 [##################..] 90%
2019/10/29 21:40:01 [###################.] 95%
2019/10/29 21:40:04 [####################] 100%
===== Stats =====
Executed requests: 100
Time taken to complete: 1m17.630595087s
Requests per second: 6.4408
Avg response time: 776.30595ms
===== Status 200 =====
40 requests, with avg response time of 1.48163871s
And the following distribution:
- The fastest request took 849.443863ms
- 20% of requests under 1.05389514s
- 40% of requests under 1.192986887s
- 60% of requests under 1.590280578s
- 80% of requests under 1.95515696s
- The slowest request took 2.43441918s
===== Non Success Status =====
Status 429: 50 requests
===== Errors =====
- Request timeout: 10 errors
```

License
-------
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.md).
