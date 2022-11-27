package main

import (
	"context"
	"fmt"
	"os"

	"dagger.io/dagger"
)

func main() {
	ctx := context.Background()

	// クライアントを初期化して Dagger Engine に接続する
	// dagger.WithLogOutput でログの出力先を指定できる
	client, err := dagger.Connect(ctx, dagger.WithLogOutput(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// ビルドマトリックス
	oses := []string{"linux", "darwin"}
	arches := []string{"amd64", "arm64"}

	// Docker イメージを取得する
	container := client.Container().From("golang:1.19")

	// カレントディレクトリをコンテナにマウントする
	src := client.Host().Directory(".")
	container = container.
		WithMountedDirectory("/src", src).
		WithWorkdir("/src")

	// テストを実行
	container = container.WithExec([]string{"go", "test", "-v", "./..."})

	// ビルド配置用の空のディレクトリを用意
	outputs := client.Directory()
	for _, goos := range oses {
		for _, goarch := range arches {
			// ビルド先のディレクトリ
			path := fmt.Sprintf("build/%s/%s/", goos, goarch)

			// GOOS, GOARCH 環境変数を設定
			container = container.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch)

			// ビルド
			container = container.WithExec([]string{"go", "build", "-o", path})

			// アウトプットにビルド先のディレクトリを追加
			outputs = outputs.WithDirectory(path, container.Directory(path))
		}
	}

	// パイプラインを実行する
	// Export でビルド先のディレクトリをホストに書き込む
	if _, err := outputs.Export(ctx, "."); err != nil {
		panic(err)
	}
}
