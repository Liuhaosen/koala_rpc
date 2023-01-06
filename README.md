# koala_rpc
grpc微服务框架
etcd作为服务注册中心

使用调用链方式执行中间件
  中间件类型:
    服务监控中间件: prometheus
    限流器中间件:  rate令牌桶限流器
    日志中间件: 自定义手写log作为日志记录, 分为console控制台打印日志和文件日志, 文件日志又加了access访问日志
    分布式追踪中间件: jaeger,zipkin


使用:
  1. 服务端代码生成:
    cd tools
    cd koala
    mkdir output
    go build
    cd output
    cp ..\koala.exe .
    cp ..\hello.proto .
    koala.exe -s -f .\hello.proto
    cd main 
    go build 
    main.exe
    
  2. 客户端代码生成
    cd tools
    cd koala
    mkdir client_example
    go build
    cd client_example
    cp ..\koala.exe .
    cp ..\hello.proto .
    koala.exe -c -f .\hello.proto
    go run main.go
    
  
  
  
