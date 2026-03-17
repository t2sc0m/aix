# aix — AI eXchange CLI

コンテキストファイル注入によりAI CLIツールへタスクを委任するCLI。

```
aix "このコードのバグを見つけて" -c expert-prompt.md -f target.go
```

## なぜ作ったか

MCPベースのAI委任はオーバーヘッドが大きい — サーバー設定、プロトコルネゴシエーション、ツール定義によるコンテキストウィンドウの消費。`aix`はこれを直接サブプロセス呼び出しに置き換える。Expertプロンプトはプロトコル層ではなくコンテキストファイル（`-c`フラグ）で注入。

## インストール

```bash
go install github.com/t2sc0m/aix@latest
```

前提: [Codex CLI](https://github.com/openai/codex) がインストール・認証済みであること。

## 使い方

```bash
# 基本プロンプト
aix "このエラーを説明して: connection refused"

# Expertコンテキストファイル（MCP developer-instructionsの代替）
aix "このプランをレビューして" -c prompts/plan-reviewer.md -f plan.md

# 複数ファイル添付
aix "セキュリティの問題を見つけて" -c prompts/security-analyst.md -f auth.go -f handler.go

# stdin パイプ
echo "なぜ遅い？" | aix -f slow-query.sql

# モデル指定
aix "分析して" -m o3

# サンドボックスモード
aix "このバグを修正して" -s workspace-write -f broken.go

# ステータス確認
aix status
```

## フラグ

| フラグ | 短縮 | デフォルト | 説明 |
|--------|------|------------|------|
| `--context` | `-c` | | コンテキストファイル（Expertプロンプト注入） |
| `--file` | `-f` | | ファイル添付（繰り返し可） |
| `--model` | `-m` | | モデル指定 |
| `--sandbox` | `-s` | `read-only` | `read-only`, `workspace-write`, `danger-full-access` |
| `--cwd` | | | codex作業ディレクトリ |
| `--raw` | | `false` | stdout/stderrをそのまま出力 |
| `--timeout` | `-t` | `300` | タイムアウト（秒） |

## 設定

任意 `~/.config/aix/config.yaml`:

```yaml
timeout: 600
sandbox: read-only
adapters:
  codex:
    enabled: true
    model: o3
```

優先順位: CLIフラグ > config.yaml > デフォルト値。

## 終了コード

| コード | 意味 |
|--------|------|
| 0 | 成功 |
| 1 | 一般エラー |
| 2 | Codex未インストール |
| 3 | 認証失敗 |
| 4 | Codex実行失敗 |
| 5 | タイムアウト |

## アーキテクチャ

```
cmd/        Cobra CLIコマンド（ask, status）
adapter/    AI CLIバックエンドインターフェース + Codex実装
prompt/     プロンプト組立 + サイズ検証
runner/     サブプロセス抽象化（モックテスト対応）
config/     YAML設定ローダー + マージ優先順位
```

## ライセンス

MIT
