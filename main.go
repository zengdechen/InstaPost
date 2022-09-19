package main

import (
    "fmt"
    "log"
    "net/http"

    "around/backend"
    "around/handler"
)


func main() {
    fmt.Println("started-service")

    backend.InitElasticsearchBackend()   // 创建ES
    backend.InitGCSBackend()             // 创建GCS

    log.Fatal(http.ListenAndServe(":8080", handler.InitRouter())) // 启动 server
}
