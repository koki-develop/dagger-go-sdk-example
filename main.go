package main

import (
	"context"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// クライアントを初期化
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// Docker イメージを取得する
	container := client.Container().From("golang:1.19")

	// カレントディレクトリをコンテナにマウントする
	src := client.Host().Directory(".")
	container = container.
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	// ビルド先のディレクトリ
	path := "build/"

	// 実行するコマンドを設定する
	container = container.
		WithExec([]string{"go", "test", "-v", "./..."}).
		WithExec([]string{"go", "build", "-o", path})

	// パイプラインを実行する
	// Export でビルド先のディレクトリをホストに書き込む
	output := container.Directory(path)
	if _, err := output.Export(ctx, path); err != nil {
		panic(err)
	}
}
