DevTools Manager

概要
DevTools Manager は、トークン管理を簡単に行うためのアプリケーションです。ユーザーはトークンの追加、更新、削除、CSV ファイルのアップロードおよびダウンロードを行うことができます。

ローカル環境での実行方法
フロントエンド

1. フロントエンドディレクトリに移動します。

```
cd frontend
```

2. 依存関係をインストールします。

```
npm install
```

3. アプリケーションを起動します。

```
npm run dev
```

バックエンド

1. バックエンドディレクトリに移動します。

```
cd backend
```

2. アプリケーションを実行します。

```
go run main.go
```

.env ファイルの設定
.env ファイルをルートディレクトリに作成し、以下の内容を追加します。

env

```
REACT_APP_API_URL=http://localhost:8080
```

[デプロイ済みアプリ](https://frontend-service-632501898277.asia-east1.run.app)
デプロイ済みのアプリケーションは以下の URL でアクセスできます。

使用方法
トークン管理

### CSV ファイルのアップロード

1. アップロードする CSV ファイルを選択し、「アップロード」ボタンをクリックします。

### トークンの追加

1. フォームに必要な情報を入力し、「トークンを追加」ボタンをクリックします。

### トークンの更新

1. プロジェクト名を入力し、「トークンを取得」ボタンをクリックします。
2. フォームに表示されたトークンの情報を編集し、「トークンを更新」ボタンをクリックします。

### トークンの削除

1. 削除したいトークンの行番号を入力し、「トークンを削除」ボタンをクリックします。

### CSV ファイルとして保存

1. 「CSV として保存」ボタンをクリックし、ファイル名を入力して保存します。

[作業の記録](https://qiita.com/ryotaNS/items/7251435810e40afab9db)

