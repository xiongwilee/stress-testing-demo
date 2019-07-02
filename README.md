[toc]

## 压测

### 测试示例

1、 PHP 原生(起了两个fpm进程)

```php
<?php
usleep(50000); // 模拟业务逻辑操作数据耗时等 50ms
echo str_repeat('s', 4096); // 模拟接口返回4kb数据
```

2、Laravel
```php
<?php
Route::get('/', function () {
    usleep(50000);
    return str_repeat('s', 4096);
});
```

3、Golang原生
```golang
package main

import (
    "net/http"
    "log"  
    "time"
    "strings"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
    time.Sleep(time.Millisecond * 50)
    strings.Repeat("s", 4096)
}

func main(){
    http.HandleFunc("/", myHandler)		
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

4、Gin
```golang
package main
import (
    "github.com/gin-gonic/gin"
    "time"
    "strings"
)
func main()  {
    gin.SetMode(gin.ReleaseMode)
    router := gin.New()
    router.GET("/", func(c *gin.Context) {
        time.Sleep(time.Millisecond * 50)
        c.String(200, strings.Repeat("s", 4096))
    })
    router.Run(":8888")
}
```

5、Nodejs(4个进程)
```javascript
var http = require('http');

http.createServer(function (request, response) {
    response.writeHead(200);
    setTimeout(()=>{
        response.end(new Array(4096).fill('s').join(''));
}).listen(8888);
    }, 50);
```

### 压测结论

1. 在高并发场景下，相比于PHP，Golang/Nodejs表现还是不错的
2. gin框架的引入不会给性能带来明显影响
3. Nginx做Golang的反向代理性能还是很不错的，甚至没有延迟
4. 同机房一定要用内网IP！要用内网IP！内网IP！

### 同机房非统计器压测详情

> 以下结果为同机房非同机器压测，延迟非常明显，通过同机器压测排除业务代码的问题；同机房不同机器的网络带宽是不是有什么问题（找到原因了，外网IP导致的延迟）？
> 中间的数字为平均每个请求的耗时，单位为ms，理想情况为50ms

**外网IP**

||-c1 -n1000|-c5 -n1000|-c10 -n1000|-c20 -n1000|-c30 -n1000|-c50 -n1000|-c100 -n1000|
|-|-|-|-|-|-|-|-|
|php|51.641|152.957|309.705|663.734|1101.118|2739.170|5698.177|
|laravel|67.825|180.184|361.584|749.355|1248.533|2750.871|5874.137|
|golang+nginx|51.710|152.619|309.018|624.468|941.871|1783.164|5489.405|
|golang|51.455|152.037|307.032|616.733|934.472|1918.292|5377.041|
|gin|`51.440`|151.515|308.533|617.120|923.178|1693.219|5546.442|
|nodejs|52.436|`151.412`|`303.335`|`609.632`|`912.297`|`1634.929`|`5376.184`|

**内网IP**

||-c1 -n1000|-c5 -n1000|-c10 -n1000|-c20 -n1000|-c30 -n1000|-c50 -n1000|-c100 -n1000|-c200 -n1000|-c200 -n100000|-c300 -n100000|
|-|-|-|-|-|-|-|-|-|-|-|
|php|51.476|52.564|106.217|208.203|312.735|531.305|1041.147||||
|laravel|69.921|72.547|144.447|285.066|426.115|690.230|1408.361||||
|golang+nginx|51.700|51.754|52.191|52.016|53.004|54.470|`58.321`|100.033|54.924|59.608|
|golang|`51.271`|51.702|52.073|52.696|54.119|54.773|60.854|`93.973`|53.878|`55.954`|
|gin|51.347|`51.444`|52.071|52.771|52.838|`52.060`|69.925|107.791|54.041|56.296|
|nodejs|51.727|51.825|`51.740`|`51.879`|`52.702`|52.169|59.288|94.520|`52.490`|57.088|

#### -c1 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      1
Time taken for tests:   51.641 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    19.36 [#/sec] (mean)
Time per request:       51.641 [ms] (mean)
Time per request:       51.641 [ms] (mean, across all concurrent requests)
Transfer rate:          80.05 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.3      0       6
Processing:    51   51   0.4     51      56
Waiting:       51   51   0.4     51      56
Total:         51   52   0.5     52      57

Percentage of the requests served within a certain time (ms)
  50%     52
  66%     52
  75%     52
  80%     52
  90%     52
  95%     52
  98%     53
  99%     54
 100%     57 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      1
Time taken for tests:   67.825 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000648 bytes
HTML transferred:       4096000 bytes
Requests per second:    14.74 [#/sec] (mean)
Time per request:       67.825 [ms] (mean)
Time per request:       67.825 [ms] (mean, across all concurrent requests)
Transfer rate:          72.00 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.3      0       7
Processing:    65   67   4.7     66     117
Waiting:       65   67   4.7     66     117
Total:         66   68   4.8     67     117

Percentage of the requests served within a certain time (ms)
  50%     67
  66%     67
  75%     68
  80%     68
  90%     70
  95%     72
  98%     76
  99%     99
 100%    117 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      1
Time taken for tests:   51.710 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    19.34 [#/sec] (mean)
Time per request:       51.710 [ms] (mean)
Time per request:       51.710 [ms] (mean, across all concurrent requests)
Transfer rate:          79.96 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.2      0       7
Processing:    51   51   0.5     51      57
Waiting:       51   51   0.5     51      57
Total:         51   52   0.5     52      59

Percentage of the requests served within a certain time (ms)
  50%     52
  66%     52
  75%     52
  80%     52
  90%     52
  95%     52
  98%     53
  99%     54
 100%     59 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      1
Time taken for tests:   51.455 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    19.43 [#/sec] (mean)
Time per request:       51.455 [ms] (mean)
Time per request:       51.455 [ms] (mean, across all concurrent requests)
Transfer rate:          79.58 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.6      0      19
Processing:    51   51   0.5     51      59
Waiting:       51   51   0.5     51      59
Total:         51   51   0.8     51      70

Percentage of the requests served within a certain time (ms)
  50%     51
  66%     51
  75%     51
  80%     51
  90%     52
  95%     52
  98%     53
  99%     54
 100%     70 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      1
Time taken for tests:   51.440 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    19.44 [#/sec] (mean)
Time per request:       51.440 [ms] (mean)
Time per request:       51.440 [ms] (mean, across all concurrent requests)
Transfer rate:          79.60 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.2      0       4
Processing:    51   51   0.4     51      56
Waiting:       51   51   0.4     51      56
Total:         51   51   0.5     51      56

Percentage of the requests served within a certain time (ms)
  50%     51
  66%     51
  75%     51
  80%     51
  90%     52
  95%     52
  98%     53
  99%     54
 100%     56 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        php.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      1
Time taken for tests:   52.436 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    19.07 [#/sec] (mean)
Time per request:       52.436 [ms] (mean)
Time per request:       52.436 [ms] (mean, across all concurrent requests)
Transfer rate:          77.68 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.2      0       5
Processing:    50   52   1.6     52      69
Waiting:       50   52   1.3     52      67
Total:         51   52   1.6     52      70

Percentage of the requests served within a certain time (ms)
  50%     52
  66%     52
  75%     52
  80%     52
  90%     53
  95%     53
  98%     55
  99%     57
 100%     70 (longest request)
```


