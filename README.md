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
   beast run (-n <number of requests> | -t <test duration>)
             [-c <number of concurrent requests>]
             [-config <configFile>] [-data <dataFile>]
             [-output <outputFile>] <templateFile>

Where:
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
            -output      string CVS file with detailed execution results
            templateFile string JSON/YAML file with details about the request to test
```

Execution Output
----------------
```bash
$ beast run -n 100000 -c 100 -config config.json -data ../test_data.csv -output get_100.csv apps_get.yaml
===== System =====
Operating System: darwin
System Architecture: amd64
Logical CPUs: 12
===== Test =====
Request template: apps_get.yaml
Sample Data: ../test_data.csv
Configuration: config.json
Number of requests: 10000
Number of concurrent requests: 100
===== Preparing =====
- Reading configuration
- Loading request template
- Loading data file
===== Executing =====
2020/07/07 23:03:33 [#...................] 5%
2020/07/07 23:03:37 [##..................] 10%
2020/07/07 23:03:40 [###.................] 15%
2020/07/07 23:03:44 [####................] 20%
2020/07/07 23:03:47 [#####...............] 25%
2020/07/07 23:03:49 [######..............] 30%
2020/07/07 23:03:52 [#######.............] 35%
2020/07/07 23:03:54 [########............] 40%
2020/07/07 23:03:56 [#########...........] 45%
2020/07/07 23:03:59 [##########..........] 50%
2020/07/07 23:04:01 [###########.........] 55%
2020/07/07 23:04:04 [############........] 60%
2020/07/07 23:04:06 [#############.......] 65%
2020/07/07 23:04:09 [##############......] 70%
2020/07/07 23:04:11 [###############.....] 75%
2020/07/07 23:04:14 [################....] 80%
2020/07/07 23:04:16 [#################...] 85%
2020/07/07 23:04:19 [##################..] 90%
2020/07/07 23:04:21 [###################.] 95%
2020/07/07 23:04:23 [####################] 100%
===== Stats =====
Executed requests: 100000
Time taken to complete: 55.167930783s
Requests per second: 1860.9974
Avg response time: 53.734627ms
===== Status 200 =====
100000 requests, with avg response time of 53.734627ms
And the following distribution:
- The fastest request took 6.988515ms
- 20% of requests under 30.239955ms
- 40% of requests under 40.23343ms
- 60% of requests under 52.59765ms
- 80% of requests under 72.891344ms
- 90% of requests under 91.691174ms
- 95% of requests under 110.193479ms
- 99% of requests under 152.862198ms
- The slowest request took 311.154579ms
===== Output File =====
Output file 'get_100.csv' was successfully generated
```

License
-------
Any contributions made under this project will be governed by the [Apache License 2.0](./LICENSE.md).
