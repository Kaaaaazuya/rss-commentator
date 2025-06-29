# 🚀 RSS Commentator

<!-- TODO: header 画像を作成したら追加する -->
<!-- ![Header](./assets/header.svg) -->

## 📖 概要

RSS Commentatorは、技術系RSSフィードから最新情報を取得し、DeepSeek APIを活用して自動的に要約・分類するサーバーレスアプリケーションです。エンジニアが効率的に技術情報をキャッチアップできるよう設計されています。

## ✨ 特徴

- **自動要約**: DeepSeek APIによる高精度な要約
- **自動分類**: AIによるタグ付けとスコアリング
- **通知機能**: LINE通知による重要な更新のリアルタイム配信
- **サーバーレスアーキテクチャ**: AWS Lambda + DynamoDBによる効率的な運用

## 🏗️ アーキテクチャ

### システム構成

- **RSS Fetcher**: RSSフィードから記事を取得
- **Summarizer**: DeepSeek APIを使用した記事要約とタグ付け
- **Summary Notifier**: LINE通知による要約配信

### 技術スタック

**バックエンド**
- Go 1.20+ (Lambda Functions)
- AWS Lambda (コンテナイメージ)
- DynamoDB (記事・タグ・関連データ)

**インフラストラクチャ**
- AWS CDK v2 (TypeScript)
- Amazon ECR (コンテナレジストリ)
- AWS EventBridge (イベント駆動)

**外部API**
- DeepSeek API (要約・分類)
- LINE Messaging API (通知)

**開発ツール**
- Task Runner (ビルド・デプロイ自動化)
- Docker & Docker Compose (ローカル開発)
- DynamoDB Local (ローカルデータベース)

## 🚦 開発環境構築

### 前提条件

- Node.js v18+ (CDKに必要)
- Go 1.20+ (Lambda関数開発)
- Docker & Docker Compose (ローカル環境)
- AWS CLI v2 (デプロイに必要)
- [Task Runner](https://taskfile.dev/) (推奨)

### クイックスタート

```bash
# 1. リポジトリをクローン
git clone https://github.com/Kaaaaazuya/rss-commentator.git
cd rss-commentator

# 2. CDK依存関係をインストール
cd infra
npm install
cd ..

# 3. ローカル開発環境を起動
task dev:start

# 4. Lambda関数をローカルでテスト
task lambda:local
task lambda:invoke:rss-fetcher
```

## 🔧 開発コマンド

### ローカル開発

```bash
# 開発環境の起動・停止
task dev:start
task dev:stop
task dev:restart

# Lambda関数のローカル実行
task lambda:local
task lambda:invoke:rss-fetcher
task lambda:invoke:summarizer
task lambda:invoke:notifier
```

### ビルド・デプロイ

```bash
# すべてのLambda関数をビルド
task build-all-lambdas

# Dockerイメージをビルド
task build-all-docker

# AWS にデプロイ
task deploy

# 高速デプロイ（Lambda関数のみ更新）
task deploy:hotswap
```

### インフラ管理

```bash
# CDK操作
task cdk:synth    # CloudFormation生成
task cdk:diff     # 変更差分表示
task cdk:destroy  # スタック削除

# ECR管理
task ecr:create:all  # リポジトリ作成
task ecr:login       # ECRログイン
task push-all        # 全イメージをプッシュ
```

## 📁 プロジェクト構造

```
rss-commentator/
├── go/
│   ├── lambda/                    # Lambda関数
│   │   ├── rss-fetcher/          # RSS取得
│   │   ├── summarizer/           # 要約・分類
│   │   └── summary-notifier/     # LINE通知
│   └── shared/                   # 共有モジュール
│       ├── models/               # データモデル
│       ├── db/                   # DynamoDB操作
│       ├── config/               # 設定管理
│       └── repos/               # リポジトリパターン
├── infra/                        # AWS CDK (TypeScript)
├── docker/                       # Docker Compose設定
├── docs/                         # 設計ドキュメント
└── tools/                        # 開発ツール
```

## 🗄️ データベース設計

### DynamoDB テーブル

- **articles**: 要約記事の保存
  - PK: `url_hash` (UUIDv7), 属性: `url`, `title`, `summary`, `tags`, `created_at`
- **tags**: タグ情報
  - PK: `tag_name`, 属性: `created_at`
- **article_tags**: 記事とタグの関連（スコア付き）
  - PK: `tag_name`, SK: `article_id`, 属性: `score`

## 🔧 ローカル開発のセットアップ

### DynamoDB Local

```bash
# テーブル作成
docker/dynamodb/create-tables.sh

# シードデータ投入
docker/dynamodb/insert-tags.sh

# テーブル削除
docker/dynamodb/delete-tables.sh
```

### 環境変数

Lambda関数で使用される主要な環境変数：

- `AWS_REGION`: デフォルト `ap-northeast-1`
- `AWS_ACCOUNT_ID`: AWS CLI経由で自動取得
- DeepSeek API設定（要約処理用）
- DynamoDBテーブル名（CDK経由で設定）

## 📋 利用可能なタスク

すべてのタスクを確認:
```bash
task --list
```

主要なタスク:
- `task deploy` - AWS にデプロイ
- `task dev:start` - ローカル環境起動
- `task build-all-lambdas` - 全Lambda関数をビルド
- `task ecr:create:all` - ECRリポジトリ作成
- `task push-all` - 全イメージをECRにプッシュ

## 🤝 コントリビューション

1. Issue を作成して変更内容を相談
2. フィーチャーブランチを作成
3. 変更を実装しテストを実行
4. Pull Request を作成

## 📄 ライセンス

MIT License
