# 使用說明

Fever Pass 是一個數位體溫記錄系統，可使用此系統取代紙本體溫記錄表。

## License 
本網站由 Linchpins 團隊開發，以 GPL v3 釋出。

## 開發者
林宏信：資料庫、後端與網頁模板

張旭誠：網頁設計、文檔

## 依賴

- [MariaDB](https://mariadb.org/download) (or MySQL)

## 安裝

從 [Release](https://github.com/Linchpins-team/fever-pass/releases) 下載對應平台之安裝包，將其解壓縮。
### 從原始碼安裝
安裝[Go](https://golang.org)（1.13 以上）、

編譯：

	$ go build

## 佈署

### 設定檔
設定預設設定檔，設定檔會存在 config.toml，可以用 -conf 自訂設定檔位置：

	$ ./fever-pass -g

編輯設定檔：

```
Mode = "release" // release 使用 MySQL，debug 使用 sqlite 並存於 /tmp 中

[Site]
  Name = "Fever Pass"
  Icon = "/static/img/logo.png"

[Server]
  Base = "http://localhost:8080"
  Port = 8080

[Database]
  Host = "localhost"
  Name = "fever_pass"
  User = "fever_pass_user"
  Password = "password"
```

### 建立資料庫

使用 -init 可產生建立資料庫 SQL 語句：

```
$ ./fever-pass -init
Copy the following code to sql.

CREATE DATABASE IF NOT EXISTS fever_pass ;
DROP USER IF EXISTS 'fever_pass_user'@'localhost';
FLUSH PRIVILEGES;
CREATE USER 'fever_pass_user'@'localhost' IDENTIFIED BY 'password'; 
GRANT ALL PRIVILEGES ON fever_pass . * TO 'fever_pass_user'@'localhost'; 
FLUSH PRIVILEGES;
```

將其貼入具備權限的 SQL 中，即可建立資料庫，也可自行手動修改。

### 設定管理員密碼

使用 -s 設定 admin 密碼，注意：密碼將會明文顯示。

```
./fever-pass -s
admin password: password
```

事後若是忘記密碼，也可用相同方式修改。

### 啟動
	
	./fever-pass

加密金鑰將會自動生成，儲存於 .env 中供下次存取。

## 建置測試環境

測試帳號資料在 testdata 之中，可以使用 -mock 選項建立幫每個帳號填入假體溫資料

	$ ./fever-pass -mock

---

若需要代為或協助佈署，請來信 linchpins-team@protonmail.com 商談。
