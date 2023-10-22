# Enablement Bootcamp for Gopherizing (2023.11)
画像変換バッチの作成

## Setting
### 入力用と出力用のGCSバケットの設定
1. 以下URLに従って2つバケットを作成します. その際出力用のバケットはインターネットに公開できるようにします. 

	https://cloud.google.com/storage/docs/creating-buckets?hl=ja

2. 以下URLに従ってサービスアカウントのキーを作成します. キーをjsonで取得し **credentials.json** と命名しルートディレクトリに追加します. 

	https://cloud.google.com/iam/docs/creating-managing-service-account-keys?hl=ja

3. バケットへの権限を追加します. サービスアカウントに **Storage レガシーオブジェクト オーナー** を追加. 出力用バケットには加えて, AllUsersに **Storage レガシーオブジェクト 読み取り** を追加します. 

4. .envファイルにそれぞれのバケット名を追加します. .env.sampleにしたがって追加してください. 

### mysqlのデータベースの設定
1. ローカルにmysqlデータベースを作成します. 作成したデータベースの情報を.envファイルに記載します. 

2. マイグレーションを行います. 以下のコマンドで実行できます. 

		$ migrate -database="postgres://ユーザー名:パスワード@ホスト名:ポート番号/データベース名?sslmode=disable" -path=database/migrations/ up 1

## demo

1. cmd/server, cmd/clientでそれぞれ以下を実行し, サーバー, クライアントを起動します. 

		go run main.go

2. cmd/client のコマンドラインでサーバーの操作を行うことができます. 起動すると以下の画面が表示されます. 

		start gRPC Client.
		1: create image
		2: convert images
		3: image list
		4: exit
		please enter > 

	それぞれ以下のような機能を持ちます. 

	image list：全ての登録された情報を取得, 表示

	convert images：データベースの情報に従ってバッチ処理を行う

	create image：画像情報の新規作成

3. create imageで画像データを追加します. 以下の手順でclient/img内のファイル名を指定し, 入力することで画像のデータが追加されます. 

		create image
		image file path: 7.png
		resize width percent: 50
		resize height percent: 50
		encode format (1: JPEG or 2: PNG): 1

		done

4. list imageで画像データを確認します. 以下のように情報が取得できます. 

		image 1: id:1  object_name:"1.png"  resize_width_percent:50  resize_height_percent:50  encode_format:JPEG  status:Waiting

5. convert imageでバッチ処理を行います. これを行うことで指定のサイズ, フォーマットに変更できます. 

6. list imageで画像データの変換を確認します. 変換すると以下のように出力されます

		image 7: id:7  object_name:"7.png"  resize_width_percent:50  resize_height_percent:50  encode_format:JPEG  status:Succeeded  converted_image_url:<converted_image_url>

	converted_image_urlをctrl clickすることで画像がダウンロードできます. 
