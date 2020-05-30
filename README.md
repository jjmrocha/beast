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

Documentation
-------------
Please check our [wiki](https://github.com/jjmrocha/beast/wiki).

Usage Manual
------------
```bash
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
            templateFile string JSON/YAML file with details about the request to test

   run      Executes a script and presents a report with execution results
            -c           int    Number of concurrent requests (default 1)
            -n           int    Number of requests (default 1)
            -config      string Config file to setup HTTP client
            -data        string CSV file with data for request generation
            templateFile string JSON/YAML file with details about the request to test 
```

Execution Output
----------------
```bash
$ beast run -n 100 -c 5 -config config.json test.yaml
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
