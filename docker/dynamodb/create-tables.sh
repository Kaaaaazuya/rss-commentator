#!/bin/bash

# デフォルトはローカル環境用
ENDPOINT_URL=${1:-http://localhost:8000}
TABLE_DIR="./tables"

# 対象のテーブル定義ファイル（拡張可能）
TABLES=(
  "articles"
  "tags"
  "article_tags"
)

echo "📦 DynamoDB テーブル作成スクリプト開始"
echo "🔗 接続先: $ENDPOINT_URL"
echo ""

for table in "${TABLES[@]}"; do
  FILE="$TABLE_DIR/${table}-table.json"

  if [ ! -f "$FILE" ]; then
    echo "❌ テーブル定義が見つかりません: $FILE"
    continue
  fi

  EXISTS=$(aws dynamodb list-tables --endpoint-url "$ENDPOINT_URL" | grep "\"$table\"")

  if [ -z "$EXISTS" ]; then
    echo "✅ $table テーブルを作成します"
    aws dynamodb create-table \
      --cli-input-json "file://$FILE" \
      --endpoint-url "$ENDPOINT_URL"
  else
    echo "⚠️  $table テーブルは既に存在します、スキップします"
  fi

  echo ""
done

echo "🎉 テーブル作成スクリプト完了"
