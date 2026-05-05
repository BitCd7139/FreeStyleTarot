## ⚒️ 线上部署：

### 基础设施：
- 采用 DigitalOcean App Platform 的 PaaS 服务进行部署，实现无服务器化容器运行。

### DNS层面：
- 使用 Cloudflare 进行 DNS 管理，提供域名解析和安全防护。
- 通过 Cloudflare 的 CDN 加速访问速度，提升用户体验。

### 安全性：
- HTTPs: 通过 Cloudflare 提供的 SSL/TLS 证书实现全站 HTTPS 加密，确保数据传输安全。
- DDoS防护: Cloudflare 的 DDoS 防护功能帮助抵御。

### 部署流：
- 使用 Git 进行版本控制，代码托管在 GitHub 上。
- 配置 DigitalOcean App Platform 与 GitHub 仓库连接，实现代码推送后自动构建和部署。
- 在 DigitalOcean App Platform 上配置环境变量，确保敏感信息（如 API 密钥）安全存储。
- 配置./do/app.yaml 文件，定义应用的构建和运行环境，包括部署区域、资源配置等。

### 监控和维护：
- 使用 DigitalOcean 提供的监控工具，实时监控应用的性能和资源
- 定期检查 Cloudflare 的安全报告，确保没有异常流量或攻击行为。