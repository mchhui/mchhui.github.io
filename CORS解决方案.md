# CORS跨域问题解决方案

## 问题描述

当Hugo博客（通常运行在1313端口）尝试访问后端API（运行在1066端口）时，浏览器会阻止跨域请求，显示CORS错误。

## 解决方案

### 方案1：后端添加CORS头（推荐）

在后端API服务器中添加CORS头，允许跨域请求：

```go
// Go语言示例
func enableCORS(w http.ResponseWriter) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handler(w http.ResponseWriter, r *http.Request) {
    enableCORS(w)
    
    if r.Method == "OPTIONS" {
        return
    }
    
    // 处理实际请求
}
```

```python
# Python Flask示例
from flask import Flask
from flask_cors import CORS

app = Flask(__name__)
CORS(app)  # 允许所有域名的跨域请求

# 或者更精确的配置
CORS(app, origins=["http://localhost:1313", "https://your-domain.com"])
```

```javascript
// Node.js Express示例
const express = require('express');
const cors = require('cors');
const app = express();

app.use(cors()); // 允许所有域名的跨域请求

// 或者更精确的配置
app.use(cors({
    origin: ['http://localhost:1313', 'https://your-domain.com'],
    methods: ['GET', 'POST', 'OPTIONS'],
    allowedHeaders: ['Content-Type']
}));
```

### 方案2：使用Hugo的代理功能

在 `hugo.toml` 中添加代理配置：

```toml
[[module.mounts]]
source = "static"
target = "static"

[[module.mounts]]
source = "layouts"
target = "layouts"

# 添加代理配置
[[module.proxy]]
from = "/api"
to = "http://localhost:1066"
```

然后修改评论组件使用相对路径：

```javascript
const apiBaseUrl = ""; // 使用空字符串，使用相对路径
```

### 方案3：使用Nginx反向代理

配置Nginx将API请求代理到后端：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    # 博客静态文件
    location / {
        root /path/to/hugo/public;
        try_files $uri $uri/ =404;
    }
    
    # API代理
    location /api/ {
        proxy_pass http://localhost:1066/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 方案4：开发环境使用浏览器禁用CORS（不推荐用于生产）

启动Chrome时添加参数：
```bash
chrome.exe --user-data-dir=/tmp/chrome_dev_test --disable-web-security
```

## 当前实现

评论组件已经添加了CORS处理：

1. **请求头设置**：添加了 `mode: 'cors'` 和 `credentials: 'omit'`
2. **错误处理**：检测CORS错误并显示友好提示
3. **调试信息**：在控制台输出详细的请求信息

## 推荐步骤

1. **首选方案1**：在后端添加CORS头，这是最标准的解决方案
2. **备选方案2**：如果无法修改后端，使用Hugo代理功能
3. **生产环境**：使用Nginx反向代理，统一域名

## 测试方法

1. 打开浏览器开发者工具
2. 查看Network标签页
3. 尝试发送评论
4. 检查是否有CORS错误
5. 查看控制台输出的调试信息

## 注意事项

- 生产环境不要使用 `Access-Control-Allow-Origin: *`
- 应该明确指定允许的域名
- 考虑添加认证和CSRF保护