#### -c5 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      5
Time taken for tests:   30.591 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.69 [#/sec] (mean)
Time per request:       152.957 [ms] (mean)
Time per request:       30.591 [ms] (mean, across all concurrent requests)
Transfer rate:          135.13 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    3  54.7      0    1000
Processing:    51  148 156.2     51     895
Waiting:       51  131 125.5     51     695
Total:         51  151 163.7     52    1051

Percentage of the requests served within a certain time (ms)
  50%     52
  66%    252
  75%    252
  80%    253
  90%    256
  95%    465
  98%    683
  99%    692
 100%   1051 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      5
Time taken for tests:   36.037 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000906 bytes
HTML transferred:       4096000 bytes
Requests per second:    27.75 [#/sec] (mean)
Time per request:       180.184 [ms] (mean)
Time per request:       36.037 [ms] (mean, across all concurrent requests)
Transfer rate:          135.52 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    9  94.3      0    1000
Processing:    64  170 169.6     68    1752
Waiting:       64  145 131.5     67    1544
Total:         64  180 192.5     68    1753

Percentage of the requests served within a certain time (ms)
  50%     68
  66%    268
  75%    269
  80%    271
  90%    475
  95%    480
  98%    708
  99%   1064
 100%   1753 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      5
Time taken for tests:   30.524 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.76 [#/sec] (mean)
Time per request:       152.619 [ms] (mean)
Time per request:       30.524 [ms] (mean, across all concurrent requests)
Transfer rate:          135.46 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    7  83.3      0    1000
Processing:    51  145 148.8     51     892
Waiting:       51  126 114.1     51     691
Total:         51  152 169.1     52    1260

Percentage of the requests served within a certain time (ms)
  50%     52
  66%    252
  75%    252
  80%    253
  90%    258
  95%    464
  98%    680
  99%    889
 100%   1260 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      5
Time taken for tests:   30.407 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.89 [#/sec] (mean)
Time per request:       152.037 [ms] (mean)
Time per request:       30.407 [ms] (mean, across all concurrent requests)
Transfer rate:          134.66 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    4  63.1      0    1000
Processing:    51  147 155.2     51    1722
Waiting:       51  128 117.8     51     695
Total:         51  152 165.2     51    1722

Percentage of the requests served within a certain time (ms)
  50%     51
  66%    252
  75%    252
  80%    253
  90%    257
  95%    461
  98%    672
  99%    696
 100%   1722 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      5
Time taken for tests:   30.303 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    33.00 [#/sec] (mean)
Time per request:       151.515 [ms] (mean)
Time per request:       30.303 [ms] (mean, across all concurrent requests)
Transfer rate:          135.13 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    5  70.5      0    1000
Processing:    51  146 153.7     51    1085
Waiting:       51  128 117.7     51     694
Total:         51  151 167.5     51    1253

Percentage of the requests served within a certain time (ms)
  50%     51
  66%    252
  75%    252
  80%    253
  90%    256
  95%    463
  98%    672
  99%    880
 100%   1253 (longest request)
```

##### nodejs
```
Server Software:        nginx/1.10.3
Server Hostname:        nodejs.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      5
Time taken for tests:   30.282 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    33.02 [#/sec] (mean)
Time per request:       151.412 [ms] (mean)
Time per request:       30.282 [ms] (mean, across all concurrent requests)
Transfer rate:          135.22 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    7  83.3      0    1000
Processing:    51  144 151.0     52    1320
Waiting:       51  127 116.5     52     691
Total:         51  151 171.1     53    1461

Percentage of the requests served within a certain time (ms)
  50%     53
  66%    252
  75%    256
  80%    256
  90%    260
  95%    464
  98%    680
  99%    892
 100%   1461 (longest request)
```

#### -c10 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      10
Time taken for tests:   30.970 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.29 [#/sec] (mean)
Time per request:       309.705 [ms] (mean)
Time per request:       30.970 [ms] (mean, across all concurrent requests)
Transfer rate:          133.48 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   15 121.4      0    1000
Processing:    51  290 335.8    153    3188
Waiting:       51  193 214.7    150    3187
Total:         51  305 355.8    252    3188

Percentage of the requests served within a certain time (ms)
  50%    252
  66%    257
  75%    460
  80%    465
  90%    716
  95%   1051
  98%   1464
  99%   1724
 100%   3188 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      10
Time taken for tests:   36.158 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000548 bytes
HTML transferred:       4096000 bytes
Requests per second:    27.66 [#/sec] (mean)
Time per request:       361.584 [ms] (mean)
Time per request:       36.158 [ms] (mean, across all concurrent requests)
Transfer rate:          135.05 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   13 113.2      0    1000
Processing:    64  346 346.4    268    3264
Waiting:       64  226 205.5    266    1580
Total:         65  360 361.3    268    3264

Percentage of the requests served within a certain time (ms)
  50%    268
  66%    475
  75%    480
  80%    684
  90%    892
  95%   1112
  98%   1476
  99%   1557
 100%   3264 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      10
Time taken for tests:   30.902 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.36 [#/sec] (mean)
Time per request:       309.018 [ms] (mean)
Time per request:       30.902 [ms] (mean, across all concurrent requests)
Transfer rate:          133.80 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    6  77.2      0    1000
Processing:    51  299 251.7    259    1952
Waiting:       51  193 187.0    251    1503
Total:         51  305 258.3    260    1952

Percentage of the requests served within a certain time (ms)
  50%    260
  66%    264
  75%    464
  80%    468
  90%    680
  95%    700
  98%    896
  99%   1104
 100%   1952 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      10
Time taken for tests:   30.703 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.57 [#/sec] (mean)
Time per request:       307.032 [ms] (mean)
Time per request:       30.703 [ms] (mean, across all concurrent requests)
Transfer rate:          133.36 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   15 121.5      0    1000
Processing:    51  290 291.7    252    2983
Waiting:       51  214 213.7    251    1513
Total:         51  305 305.1    252    2984

Percentage of the requests served within a certain time (ms)
  50%    252
  66%    460
  75%    464
  80%    465
  90%    684
  95%    888
  98%   1051
  99%   1252
 100%   2984 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      10
Time taken for tests:   30.853 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.41 [#/sec] (mean)
Time per request:       308.533 [ms] (mean)
Time per request:       30.853 [ms] (mean, across all concurrent requests)
Transfer rate:          132.72 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   11 104.2      0    1000
Processing:    51  293 251.6    254    1738
Waiting:       51  215 193.8    251    1527
Total:         51  304 273.2    254    2504

Percentage of the requests served within a certain time (ms)
  50%    254
  66%    460
  75%    464
  80%    464
  90%    672
  95%    692
  98%    900
  99%   1114
 100%   2504 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      10
Time taken for tests:   30.333 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.97 [#/sec] (mean)
Time per request:       303.335 [ms] (mean)
Time per request:       30.333 [ms] (mean, across all concurrent requests)
Transfer rate:          134.28 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   21 143.4      0    1001
Processing:    50  280 274.2    253    1716
Waiting:       50  199 192.1    252    1531
Total:         51  301 304.6    255    2520

Percentage of the requests served within a certain time (ms)
  50%    255
  66%    460
  75%    464
  80%    468
  90%    676
  95%    893
  98%   1080
  99%   1304
 100%   2520 (longest request)
```

#### -c20 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      20
Time taken for tests:   33.187 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    30.13 [#/sec] (mean)
Time per request:       663.734 [ms] (mean)
Time per request:       33.187 [ms] (mean, across all concurrent requests)
Transfer rate:          124.56 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   35 183.7      0    1000
Processing:    51  577 1247.6    259   27228
Waiting:       51  309 608.7    252   13119
Total:         51  613 1259.4    303   27228

Percentage of the requests served within a certain time (ms)
  50%    303
  66%    666
  75%    700
  80%    896
  90%   1464
  95%   1736
  98%   3164
  99%   3808
 100%  27228 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      20
Time taken for tests:   37.468 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000646 bytes
HTML transferred:       4096000 bytes
Requests per second:    26.69 [#/sec] (mean)
Time per request:       749.355 [ms] (mean)
Time per request:       37.468 [ms] (mean, across all concurrent requests)
Transfer rate:          130.34 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   21 143.3      0    1000
Processing:    65  702 751.6    480    7486
Waiting:       65  391 542.8    270    6654
Total:         65  723 758.9    488    7486

Percentage of the requests served within a certain time (ms)
  50%    488
  66%    888
  75%   1068
  80%   1120
  90%   1567
  95%   1936
  98%   2384
  99%   3465
 100%   7486 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      20
Time taken for tests:   31.223 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.03 [#/sec] (mean)
Time per request:       624.468 [ms] (mean)
Time per request:       31.223 [ms] (mean, across all concurrent requests)
Transfer rate:          132.43 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   19 136.5      0    1004
Processing:    51  593 656.0    464    6656
Waiting:       51  400 547.0    255    6656
Total:         51  612 655.1    465    6657

Percentage of the requests served within a certain time (ms)
  50%    465
  66%    676
  75%    880
  80%    896
  90%   1108
  95%   1528
  98%   2157
  99%   3600
 100%   6657 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      20
Time taken for tests:   30.837 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.43 [#/sec] (mean)
Time per request:       616.733 [ms] (mean)
Time per request:       30.837 [ms] (mean, across all concurrent requests)
Transfer rate:          132.79 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   39 193.6      0    1002
Processing:    51  569 554.7    463    6626
Waiting:       51  372 436.2    254    6626
Total:         51  608 565.6    464    6626

Percentage of the requests served within a certain time (ms)
  50%    464
  66%    676
  75%    884
  80%    908
  90%   1508
  95%   1712
  98%   1732
  99%   2152
 100%   6626 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      20
Time taken for tests:   30.856 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.41 [#/sec] (mean)
Time per request:       617.120 [ms] (mean)
Time per request:       30.856 [ms] (mean, across all concurrent requests)
Transfer rate:          132.70 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   24 153.0      0    1000
Processing:    51  585 505.5    464    3375
Waiting:       51  384 386.8    255    3167
Total:         51  609 506.5    468    3375

Percentage of the requests served within a certain time (ms)
  50%    468
  66%    682
  75%    884
  80%    898
  90%   1344
  95%   1530
  98%   1917
  99%   2148
 100%   3375 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        php.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      20
Time taken for tests:   30.482 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.81 [#/sec] (mean)
Time per request:       609.632 [ms] (mean)
Time per request:       30.482 [ms] (mean, across all concurrent requests)
Transfer rate:          133.63 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   30 170.6      0    1004
Processing:    50  571 562.4    464    3805
Waiting:       50  381 455.8    255    3192
Total:         51  601 563.8    468    3806

Percentage of the requests served within a certain time (ms)
  50%    468
  66%    685
  75%    880
  80%    900
  90%   1515
  95%   1712
  98%   1922
  99%   2156
 100%   3806 (longest request)
```

#### -c30 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      30
Time taken for tests:   36.704 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    27.25 [#/sec] (mean)
Time per request:       1101.118 [ms] (mean)
Time per request:       36.704 [ms] (mean, across all concurrent requests)
Transfer rate:          112.63 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   37 188.7      0    1000
Processing:    51  889 1570.3    462   28383
Waiting:       51  431 800.5    253   13184
Total:         51  926 1583.9    464   28383

Percentage of the requests served within a certain time (ms)
  50%    464
  66%    876
  75%   1104
  80%   1464
  90%   2133
  95%   3392
  98%   4016
  99%   7092
 100%  28383 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      30
Time taken for tests:   41.618 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000504 bytes
HTML transferred:       4096000 bytes
Requests per second:    24.03 [#/sec] (mean)
Time per request:       1248.533 [ms] (mean)
Time per request:       41.618 [ms] (mean, across all concurrent requests)
Transfer rate:          117.34 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   36 186.3      0    1004
Processing:    65 1052 1154.7    688   14145
Waiting:       65  518 704.9    272    6720
Total:         65 1089 1162.1    692   14148

Percentage of the requests served within a certain time (ms)
  50%    692
  66%   1102
  75%   1328
  80%   1552
  90%   2367
  95%   3408
  98%   4048
  99%   6704
 100%  14148 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      30
Time taken for tests:   31.396 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    31.85 [#/sec] (mean)
Time per request:       941.871 [ms] (mean)
Time per request:       31.396 [ms] (mean, across all concurrent requests)
Transfer rate:          131.70 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   55 228.0      0    1004
Processing:    51  860 1013.2    670   13192
Waiting:       51  556 875.4    255   13192
Total:         51  915 1003.3    673   13193

Percentage of the requests served within a certain time (ms)
  50%    673
  66%   1052
  75%   1252
  80%   1508
  90%   1744
  95%   3172
  98%   3594
  99%   3836
 100%  13193 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        php.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      30
Time taken for tests:   31.149 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.10 [#/sec] (mean)
Time per request:       934.472 [ms] (mean)
Time per request:       31.149 [ms] (mean, across all concurrent requests)
Transfer rate:          131.46 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   42 200.6      0    1003
Processing:    51  870 1003.5    667    8118
Waiting:       51  520 727.4    252    6671
Total:         51  912 999.8    669    8119

Percentage of the requests served within a certain time (ms)
  50%    669
  66%   1051
  75%   1309
  80%   1516
  90%   1936
  95%   3181
  98%   3398
  99%   4032
 100%   8119 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        php.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      30
Time taken for tests:   30.773 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.50 [#/sec] (mean)
Time per request:       923.178 [ms] (mean)
Time per request:       30.773 [ms] (mean, across all concurrent requests)
Transfer rate:          133.06 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   24 153.0      0    1001
Processing:    51  883 2192.0    467   26558
Waiting:       51  616 2125.9    255   26141
Total:         51  908 2188.7    668   26559

Percentage of the requests served within a certain time (ms)
  50%    668
  66%    880
  75%   1088
  80%   1112
  90%   1712
  95%   1940
  98%   3377
  99%   4639
 100%  26559 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        php.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      30
Time taken for tests:   30.410 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    32.88 [#/sec] (mean)
Time per request:       912.297 [ms] (mean)
Time per request:       30.410 [ms] (mean, across all concurrent requests)
Transfer rate:          133.94 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   22 146.7      0    1000
Processing:    50  877 908.1    671    6768
Waiting:       50  525 696.1    255    6560
Total:         51  899 909.1    672    6768

Percentage of the requests served within a certain time (ms)
  50%    672
  66%   1087
  75%   1504
  80%   1520
  90%   1944
  95%   3164
  98%   3392
  99%   3600
 100%   6768 (longest request)
```

#### -c50 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      50
Time taken for tests:   54.783 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.25 [#/sec] (mean)
Time per request:       2739.170 [ms] (mean)
Time per request:       54.783 [ms] (mean, across all concurrent requests)
Transfer rate:          75.46 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   51 220.0      0    1004
Processing:    51 1551 3482.5    792   54024
Waiting:       51  669 2092.2    256   53837
Total:         51 1603 3484.4    878   54025

Percentage of the requests served within a certain time (ms)
  50%    878
  66%   1107
  75%   1536
  80%   1920
  90%   3379
  95%   6696
  98%   7920
  99%  13568
 100%  54025 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      50
Time taken for tests:   55.017 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000686 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.18 [#/sec] (mean)
Time per request:       2750.871 [ms] (mean)
Time per request:       55.017 [ms] (mean, across all concurrent requests)
Transfer rate:          88.76 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   50 222.6      0    2008
Processing:    64 1834 4171.1    836   54463
Waiting:       64  847 2623.5    272   54248
Total:         65 1884 4177.0    894   54463

Percentage of the requests served within a certain time (ms)
  50%    894
  66%   1227
  75%   1732
  80%   2068
  90%   3724
  95%   6805
  98%  13468
  99%  20819
 100%  54463 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      50
Time taken for tests:   35.663 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    28.04 [#/sec] (mean)
Time per request:       1783.164 [ms] (mean)
Time per request:       35.663 [ms] (mean, across all concurrent requests)
Transfer rate:          115.94 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   47 216.5      0    2021
Processing:    51 1473 2369.5    676   27507
Waiting:       51  924 2038.0    255   26640
Total:         51 1521 2363.4    880   27508

Percentage of the requests served within a certain time (ms)
  50%    880
  66%   1512
  75%   1729
  80%   2164
  90%   3380
  95%   4016
  98%   7559
  99%  13920
 100%  27508 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      50
Time taken for tests:   38.366 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    26.06 [#/sec] (mean)
Time per request:       1918.292 [ms] (mean)
Time per request:       38.366 [ms] (mean, across all concurrent requests)
Transfer rate:          106.73 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   63 247.2      0    2024
Processing:    51 1472 2815.0    668   27276
Waiting:       51  716 1564.9    252   13288
Total:         51 1536 2819.1    671   27276

Percentage of the requests served within a certain time (ms)
  50%    671
  66%   1101
  75%   1543
  80%   1940
  90%   3396
  95%   6912
  98%   9471
  99%  13888
 100%  27276 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      50
Time taken for tests:   33.864 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    29.53 [#/sec] (mean)
Time per request:       1693.219 [ms] (mean)
Time per request:       33.864 [ms] (mean, across all concurrent requests)
Transfer rate:          120.92 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   74 265.7      0    2032
Processing:    51 1454 2922.9    667   27280
Waiting:       51  848 2615.6    252   26642
Total:         51 1529 2928.7    672   27282

Percentage of the requests served within a certain time (ms)
  50%    672
  66%   1088
  75%   1520
  80%   1920
  90%   3392
  95%   6704
  98%   8552
  99%  13504
 100%  27282 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      50
Time taken for tests:   32.699 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    30.58 [#/sec] (mean)
Time per request:       1634.929 [ms] (mean)
Time per request:       32.699 [ms] (mean, across all concurrent requests)
Transfer rate:          124.57 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   56 229.9      0    1000
Processing:    51 1451 2152.1    675   26607
Waiting:       50  876 1845.2    256   26606
Total:         51 1507 2153.0    880   26607

Percentage of the requests served within a certain time (ms)
  50%    880
  66%   1260
  75%   1724
  80%   2352
  90%   3600
  95%   6544
  98%   7307
  99%   9968
 100%  26607 (longest request)
```

#### -c100 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      100
Time taken for tests:   56.982 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4233000 bytes
HTML transferred:       4096000 bytes
Requests per second:    17.55 [#/sec] (mean)
Time per request:       5698.177 [ms] (mean)
Time per request:       56.982 [ms] (mean, across all concurrent requests)
Transfer rate:          72.55 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   64 242.9      0    1002
Processing:    51 3503 8510.6   1087   56270
Waiting:       51 1641 6161.6    256   54574
Total:         51 3567 8508.5   1088   56271

Percentage of the requests served within a certain time (ms)
  50%   1088
  66%   1712
  75%   2356
  80%   3376
  90%   7408
  95%  14189
  98%  53718
  99%  54374
 100%  56271 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      100
Time taken for tests:   58.741 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      5000756 bytes
HTML transferred:       4096000 bytes
Requests per second:    17.02 [#/sec] (mean)
Time per request:       5874.137 [ms] (mean)
Time per request:       58.741 [ms] (mean, across all concurrent requests)
Transfer rate:          83.14 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   47 209.5      0    1003
Processing:    65 4065 9491.0   1156   56979
Waiting:       65 1745 6292.3    280   56689
Total:         65 4111 9489.8   1322   56980

Percentage of the requests served within a certain time (ms)
  50%   1322
  66%   1964
  75%   2560
  80%   3604
  90%   7588
  95%  15380
  98%  55468
  99%  56691
 100%  56980 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      100
Time taken for tests:   54.894 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.22 [#/sec] (mean)
Time per request:       5489.405 [ms] (mean)
Time per request:       54.894 [ms] (mean, across all concurrent requests)
Transfer rate:          75.32 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   59 233.6      0    1008
Processing:    51 4228 12369.4    668   54829
Waiting:       51 3730 12401.7    255   54829
Total:         51 4286 12355.8    672   54830

Percentage of the requests served within a certain time (ms)
  50%    672
  66%   1088
  75%   1520
  80%   1932
  90%   4836
  95%  52840
  98%  53256
  99%  54822
 100%  54830 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      100
Time taken for tests:   53.770 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.60 [#/sec] (mean)
Time per request:       5377.041 [ms] (mean)
Time per request:       53.770 [ms] (mean, across all concurrent requests)
Transfer rate:          76.15 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   65 244.5      0    1000
Processing:    51 3847 10653.0    675   53766
Waiting:       51 3173 10652.0    254   53351
Total:         51 3912 10639.2    900   53768

Percentage of the requests served within a certain time (ms)
  50%    900
  66%   1668
  75%   3164
  80%   3376
  90%   6712
  95%  14008
  98%  53355
  99%  53561
 100%  53768 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      100
Time taken for tests:   55.464 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.03 [#/sec] (mean)
Time per request:       5546.442 [ms] (mean)
Time per request:       55.464 [ms] (mean, across all concurrent requests)
Transfer rate:          73.83 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   58 233.7      0    1003
Processing:    51 4098 11797.0    671   55462
Waiting:       51 3579 11782.5    253   54823
Total:         51 4156 11782.4    876   55464

Percentage of the requests served within a certain time (ms)
  50%    876
  66%   1508
  75%   1920
  80%   3164
  90%   4656
  95%  27079
  98%  55185
  99%  55239
 100%  55464 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      100
Time taken for tests:   53.762 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.60 [#/sec] (mean)
Time per request:       5376.184 [ms] (mean)
Time per request:       53.762 [ms] (mean, across all concurrent requests)
Transfer rate:          75.76 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   35 181.2      0    1001
Processing:    51 4155 11903.7    676   53758
Waiting:       50 3568 11865.2    256   53119
Total:         51 4189 11895.7    880   53761

Percentage of the requests served within a certain time (ms)
  50%    880
  66%   1504
  75%   1932
  80%   2368
  90%   4016
  95%  53065
  98%  53479
  99%  53694
 100%  53761 (longest request)
```


#### -c200 -n1000

##### php
```
Server Software:        nginx/1.10.3
Server Hostname:        php.stresstesting
Server Port:            80

Document Path:          /
Document Length:        173 bytes

Concurrency Level:      200
Time taken for tests:   32.389 seconds
Complete requests:      1000
Failed requests:        719
   (Connect: 0, Receive: 0, Length: 719, Exceptions: 0)
Write errors:           0
Non-2xx responses:      281
Total transferred:      3134852 bytes
HTML transferred:       2993637 bytes
Requests per second:    30.87 [#/sec] (mean)
Time per request:       6477.872 [ms] (mean)
Time per request:       32.389 [ms] (mean, across all concurrent requests)
Transfer rate:          94.52 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   81 271.3      0    1010
Processing:    18 4787 8383.4   1088   30060
Waiting:       18 2951 7076.7    256   28544
Total:         20 4868 8395.1   1092   30061

Percentage of the requests served within a certain time (ms)
  50%   1092
  66%   2307
  75%   4016
  80%   6848
  90%  26168
  95%  27940
  98%  28292
  99%  28514
 100%  30061 (longest request)
```

##### laravel
```
Server Software:        nginx/1.10.3
Server Hostname:        laravel5.stresstesting
Server Port:            80

Document Path:          /
Document Length:        173 bytes

Concurrency Level:      200
Time taken for tests:   35.094 seconds
Complete requests:      1000
Failed requests:        720
   (Connect: 0, Receive: 0, Length: 720, Exceptions: 0)
Write errors:           0
Non-2xx responses:      280
Total transferred:      3691402 bytes
HTML transferred:       2997560 bytes
Requests per second:    28.50 [#/sec] (mean)
Time per request:       7018.758 [ms] (mean)
Time per request:       35.094 [ms] (mean, across all concurrent requests)
Transfer rate:          102.72 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   84 275.6      0    1004
Processing:     1 5339 9047.9   1120   31250
Waiting:        1 2872 6972.6    283   28689
Total:          1 5423 9072.6   1277   31254

Percentage of the requests served within a certain time (ms)
  50%   1277
  66%   2576
  75%   4224
  80%   6952
  90%  26833
  95%  28182
  98%  28768
  99%  29177
 100%  31254 (longest request)
```

##### golang+nginx
```
Server Software:        nginx/1.10.3
Server Hostname:        golang.stresstesting
Server Port:            80

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      200
Time taken for tests:   54.468 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4234000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.36 [#/sec] (mean)
Time per request:       10893.667 [ms] (mean)
Time per request:       54.468 [ms] (mean, across all concurrent requests)
Transfer rate:          75.91 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   68 249.9      0    1001
Processing:    51 7368 14655.5    896   54463
Waiting:       51 6317 14437.9    256   52992
Total:         51 7436 14633.0   1088   54467

Percentage of the requests served within a certain time (ms)
  50%   1088
  66%   2712
  75%   4160
  80%   7124
  90%  27474
  95%  53410
  98%  53841
  99%  54466
 100%  54467 (longest request)
```

##### golang
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3001

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      200
Time taken for tests:   55.328 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.07 [#/sec] (mean)
Time per request:       11065.614 [ms] (mean)
Time per request:       55.328 [ms] (mean, across all concurrent requests)
Transfer rate:          74.01 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   60 235.2      0    1000
Processing:    51 7920 16178.0    880   55322
Waiting:       51 6875 16088.6    255   54682
Total:         51 7980 16156.5   1084   55325

Percentage of the requests served within a certain time (ms)
  50%   1084
  66%   1968
  75%   3393
  80%   6735
  90%  27740
  95%  53278
  98%  55100
  99%  55324
 100%  55325 (longest request)
```

##### gin
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3002

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      200
Time taken for tests:   55.310 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4193000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.08 [#/sec] (mean)
Time per request:       11062.061 [ms] (mean)
Time per request:       55.310 [ms] (mean, across all concurrent requests)
Transfer rate:          74.03 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   69 251.6      0    1006
Processing:    51 7485 15188.0    896   55305
Waiting:       51 6725 15040.7    255   53833
Total:         51 7553 15164.8   1088   55307

Percentage of the requests served within a certain time (ms)
  50%   1088
  66%   2344
  75%   3816
  80%   6708
  90%  27659
  95%  54250
  98%  55306
  99%  55307
 100%  55307 (longest request)
```

##### nodejs
```
Server Software:
Server Hostname:        golang.stresstesting
Server Port:            3000

Document Path:          /
Document Length:        4096 bytes

Concurrency Level:      200
Time taken for tests:   54.884 seconds
Complete requests:      1000
Failed requests:        0
Write errors:           0
Total transferred:      4171000 bytes
HTML transferred:       4096000 bytes
Requests per second:    18.22 [#/sec] (mean)
Time per request:       10976.872 [ms] (mean)
Time per request:       54.884 [ms] (mean, across all concurrent requests)
Transfer rate:          74.21 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0   55 225.9      0    1004
Processing:    50 7137 14463.7   1087   54880
Waiting:       50 6083 14146.8    259   53409
Total:         51 7192 14444.6   1088   54883

Percentage of the requests served within a certain time (ms)
  50%   1088
  66%   1744
  75%   3400
  80%   6768
  90%  28002
  95%  53412
  98%  54258
  99%  54259
 100%  54883 (longest request)
```

