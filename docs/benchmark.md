# benchmark

## 说明
参考[traefik benchmark](https://doc.traefik.io/traefik/v1.4/benchmarks/)方式进行

## 环境
### 虚拟机和测试机信息
4台腾讯云标准S6服务器，配置信息如下

- 8 GB RAM
- 4 CPU Cores
- 50 GB SSD
- Ubuntu 22.04 LTS 64-bit

其中一台作为gateway, 两台作为后端分别为whoami1, whoami2

最后一台作为测试机器

### 测试工具
1. 使用[wrk](https://github.com/wg/wrk)进行基准测试
2. 两台虚拟机使用[whoami](https://github.com/traefik/whoami)作为后端服务
3. traefik版本为2.4.x
4. nginx版本为: nginx/1.18.0 (Ubuntu)

### 通讯方式
vms之间全部通过腾讯云内网IP通讯, 腾讯控制台指明了内网带宽为6GBPS

### 测试命令
在测试vm上执行
```shell
wrk -t20 -c1000 -d60s --latency  http://{gatewayIP}:8001/bench
```

## 配置
4台服务器上执行如下命令
```shell
sysctl -w fs.file-max="9999999"
sysctl -w fs.nr_open="9999999"
sysctl -w net.core.netdev_max_backlog="4096"
sysctl -w net.core.rmem_max="16777216"
sysctl -w net.core.somaxconn="65535"
sysctl -w net.core.wmem_max="16777216"
sysctl -w net.ipv4.ip_local_port_range="1025       65535"
sysctl -w net.ipv4.tcp_fin_timeout="30"
sysctl -w net.ipv4.tcp_keepalive_time="30"
sysctl -w net.ipv4.tcp_max_syn_backlog="20480"
sysctl -w net.ipv4.tcp_max_tw_buckets="400000"
sysctl -w net.ipv4.tcp_no_metrics_save="1"
sysctl -w net.ipv4.tcp_syn_retries="2"
sysctl -w net.ipv4.tcp_synack_retries="2"
sysctl -w net.ipv4.tcp_tw_recycle="1"
sysctl -w net.ipv4.tcp_tw_reuse="1"
sysctl -w vm.min_free_kbytes="65536"
sysctl -w vm.overcommit_memory="1"
ulimit -n 9999999
```

whoami服务器
```shell
docker run -d -p 8000:80 --name iamfoo traefik/whoami
```

### Nginx
/etc/nginx/nginx.conf
```text
user ubuntu;
worker_processes auto;
worker_rlimit_nofile 200000;
pid /var/run/nginx.pid;

events {
    worker_connections 10000;
    use epoll;
    multi_accept on;
}

http {
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 300;
    keepalive_requests 10000;
    types_hash_max_size 2048;

    open_file_cache max=200000 inactive=300s;
    open_file_cache_valid 300s;
    open_file_cache_min_uses 2;
    open_file_cache_errors on;

    server_tokens off;
    dav_methods off;

    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    access_log /var/log/nginx/access.log combined;
    error_log /var/log/nginx/error.log warn;

    gzip off;
    gzip_vary off;

    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*.conf;
}
```

/etc/nginx/conf.d/whoami.conf
```text
upstream whoami {
    server {whami1IP}:8000;
    server {whami2IP}:8000;
    keepalive 300;
}

server {
    listen 8001;
    access_log off;
    error_log /dev/null crit;
    location / {
        proxy_pass http://whoami;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    proxy_set_header  X-Forwarded-Host $host;
    }
}
```

### traefik
/compose/traefik/docker-compose.yaml
```yaml
version: '3'

services:
  proxy:
    image: traefik:v2.4
    restart: always
    ports:
      - "8080:8080"
      - "8002:8002"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./traefik.toml:/etc/traefik/traefik.toml
      - ./conf:/dynamic/conf
```

/compose/traefik/traefik.toml
```toml
[api]
  insecure = true

[entryPoints]
  [entryPoints.web]
    address = ":8002"

[providers.docker]
  endpoint = "unix:///var/run/docker.sock"
  [providers.file]
    directory = "/dynamic/conf"
```

/compose/traefik/conf/whoami.toml

```toml
[http]
  [http.routers]
     [http.routers.whoami]
      rule = "PathPrefix(`/`)"
      service = "whoami-service"
      entryPoints = ["web"]
  [http.serversTransports.mytransport]
    maxIdleConnsPerHost = 100000

  [http.services]
    [http.services.whoami-service.loadBalancer]
    serversTransport = "mytransport"
      [[http.services.whoami-service.loadBalancer.servers]]
        url = "http://{whami1IP}:8000"
      [[http.services.whoami-service.loadBalancer.servers]]
        url = "http://{whami2IP}:8000"
```


### api gateway
/compose/api-gateawy/config.yaml
```yaml
log:
  level: info
  path: stdout
  access_log:
    enable: false

gateway:
  graceTimeOut: 50ns
  listenUDPTimeOut: 30s
```

/compose/api-gateawy/endpoints.yaml
```yaml
endpoints:
  - name: benchmark
    listen_port: 8003
    type: tcp
    routers:
      - type: http
        rule: "PathPrefix(`/`)"
        upstream:
          type: url
          paths:
            - http://{whami1IP}:8000/
            - http://{whami2IP}:8000/
          load_balance: robin
          maxIdleConnsPerHost: 100000
```


## 测试结果
```text
wrk -t20 -c1000 -d60s --latency  http://10.206.16.10:8001/bench
Running 1m test @ http://10.206.16.10:8001/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     6.81ms    3.66ms  52.36ms   70.19%
    Req/Sec     7.42k     1.72k   33.90k    64.32%
  Latency Distribution
     50%    6.80ms
     75%   10.08ms
     90%   10.56ms
     99%   14.40ms
  8864116 requests in 1.00m, 1.16GB read
Requests/sec: 147491.13
Transfer/sec:     19.83MB


wrk -t20 -c1000 -d60s --latency  http://10.206.16.10:8002/bench
Running 1m test @ http://10.206.16.10:8002/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    28.98ms    7.55ms 236.87ms   82.26%
    Req/Sec     1.74k   144.00     3.52k    75.11%
  Latency Distribution
     50%   28.41ms
     75%   31.92ms
     90%   35.82ms
     99%   50.51ms
  2081968 requests in 1.00m, 202.52MB read
Requests/sec:  34647.36
Transfer/sec:      3.37MB


wrk -t20 -c1000 -d60s --latency  http://10.206.16.10:8003/bench
Running 1m test @ http://10.206.16.10:8003/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    21.66ms    5.42ms 196.03ms   84.43%
    Req/Sec     2.33k   183.95     7.83k    83.10%
  Latency Distribution
     50%   21.38ms
     75%   23.80ms
     90%   26.17ms
     99%   36.40ms
  2783878 requests in 1.00m, 270.80MB read
Requests/sec:  46329.19
Transfer/sec:      4.51MB
```

## 结论
从上面测试结果可以看到nginx的性能领先非常多基本上是traefik和api_gateway的4-5倍，但是仍然在一个数量级

当前api_gateway的性能实际上比traefik还要高一些，当然可能是测试方法或者配置的问题，而且traefik的功能相对更加强大。

总体来说api_gateway的性能是可以保障的。
