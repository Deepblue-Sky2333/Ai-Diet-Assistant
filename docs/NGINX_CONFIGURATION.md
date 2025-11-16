# Nginx 配置指南

本文档提供了 AI Diet Assistant 后端 API 服务的 Nginx 配置示例和最佳实践。

## 概述

AI Diet Assistant 后端是一个纯 API 服务，不包含前端代码。在生产环境中，建议使用 Nginx 作为反向代理，处理以下功能：

- **反向代理**：将请求转发到后端 API 服务
- **CORS 处理**：统一处理跨域请求
- **SSL/TLS 终止**：处理 HTTPS 加密
- **负载均衡**：支持多实例部署
- **静态资源缓存**：优化性能
- **请求限流**：防止滥用

## 基础配置

### 最小配置示例

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    # 代理到后端 API
    location / {
        proxy_pass http://localhost:9090;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 完整配置示例

```nginx
# 上游服务器配置
upstream diet_assistant_backend {
    # 单实例配置
    server localhost:9090;
    
    # 多实例负载均衡配置（可选）
    # server localhost:9090 weight=1;
    # server localhost:9091 weight=1;
    # server localhost:9092 weight=1;
    
    # 健康检查
    keepalive 32;
}

# HTTP 服务器配置
server {
    listen 80;
    server_name api.yourdomain.com;

    # 访问日志
    access_log /var/log/nginx/diet-assistant-access.log;
    error_log /var/log/nginx/diet-assistant-error.log;

    # 请求体大小限制（用于文件上传）
    client_max_body_size 10M;

    # CORS 配置
    add_header 'Access-Control-Allow-Origin' '$http_origin' always;
    add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
    add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Requested-With' always;
    add_header 'Access-Control-Allow-Credentials' 'true' always;
    add_header 'Access-Control-Max-Age' '3600' always;

    # 处理 OPTIONS 预检请求
    if ($request_method = 'OPTIONS') {
        add_header 'Access-Control-Allow-Origin' '$http_origin' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Requested-With' always;
        add_header 'Access-Control-Allow-Credentials' 'true' always;
        add_header 'Access-Control-Max-Age' '3600' always;
        add_header 'Content-Type' 'text/plain; charset=utf-8';
        add_header 'Content-Length' '0';
        return 204;
    }

    # 代理到后端 API
    location / {
        proxy_pass http://diet_assistant_backend;
        
        # 代理头设置
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
        
        # 缓冲设置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
        
        # WebSocket 支持（如果需要）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 健康检查端点（不需要 CORS）
    location /health {
        proxy_pass http://diet_assistant_backend/health;
        proxy_set_header Host $host;
        access_log off;
    }
}
```

## HTTPS 配置

### 使用 Let's Encrypt 证书

```bash
# 安装 Certbot
sudo apt-get install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d api.yourdomain.com

# 自动续期
sudo certbot renew --dry-run
```

### HTTPS 配置示例

```nginx
# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name api.yourdomain.com;
    
    # 重定向到 HTTPS
    return 301 https://$server_name$request_uri;
}

# HTTPS 服务器配置
server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL 证书配置
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    
    # SSL 安全配置
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # HSTS（可选，但推荐）
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # 访问日志
    access_log /var/log/nginx/diet-assistant-access.log;
    error_log /var/log/nginx/diet-assistant-error.log;

    # 请求体大小限制
    client_max_body_size 10M;

    # CORS 配置
    add_header 'Access-Control-Allow-Origin' '$http_origin' always;
    add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
    add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Requested-With' always;
    add_header 'Access-Control-Allow-Credentials' 'true' always;
    add_header 'Access-Control-Max-Age' '3600' always;

    # 处理 OPTIONS 预检请求
    if ($request_method = 'OPTIONS') {
        add_header 'Access-Control-Allow-Origin' '$http_origin' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-Requested-With' always;
        add_header 'Access-Control-Allow-Credentials' 'true' always;
        add_header 'Access-Control-Max-Age' '3600' always;
        add_header 'Content-Type' 'text/plain; charset=utf-8';
        add_header 'Content-Length' '0';
        return 204;
    }

    # 代理到后端 API
    location / {
        proxy_pass http://localhost:9090;
        
        # 代理头设置
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 健康检查端点
    location /health {
        proxy_pass http://localhost:9090/health;
        proxy_set_header Host $host;
        access_log off;
    }
}
```

## 高级配置

### 请求限流

```nginx
# 在 http 块中定义限流区域
http {
    # 限制每个 IP 每秒 10 个请求
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
    
    # 限制每个 IP 的并发连接数
    limit_conn_zone $binary_remote_addr zone=conn_limit:10m;
}

# 在 server 块中应用限流
server {
    # ... 其他配置 ...
    
    location / {
        # 应用请求限流（允许突发 20 个请求）
        limit_req zone=api_limit burst=20 nodelay;
        
        # 应用连接限流（每个 IP 最多 10 个并发连接）
        limit_conn conn_limit 10;
        
        proxy_pass http://localhost:9090;
        # ... 其他代理配置 ...
    }
}
```

### 负载均衡

```nginx
upstream diet_assistant_backend {
    # 轮询（默认）
    server localhost:9090;
    server localhost:9091;
    server localhost:9092;
    
    # 或使用 IP Hash（同一客户端总是访问同一服务器）
    # ip_hash;
    
    # 或使用最少连接
    # least_conn;
    
    # 健康检查
    keepalive 32;
}

server {
    # ... 其他配置 ...
    
    location / {
        proxy_pass http://diet_assistant_backend;
        # ... 其他代理配置 ...
    }
}
```

### 缓存配置

```nginx
# 在 http 块中定义缓存路径
http {
    proxy_cache_path /var/cache/nginx/diet-assistant 
                     levels=1:2 
                     keys_zone=api_cache:10m 
                     max_size=1g 
                     inactive=60m 
                     use_temp_path=off;
}

server {
    # ... 其他配置 ...
    
    # 缓存静态数据（如食材列表）
    location ~* ^/api/v1/(foods|settings) {
        proxy_pass http://localhost:9090;
        
        # 启用缓存
        proxy_cache api_cache;
        proxy_cache_valid 200 5m;
        proxy_cache_key "$scheme$request_method$host$request_uri";
        
        # 缓存头
        add_header X-Cache-Status $upstream_cache_status;
        
        # 其他代理配置
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # 不缓存动态数据
    location / {
        proxy_pass http://localhost:9090;
        proxy_no_cache 1;
        proxy_cache_bypass 1;
        # ... 其他代理配置 ...
    }
}
```

### 安全加固

```nginx
server {
    # ... 其他配置 ...
    
    # 隐藏 Nginx 版本号
    server_tokens off;
    
    # 安全头
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    
    # 限制请求方法
    if ($request_method !~ ^(GET|POST|PUT|DELETE|OPTIONS)$) {
        return 405;
    }
    
    # 阻止常见攻击路径
    location ~ /\. {
        deny all;
        access_log off;
        log_not_found off;
    }
    
    # ... 其他配置 ...
}
```

## 部署步骤

### 1. 安装 Nginx

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install nginx

# CentOS/RHEL
sudo yum install nginx

# macOS
brew install nginx
```

### 2. 创建配置文件

```bash
# 创建配置文件
sudo vim /etc/nginx/sites-available/diet-assistant

# 创建符号链接
sudo ln -s /etc/nginx/sites-available/diet-assistant /etc/nginx/sites-enabled/

# 或者直接编辑主配置文件
sudo vim /etc/nginx/nginx.conf
```

### 3. 测试配置

```bash
# 测试配置文件语法
sudo nginx -t

# 如果测试通过，输出：
# nginx: configuration file /etc/nginx/nginx.conf test is successful
```

### 4. 重启 Nginx

```bash
# 重启 Nginx
sudo systemctl restart nginx

# 或重新加载配置（不中断服务）
sudo systemctl reload nginx

# 查看状态
sudo systemctl status nginx
```

### 5. 配置防火墙

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow 'Nginx Full'
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
sudo firewall-cmd --reload
```

### 6. 验证配置

```bash
# 测试 HTTP 访问
curl http://api.yourdomain.com/health

# 测试 HTTPS 访问
curl https://api.yourdomain.com/health

# 测试 CORS
curl -H "Origin: https://yourdomain.com" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Authorization, Content-Type" \
     -X OPTIONS \
     https://api.yourdomain.com/api/v1/foods
```

## 故障排查

### 查看日志

```bash
# 查看访问日志
sudo tail -f /var/log/nginx/diet-assistant-access.log

# 查看错误日志
sudo tail -f /var/log/nginx/diet-assistant-error.log

# 查看 Nginx 主日志
sudo tail -f /var/log/nginx/error.log
```

### 常见问题

#### 1. 502 Bad Gateway

**原因**：后端服务未启动或无法连接

**解决方案**：
```bash
# 检查后端服务状态
sudo systemctl status diet-assistant

# 启动后端服务
sudo systemctl start diet-assistant

# 检查端口是否监听
sudo netstat -tlnp | grep 9090
```

#### 2. CORS 错误

**原因**：CORS 配置不正确

**解决方案**：
- 确保 `add_header` 指令包含 `always` 参数
- 确保 OPTIONS 请求正确处理
- 检查 `Access-Control-Allow-Origin` 是否正确设置

#### 3. 413 Request Entity Too Large

**原因**：请求体超过限制

**解决方案**：
```nginx
# 增加请求体大小限制
client_max_body_size 20M;
```

#### 4. 504 Gateway Timeout

**原因**：后端响应超时

**解决方案**：
```nginx
# 增加超时时间
proxy_connect_timeout 120s;
proxy_send_timeout 120s;
proxy_read_timeout 120s;
```

## 性能优化

### 1. 启用 Gzip 压缩

```nginx
http {
    gzip on;
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript 
               application/json application/javascript application/xml+rss 
               application/rss+xml font/truetype font/opentype 
               application/vnd.ms-fontobject image/svg+xml;
}
```

### 2. 启用 HTTP/2

```nginx
server {
    listen 443 ssl http2;
    # ... 其他配置 ...
}
```

### 3. 优化缓冲区

```nginx
server {
    # ... 其他配置 ...
    
    location / {
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
        # ... 其他代理配置 ...
    }
}
```

### 4. 启用 Keepalive

```nginx
upstream diet_assistant_backend {
    server localhost:9090;
    keepalive 32;
}

server {
    # ... 其他配置 ...
    
    location / {
        proxy_pass http://diet_assistant_backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        # ... 其他代理配置 ...
    }
}
```

## 监控和日志

### 日志格式

```nginx
http {
    # 自定义日志格式
    log_format api_log '$remote_addr - $remote_user [$time_local] '
                       '"$request" $status $body_bytes_sent '
                       '"$http_referer" "$http_user_agent" '
                       '$request_time $upstream_response_time';
    
    access_log /var/log/nginx/diet-assistant-access.log api_log;
}
```

### 日志轮转

```bash
# 创建日志轮转配置
sudo vim /etc/logrotate.d/diet-assistant

# 内容：
/var/log/nginx/diet-assistant-*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 www-data adm
    sharedscripts
    postrotate
        [ -f /var/run/nginx.pid ] && kill -USR1 `cat /var/run/nginx.pid`
    endscript
}
```

## 最佳实践

1. **使用 HTTPS**：生产环境必须使用 HTTPS
2. **配置 CORS**：在 Nginx 层统一处理 CORS，后端不需要处理
3. **启用限流**：防止 API 滥用
4. **配置健康检查**：确保后端服务可用
5. **日志监控**：定期检查日志，及时发现问题
6. **定期更新**：保持 Nginx 和 SSL 证书更新
7. **备份配置**：定期备份 Nginx 配置文件
8. **测试配置**：修改配置后务必测试

## 参考资源

- [Nginx 官方文档](https://nginx.org/en/docs/)
- [Let's Encrypt 文档](https://letsencrypt.org/docs/)
- [Mozilla SSL 配置生成器](https://ssl-config.mozilla.org/)
- [Nginx 性能优化指南](https://www.nginx.com/blog/tuning-nginx/)

