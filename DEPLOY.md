# 部署文档

## 前提条件

### 1. 在 CPA 中创建 OpenAI 兼容提供商（必须提前完成）

**OpenCodeGoPool 通过 PUT 请求合并更新 CPA 的提供商列表，但不会在 CPA 侧创建提供商配置入口。必须先在 CPA 管理界面手动创建好对应的提供商，否则同步后 CPA 不会路由流量到该提供商。**

在 CPA 管理界面按如下参数新建一个 OpenAI 兼容提供商：

| 字段 | 值 |
|---|---|
| 名称（Provider Name） | 自定义，需与后续 `config.yaml` 中 `cpa.provider_name` 完全一致，默认为 `OpenCode Go` |
| Base URL | `https://opencode.ai/zen/go/v1` |
| API Key | 随意填写一个占位值（后续每次同步会自动覆盖为真实的账号 API Key 列表） |

创建完成后，记录下以下两个值，配置文件会用到：

- **管理端点（Endpoint）**：CPA 管理 API 的 URL，格式类似 `https://<your-cpa-host>/v0/management/openai-compatibility`
- **Bearer Token**：用于调用 CPA 管理 API 的鉴权 Token

### 2. 环境依赖

| 工具 | 版本要求 |
|---|---|
| Go | 1.21+ |
| Node.js | 18+ |

---

## 构建

```bat
# Windows 一键构建（前端 + 后端，产物为 build/opencode-pool-linux-amd64）
build.bat
```

也可以分步执行：

```bash
# 1. 构建前端
cd web
npm ci
npm run build
cd ..

# 2. 将前端产物复制到嵌入目录
# Windows
xcopy web\dist internal\frontend\dist\ /e /i /y
# Linux/macOS
cp -r web/dist internal/frontend/dist

# 3. 编译 Go 二进制（跨平台编译到 Linux）
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o build/opencode-pool ./cmd/server
```

---

## 配置

将 `config.example.yaml` 复制为 `config.yaml` 并按实际情况填写：

```yaml
server:
  address: ":8080"       # 监听地址
  password: "your-password"  # 管理界面登录密码

database:
  path: "./data/pool.db"  # SQLite 文件路径，目录需提前存在或有写权限

scraper:
  interval: "5m"   # 抓取所有账号配额和用量的间隔，默认 5 分钟
  timeout: "30s"   # 单次 HTTP 请求超时

cpa:
  endpoint: "https://<your-cpa-host>/v0/management/openai-compatibility"
  bearer_token: "<your-cpa-bearer-token>"
  provider_name: "OpenCode Go"   # 必须与在 CPA 中创建的提供商名称完全一致
  base_url: "https://opencode.ai/zen/go/v1"  # 同步给 CPA 的 Base URL，一般无需修改
```

> `cpa.provider_name` 和 `cpa.base_url` 会在首次启动时写入数据库作为默认值。之后可在管理界面的「设置」页修改，修改后数据库中的值优先于 `config.yaml`。

---

## Docker 部署

仓库内置 GitHub Actions 工作流（`.github/workflows/docker.yml`）：推送到 `main` 分支或打 `v*` 标签时自动构建 `linux/amd64` + `linux/arm64` 镜像并推送到 GHCR；Pull Request 只验证构建不推送。

镜像地址：

```
ghcr.io/handsomezhuzhu/opencodego_pool:latest      # main 分支最新
ghcr.io/handsomezhuzhu/opencodego_pool:1.2.3       # v1.2.3 标签（语义化版本）
```

### 使用 docker compose（推荐）

仓库根目录已提供 `docker-compose.yml`：

```bash
cp config.example.yaml config.yaml
# 编辑 config.yaml
docker compose up -d
```

- `config.yaml` 以只读方式挂载进容器，修改后执行 `docker compose restart` 生效
- SQLite 数据库保存在命名卷 `opencode-pool-data` 中（对应容器内 `/app/data`）

### 使用 docker run

```bash
docker run -d --name opencode-pool --restart unless-stopped \
  -p 8080:8080 \
  -v ./config.yaml:/app/config.yaml:ro \
  -v opencode-pool-data:/app/data \
  ghcr.io/handsomezhuzhu/opencodego_pool:latest
```

### 说明

- 镜像基于 Alpine，内置前端，以非 root 用户运行，监听 `8080` 端口
- GHCR 包默认为私有：在仓库 Packages 页面改为 Public 可匿名拉取，或在服务器上 `docker login ghcr.io`（使用带 `read:packages` 权限的 PAT）后拉取
- 自托管镜像也可本地构建：`docker build -t opencode-pool .`

---

## 部署与运行（裸机 / systemd）

将以下文件上传到服务器：

```
opencode-pool-linux-amd64   # 编译产物
config.yaml
```

确保 `data/` 目录存在（或 `config.yaml` 中指定的数据库路径所在目录存在）：

```bash
mkdir -p data
chmod +x opencode-pool-linux-amd64
./opencode-pool-linux-amd64
```

服务启动后会立即触发一次全量抓取，之后按 `scraper.interval` 定时抓取。

### 使用 systemd 管理进程

```ini
# /etc/systemd/system/opencode-pool.service
[Unit]
Description=OpenCode Go Pool Manager
After=network.target

[Service]
WorkingDirectory=/opt/opencode-pool
ExecStart=/opt/opencode-pool/opencode-pool-linux-amd64
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable --now opencode-pool
journalctl -u opencode-pool -f
```

---

## 首次使用

1. 打开 `http://<host>:8080`，用 `config.yaml` 中的 `server.password` 登录。
2. 进入「设置」页，确认 CPA 配置（Endpoint、Bearer Token、Provider Name、Base URL）与 CPA 管理界面中的提供商一致。
3. 进入「账号」页，点击「添加账号」，填写：
   - **Email**：opencode.ai 账号邮箱（仅作标识用）
   - **Cookie**：登录 opencode.ai 后从浏览器 DevTools 中复制完整的 Cookie 字符串
   - **Workspace ID**：工作区 ID，从 opencode.ai 工作区 URL 中获取（`/workspace/<workspace_id>/...`）
   - **API Key**：该账号对应的 API Key（`sk-...`），用于同步到 CPA；如果账号没有 API Key 则留空，该账号不会被同步到 CPA
4. 账号添加成功后，系统会立即触发一次抓取并同步 CPA。

---

## CPA 同步说明

每次以下操作发生时，系统会自动向 CPA 发起同步：

- 添加、修改、删除账号
- 启用或禁用账号
- 抓取发现某账号的 `limit_exceeded` 状态发生变化

同步逻辑：

1. 从 CPA 拉取当前完整的提供商列表（GET）
2. 找到名称匹配 `provider_name` 的提供商条目，将其 `api-key-entries` 替换为当前所有**状态为 active 且未超出配额限制**（`limit_exceeded = false`）且**填写了 API Key** 的账号
3. 将合并后的完整列表写回 CPA（PUT）

其他提供商的数据不受影响。

可在「设置」页查看最近一次同步的状态和时间，也可手动点击「立即同步」。

---

## 账号 Cookie 失效

抓取时若 Cookie 失效（HTTP 401/403 或被重定向到登录页），该账号状态会变为 `error`，需要重新从浏览器获取 Cookie 并在「账号」页更新。失效的账号不影响其他账号的抓取和 CPA 同步。
