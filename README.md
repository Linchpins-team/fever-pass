# 使用說明

Fever Pass 是一個數位體溫記錄系統，機構可使用此系統取代紙本自主體溫記錄表。

## License 
本網站由 Linchpins 團隊開發，以 GPL v3 釋出，可於 [GitHub](https://github.com/Linchpins-team/fever-pass) 取得原始碼。

## 開發者
林宏信：資料庫、後端與網頁模板
張旭誠：網頁設計

## 部署
安裝[Go](https://golang.org)

編譯：

	go build

建立設定檔：

	./fever-pass -init
	./fever-pass -conf /path/to/config.toml

接下來會詢問一系列問題：
留空為預設值。debug 模式使用 SQLite，儲存在 /tmp/gorm.sqlite，release 連線到本機的 MySQL，填寫資料庫名稱、使用者帳密後會嘗試連結，如果尚未建立可以選擇自動建立。最後設定管理員密碼。完成後設定檔會儲存於 config.toml，管理員密碼不會儲存於設定檔。

加密金鑰將會自動生成，儲存於 .env 中供下次存取。

可以透過 `quicktest.sh` 執行 curl 測試，go 測試使用 `go test`  
