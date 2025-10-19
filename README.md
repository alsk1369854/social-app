# social-app

## Workflow
<image src="https://raw.githubusercontent.com/alsk1369854/social-app/refs/heads/master/docs/images/c4-container.png" alt="c4-container.png">

## Deploy With Docker

### Optional: if your need custom base URL `social-app`

#### update `./frontend/package.json`
```json
{
 ...
 "homepage": "/social-app/public",
 ...
}
```

#### update `./frontend/.env.production`
```env
REACT_APP_BASE_URL=/social-app
REACT_APP_API_BASE_URL=/social-app/api
```

### create and update `./backend/.env`
```env
SERVER_PORT=28080
JWT_SECRET=<RUN openssl rand -base64 32>

DB_NAME=social
DB_USER=admin
DB_PASSWORD=pg123456

ADMIN_EMAIL=admin@admin.com
ADMIN_PASSWORD=admin@admin

OPENAI_API_KEY=<your-api-token>
OPENAI_BASE_URL=<your-api-base-url>
OPENAI_CHAT_MODEL=<your-chat-model-name>
```

### docker compose
```bash
docker compose up
```

### open web app `http://localhost:28080`
