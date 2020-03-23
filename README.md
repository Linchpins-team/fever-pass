# 使用說明

Fever Pass 是一個數位體溫記錄系統，可使用此系統取代紙本體溫記錄表。

## License 
本網站由 Linchpins 團隊開發，以 GPL v3 釋出。

## 開發者
林宏信：資料庫、後端與網頁模板

張旭誠：網頁設計、文檔

## 部署
安裝[Go](https://golang.org)、[SQL（MariaDB）](https://mariadb.org/download)

編譯：

	go build

建立設定檔：

	./fever-pass -init

依提示訊息完成設定：
	
	server base url (default http://localhost:8080): 
	server listen port (default: 8080): 
	database mode debug/release (default: debug): release
	MySQL database setup. Please create database before configuration
	database name: name
	database user: user
	database password: 0000
	Cannot connect to database: Error 1698: Access denied for user 'user'@'localhost', initial it now? (y/n) y
	admin password: 0000
	configurations has generated at config.toml

	./fever-pass -conf /path/to/config.toml

加密金鑰將會自動生成，儲存於 .env 中供下次存取。

可以透過 `quicktest.sh` 執行 curl 測試，go 測試使用 `go test`  

---

若需要代為部署，請來信 linchpins-team@protonmail.com 商談。
