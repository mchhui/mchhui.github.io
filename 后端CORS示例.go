package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

// 消息结构体
type Message struct {
    User string `json:"User"`
    Msg  string `json:"Msg"`
    Cate string `json:"Cate"`
    Time string `json:"Time"`
    IP   string `json:"IP"`
}

// 响应结构体
type Response struct {
    Msgs []Message `json:"Msgs"`
}

// 启用CORS的中间件
func enableCORS(w http.ResponseWriter, r *http.Request) {
    // 允许所有域名的跨域请求（开发环境）
    w.Header().Set("Access-Control-Allow-Origin", "*")
    
    // 允许的HTTP方法
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    
    // 允许的请求头
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    
    // 预检请求的缓存时间
    w.Header().Set("Access-Control-Max-Age", "86400")
}

// 处理OPTIONS预检请求
func handleOptions(w http.ResponseWriter, r *http.Request) {
    enableCORS(w, r)
    w.WriteHeader(http.StatusOK)
}

// 发送消息处理器
func msgHandler(w http.ResponseWriter, r *http.Request) {
    enableCORS(w, r)
    
    if r.Method == "OPTIONS" {
        return
    }
    
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // 解析表单数据
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form", http.StatusBadRequest)
        return
    }
    
    user := r.FormValue("user")
    msg := r.FormValue("msg")
    cate := r.FormValue("cate")
    
    if user == "" || msg == "" || cate == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }
    
    // 获取客户端IP
    clientIP := r.RemoteAddr
    if r.Header.Get("X-Forwarded-For") != "" {
        clientIP = r.Header.Get("X-Forwarded-For")
    }
    
    // 创建消息
    message := Message{
        User: user,
        Msg:  msg,
        Cate: cate,
        Time: time.Now().Format("2006-01-02 15:04:05"),
        IP:   clientIP,
    }
    
    // 这里应该将消息保存到数据库
    // 为了示例，我们只是打印到控制台
    log.Printf("新消息: %+v", message)
    
    // 返回成功响应
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "消息发送成功！")
}

// 获取消息列表处理器
func listHandler(w http.ResponseWriter, r *http.Request) {
    enableCORS(w, r)
    
    if r.Method == "OPTIONS" {
        return
    }
    
    if r.Method != "GET" {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    cate := r.URL.Query().Get("cate")
    if cate == "" {
        http.Error(w, "Missing cate parameter", http.StatusBadRequest)
        return
    }
    
    // 这里应该从数据库查询指定分类的消息
    // 为了示例，我们返回空列表
    response := Response{
        Msgs: []Message{},
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func main() {
    // 设置路由
    http.HandleFunc("/msg", msgHandler)
    http.HandleFunc("/list", listHandler)
    
    // 处理所有OPTIONS请求
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            handleOptions(w, r)
            return
        }
        http.NotFound(w, r)
    })
    
    fmt.Println("服务器启动在端口 1066...")
    fmt.Println("支持CORS跨域请求")
    fmt.Println("访问 http://localhost:1066 查看API状态")
    
    log.Fatal(http.ListenAndServe(":1066", nil))
}
