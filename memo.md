# Station07
## API仕様書の確認
- docs/openapi.yaml に書いてある仕様を見て、どのような値が必要か確認する。
### post /todos
- POSTメソッド: サーバーにデータを送信し、新しいTODOを作成するために使用
- 以下で示されるフィールドは、CreateTODORequest 構造体に反映される。
```yaml
post:
  summary: Create TODO   # エンドポイントの簡単な説明
  requestBody:           # エンドポイントに送信されるデータの形式と内容を定義
    content:
      application/json:  # リクエストボディがJSON形式であることを示している。
        schema:          # 具体的なJSONオブジェクトの形式
          type: object   # リクエストボディはJSONオブジェクトである。
          properties:    # リクエストボディのJSONオブジェクトが持つべきプロパティ（キーとそのデータ型）をリストアップ
            subject:          # フィールド名
              type: string    # subjectフィールドは文字列型
              required: true  # 必須
            description:
              type: string
              required: false # 必須ではない
```
## CreateTODORequest の定義
- 上記のyamlファイルより、利用者から受け取る値は`subject`と`description`であるので、以下のように定義。
    - model/todo.go
```go
CreateTODORequest struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}
```
## CreateTODOResponse の定義
- TODOをデータベースに保存した後に、そのTODOを利用者に返すために、`CreateTODOResponse`構造体をTODO構造体を含む形で定義
    - model/todo.go
```go
CreateTODOResponse struct {
		TODO TODO `json:"todo"`
	}
```

# Station10
- この実装により、ErrNotFound エラーが発生した際には、どのようなリソースが、どのIDに対して見つからなかったのかを具体的に示すエラーメッセージを返すことが可能になる。