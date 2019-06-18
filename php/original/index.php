<?php
usleep(50000); // 模拟业务逻辑操作数据耗时等 50ms
echo str_repeat('s', 4096); // 模拟接口返回4kb数据