#!/bin/bash

set -e

# 引数からデプロイ環境を取得（デフォルトはdev）
ENVIRONMENT=${1:-dev}
AWS_REGION=${2:-ap-northeast-1}
# 特定のLambda関数のみをビルド・デプロイする場合のオプション
# 例: ./deploy-cdk.sh dev ap-northeast-1 --bootstrap --lambda=rss-fetcher
LAMBDA_FILTER=${3:-all}

echo "=== RSS Commentatorをデプロイします（環境: $ENVIRONMENT, リージョン: $AWS_REGION） ==="

# 作業ディレクトリをルートディレクトリに設定
ROOT_DIR="$(dirname "$0")"
cd "$ROOT_DIR"

# 環境変数の設定
export ENVIRONMENT=$ENVIRONMENT
export AWS_REGION=$AWS_REGION

# CDKのデプロイ処理に移行
cd "$ROOT_DIR/infra"

# 依存関係のインストール
echo "📦 依存関係をインストールしています..."
npm ci

# ビルド
echo "🔨 TypeScriptをビルドしています..."
npm run build

# CDKブートストラップ（初回のみ必要）
if [[ "$3" == "--bootstrap" || "$4" == "--bootstrap" ]]; then
  echo "🥾 CDKブートストラップを実行しています..."
  npm run bootstrap
fi

# リソースを再利用するオプションを追加
if [[ "$3" == "--reuse-resources" || "$4" == "--reuse-resources" || "$5" == "--reuse-resources" ]]; then
  echo "♻️ 既存のリソースを再利用します"
  export REUSE_EXISTING_RESOURCE=true
fi

# スタック内容の差分を表示
echo "📊 デプロイする変更内容を確認しています..."
npm run diff

# デプロイ確認
echo "🚀 デプロイを実行しますか？ [y/N]"
read -r confirmation

if [[ "$confirmation" == "y" || "$confirmation" == "Y" ]]; then
  echo "🚀 デプロイを開始します..."
  if [[ "$3" == "--hotswap" || "$4" == "--hotswap" || "$5" == "--hotswap" ]]; then
    npm run deploy:hotswap
  else
    npm run deploy:with-approval
  fi
  echo "✅ デプロイが完了しました！"
else
  echo "❌ デプロイをキャンセルしました。"
  exit 1
fi