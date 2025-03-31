#!/bin/bash

ENDPOINT_URL=${1:-http://localhost:8000}

TABLES=(
  "articles"
  "tags"
  "article_tags"
)

echo "⚠️ DynamoDB テーブル削除スクリプト"
echo "🔗 接続先: $ENDPOINT_URL"
echo ""

for table in "${TABLES[@]}"; do
  EXISTS=$(aws dynamodb list-tables \
    --endpoint-url "$ENDPOINT_URL" \
    --query "TableNames[?@ == '$table']" \
    --output text)

  if [[ "$EXISTS" == "$table" ]]; then
    echo "🗑️ 削除中: $table"
    aws dynamodb delete-table --table-name "$table" --endpoint-url "$ENDPOINT_URL"
  else
    echo "✅ テーブル $table は存在しません、スキップ"
  fi
done

echo ""
echo "🎉 全テーブル削除完了"
