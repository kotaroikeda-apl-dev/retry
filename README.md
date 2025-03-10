## **実行方法**

```sh
go run cmd/basic/main.go # リトライの基本を実装
```

## **学習ポイント**

1. **`retry`** の仕組みを学び、最大リトライ回数や待機時間を制御する方法を習得できる。
2. **`fmt.Errorf("%w", err)`** を使ってエラーをラップし、**`errors.Is()`** でリトライ可能なエラーを判定できる。

## 作成者

- **池田虎太郎** | [GitHub プロフィール](https://github.com/kotaroikeda-apl-dev)
