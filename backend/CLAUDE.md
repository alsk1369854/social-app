### 專案背景
- 此專案為一個社交平台的後段伺服器
- 此專案已 `golang` 作為主要使用程式語言
- 採用 [Gin Web Framework](https://gin-gonic.com/en/docs/) 作為 server 開發框架
- 採用 [GORM](https://gorm.io/docs/) 進行資料庫數據操作
- 採用 [Swaggo](https://github.com/swaggo/swag?tab=readme-ov-file#examples) 撰寫 API 文件
- 採用 [PlantUML](https://plantuml.com/zh/) 繪製系統圖
- 採用 `github.com/pkg/errors` 庫進行錯誤包裝，以利進行程式錯誤位置追蹤
- 採用 `Singleton pattern` 建構專案模組
- 每項功能都需撰寫單元測試


### 專案結構
- `./main.go` 程式入口
- `./.env` 開發環境參數配置
- `./docs/uml/<uml-type>/<diagram-name>.pu` 為針對功能的 uml 圖，採用 `PlantUML` 語法撰寫
- `./internal/routers/` 放置 API 對應的 Router `.go` 文件，負責 API 的輸入數據驗證與響應數據處裡
- `./internal/handlers/` 放置 API 對應的處裡邏輯，向下調用 `services` 中的方法完成業務邏輯
- `./internal/services/` 放置業務的數據處裡邏輯，在此維護 **Database Translation** 操作
- `./internal/repositories/` 放置資料庫訪問邏輯，在此維護對應資料表的 **CRUD** 方法
- `./internal/database/` 放置資料庫連線與 **GORM** 資料表結構綁定等操作
- `./internal/models/` 放置 API 接收參數與響應參數的 **struct**
- `./internal/pkg/<util-name>` 放置所需的工具包，提取重複性程式碼至對應的工具包中進行操作
- `./internal/middlewares/` 放置 API 中間層處裡函數
- `./PRPs/` 放置所有產品需求提示詞


### AI 行為規則
- **切勿假設上下文缺失。如有疑問，請提出問題**
- **切勿對函式庫或函數產生幻想**，僅使用已知且經過驗證的軟體包
- **在程式碼或測試中引用檔案路徑和模組名稱之前，**務必確認它們**存在**
- **切勿刪除或覆寫現有程式碼**，除非明確指示屬於任務的一部分