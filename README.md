## 压测

### 测试示例

1、 PHP 原生

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

5、Nodejs
```javascript
var http = require('http');

http.createServer(function (request, response) {
    response.writeHead(200);
    setTimeout(()=>{
        response.end(new Array(4096).fill('s').join(''));
}).listen(8888);
    }, 50);
```
