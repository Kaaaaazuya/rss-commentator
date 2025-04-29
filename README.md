# 🚀 RSS Commentator

<!-- TODO: header 画像を作成したら追加する -->
<!-- ![Header](./assets/header.svg) -->

## 📖 概要

RSS Commentatorは、技術系RSSフィードから最新情報を取得し、DeepSeek APIを活用して自動的に要約・分類するツールです。エンジニアが効率的に技術情報をキャッチアップできるよう設計されています。

## ✨ 特徴

- **自動要約**: DeepSeek APIによる高精度な要約
- **多言語対応**: 英語・日本語の技術情報に対応
- **通知機能**: 重要な更新をリアルタイムで通知
- **サーバーレスアーキテクチャ**: AWS CDK v2による効率的なインフラ管理

## 🛠️ 技術スタック

### フロントエンド

- TypeScript
- React
- Next.js

### バックエンド

- Golang

### インフラ

- AWS CDK v2
- DynamoDB
- AWS Lambda

### 生成AI

- DeepSeek API

## 🚦 開発環境構築

### 前提条件

- Node.js v18+
- Go 1.20+
- Docker
- AWS CLI

### セットアップ手順

```bash
# 依存関係のインストール
npm install

# ローカル環境の起動
docker-compose -f local/docker-compose.yml up -d

# フロントエンドの起動
cd frontend
npm run dev

# バックエンドの起動
cd backend
go run main.go
```

## 🤝 コントリビューション

コントリビューションガイドラインは[CONTRIBUTING.md](./CONTRIBUTING.md)をご覧ください。

## 📄 ライセンス

MIT License
