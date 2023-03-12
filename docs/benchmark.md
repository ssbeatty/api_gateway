# benchmark

## 说明
参考[traefik benchmark](https://doc.traefik.io/traefik/v1.4/benchmarks/)方式进行

## 环境
### 虚拟机和测试机信息
3台腾讯云标准S6服务器，配置信息如下

- 8 GB RAM
- 4 CPU Cores
- 50 GB SSD
- Ubuntu 22.04 LTS 64-bit
- 带宽 20M

测试电脑
- windows10 wsl
- 16 CPU Cores
- 32 GB RAM

### 测试工具
1. 使用[wrk](https://github.com/wg/wrk)进行基准测试
2. 两台虚拟机使用[whoami](https://github.com/traefik/whoami)作为后端服务
3. traefik版本为2.4.x
4. nginx版本为: nginx/1.18.0 (Ubuntu)

### 通讯方式
- gateway和whami vms通过腾讯云内网IP通讯
- 测试机和gateway通过公网IP通讯

### 测试命令
```shell
wrk -t20 -c1000 -d60s --latency  http://{gatewayIP}:8001/bench
```

## 配置
网关服务器
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


  [http.services]
    [http.services.whoami-service.loadBalancer]
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
```


## 测试结果
```text
wrk -t20 -c1000 -d60s --latency  http://43.137.1.29:8001/bench
Running 1m test @ http://43.137.1.29:8001/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   126.55ms  283.35ms   2.00s    88.98%
    Req/Sec   179.00    114.49     1.49k    74.47%
  Latency Distribution
     50%   11.48ms
     75%   22.19ms
     90%  447.44ms
     99%    1.55s
  212250 requests in 1.00m, 28.54MB read
  Socket errors: connect 0, read 0, write 0, timeout 2152
Requests/sec:   3531.62
Transfer/sec:    486.29KB



wrk -t20 -c1000 -d60s --latency  http://43.137.1.29:8002/bench
Running 1m test @ http://43.137.1.29:8002/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   126.63ms  281.33ms   2.00s    89.03%
    Req/Sec   181.74    112.39     1.38k    73.25%
  Latency Distribution
     50%   11.67ms
     75%   36.88ms
     90%  447.70ms
     99%    1.55s
  215431 requests in 1.00m, 20.96MB read
  Socket errors: connect 0, read 0, write 0, timeout 1972
Requests/sec:   3585.38
Transfer/sec:    357.14KB


wrk -t20 -c1000 -d60s --latency  http://43.137.1.29:8003/bench
Running 1m test @ http://43.137.1.29:8003/bench
  20 threads and 1000 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   126.83ms  281.49ms   2.00s    89.01%
    Req/Sec   181.23    110.40     1.08k    70.75%
  Latency Distribution
     50%   11.48ms
     75%   30.65ms
     90%  447.82ms
     99%    1.54s
  215319 requests in 1.00m, 20.95MB read
  Socket errors: connect 0, read 1, write 0, timeout 2175
Requests/sec:   3584.76
Transfer/sec:    357.08KB
```

## 结论
从上面测试结果可以看到基本上与traefik的请求数持平，api_gateway的性能可以保障。


> 也可能是网关没有达到对应的瓶颈。 之后会使用内网vm继续测试。