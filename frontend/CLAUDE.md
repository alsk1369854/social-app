### 專案背景
- 此專案為一個社交平台的單頁式 Web UI
- 此專案以 `typescript` 作為主要使用程式語言
- 採用 [React](https://react.dev/) 作為專案開發框架
- 採用 [tailwindcss](https://tailwindcss.com/) 進行資料庫數據操作
- 採用 [PlantUML](https://plantuml.com/zh/) 繪製系統圖

### 專案結構
- `./public/` 公開資源目錄
- `./docs/backend-api/` 後端 API Swagger 文件，包含了。1. API 路徑描述: **swagger.json**; 2. API 數據類型描述: **swagger.yaml**。你也可以直接訪問後端 [Swagger](http://localhost:28080/swagger/index.html) 頁面
- `./src/` 開發目錄
- `./src/apis/` 放置後端 API 串接方法
- `./src/apis/models/` 放置 `./src/apis/` 中使用到的接口數據模型定義
- `./src/components/` 放置專案通用型組件
- `./src/hooks/` 放置專案回調工具
- `./src/models/` 放置專案使用的數據模型定義，不要在這邊定義後端 API 接口使用的數據模型
- `./src/pages/` 放置專案頁面
- `./src/pages/layouts/` 放置頁面布局 
- `./PRPs/` 放置所有產品需求提示詞

### AI 行為規則
- **切勿假設上下文缺失。如有疑問，請提出問題**
- **切勿對函式庫或函數產生幻想**，僅使用已知且經過驗證的軟體包
- **在程式碼或測試中引用檔案路徑和模組名稱之前，**務必確認它們**存在**
- **切勿刪除或覆寫現有程式碼**，除非明確指示屬於任務的一部分