
# DEMO

👉 [線上 DEMO](https://chiaming.sdpmlab.org/social-app/)

# social-app

一個含前後端的社群貼文範例專案：
- 前端：React + CRA + TailwindCSS（支援 Markdown 編輯/預覽）
- 後端：Golang + Gin + Gorm（PostgreSQL）
- AI：以 OpenAI 相容 API 產生與優化貼文內容

<image src="https://raw.githubusercontent.com/alsk1369854/social-app/refs/heads/master/docs/images/c4-container.png" alt="c4-container.png">

## 快速開始（Docker）

1) 建立環境變數檔（Docker 用，專案根目錄）
- 複製 `.env.example` 為 `.env`，依需求修改。

2) 啟動服務
```bash
docker compose up -d
```

3) 開啟應用
- Web: http://localhost:${SERVER_PORT:-28080}
- Swagger: http://localhost:${SERVER_PORT:-28080}/swagger/index.html

> 預設會一併啟動 PostgreSQL（資料會掛載到 `./postgres-data`）。

## 本機開發（不走 Docker）

後端（Go 1.24+）：
1) 啟動資料庫（擇一）
- 使用專案內開發用 compose 啟 DB：
	```bash
	docker compose -f backend/docker-compose.yml up -d
	```
- 或自行啟動本機 PostgreSQL（設定與 `.env` 對應）。

2) 設定環境變數
- 複製 `backend/.env.example` 為 `backend/.env`（本機可將 `DB_HOST=127.0.0.1`、`DB_PORT=5432`）。

3) 啟動 Server（熱重載可選）
- 安裝 air（熱重載）：參考 `backend/README.md` 安裝 CLI 後執行
	```bash
	make debug
	```
- 或直接執行
	```bash
	cd backend
	go run .
	```

前端（Node 20+）：
```bash
cd frontend
npm ci
npm start
```
開發模式下，前端使用 CRA 的 `proxy`（指向 http://localhost:28080），因此可直接透過相對路徑呼叫 API。

## 環境變數

建立 `backend/.env`（可參考 `backend/.env.example`）
```env
# Server
DEBUG_MODE=true
SERVER_HOST=0.0.0.0
SERVER_PORT=28080
JWT_SECRET=<RUN openssl rand -base64 32>

# Database
DB_HOST=db           # Docker 模式為 db；本機開發可改 127.0.0.1
DB_PORT=5432
DB_NAME=social
DB_USER=admin
DB_PASSWORD=pg123456

# 初始帳號（啟動時自動建立/覆蓋 Admin，另會建立一組 Guest: temp@temp.com / temp@temp）
ADMIN_EMAIL=admin@admin.com
ADMIN_PASSWORD=admin@admin

# AI 提供者（OpenAI 相容）
OPENAI_API_KEY=<your-api-token>
OPENAI_BASE_URL=<your-api-base-url>
OPENAI_CHAT_MODEL=<your-chat-model-name>
```

前端可在生產環境提供下列變數（`frontend/.env.production` 或建置時注入）：
```env
REACT_APP_BASE_URL=/            # 前端路由與資源的 Base，預設 /
REACT_APP_API_BASE_URL=/api     # 後端 API Base，預設 /api
```

## 自訂 Base URL（部署於子路徑，如 /social-app）

若你希望網站掛在子目錄（例如 Nginx 的 /social-app）：
1) 調整前端 homepage 與建置環境
	 - 編輯 `frontend/package.json`
		 ```json
		 {
			 "homepage": "/social-app/public"
		 }
		 ```
	 - 設定 `frontend/.env.production`
		 ```env
		 REACT_APP_BASE_URL=/social-app
		 REACT_APP_API_BASE_URL=/social-app/api
		 ```
2) 反向代理路由對應
	 - 將 `/social-app/public` 與 `/social-app` 代理到後端服務；
	 - 將 `/social-app/api` 代理到後端的 `/api`。

> 後端會將前端建置產物掛載於 `/public`，根路由 `/` 直接回傳 `public/index.html`。

## API 與文件

- 後端 API BasePath：`/api`
- Swagger（啟動後）：`/swagger/index.html`
- Swagger 原始檔：`backend/docs/swagger.yaml`、`backend/docs/swagger.json`

主要路由（部分）：
- 使用者 / 貼文 / 留言 CRUD
- AI 內容生成功能：
	- POST `/api/ai/generate/text/create-post-content`
	- POST `/api/ai/generate/text/content-optimize`

授權
- 以 `Authorization` Header 帶入存取令牌（依後端中介層驗證）。

## 專案結構（節錄）

```
social-app/
	docker-compose.yml        # 一鍵啟動 backend + db
	backend/
		Dockerfile              # 建置：先建前端，再建後端，最後單一鏡像
		main.go                 # Gin 入口；/public 靜態檔、/api 路由、Swagger
		internal/               # 資料庫、路由、服務、模型…
		docs/                   # swag 產生之 API 文件
		docker-compose.yml      # 開發用：DB/管理工具（Adminer/pgAdmin）
	frontend/
		src/                    # React App（Markdown 編輯/預覽、AI 串接）
		public/
```
