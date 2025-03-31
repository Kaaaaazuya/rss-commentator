#!/bin/bash

# 使用方法:
#   bash insert-tags.sh          # ローカル用（http://localhost:8000）
#   bash insert-tags.sh prod     # 本番用（--endpoint-url なし）

ENV=$1
SEED_FILE="seed/tags-seed.sh"
TODAY=$(date -I)  # ISO形式: 例）2025-03-30

# エンドポイント切り替え
if [[ "$ENV" == "prod" ]]; then
  ENDPOINT=""
  echo "🚀 本番環境に投入します"
else
  ENDPOINT="--endpoint-url http://localhost:8000"
  echo "🧪 ローカル環境に投入します"
fi

# 各行に --endpoint-url を追加し、created_at を今日の日付に置換
while IFS= read -r line; do
  if [[ "$line" == aws\ dynamodb\ put-item* ]]; then
    # created_at の日付部分（YYYY-MM-DD）を $TODAY に置換
    line_updated=$(echo "$line" | sed -E "s/\"created_at\":\s*\{\"S\":\s*\"[0-9\-]+\"\}/\"created_at\": {\"S\": \"$TODAY\"}/")

    echo "▶️ 実行: $line_updated $ENDPOINT"
    eval "$line_updated $ENDPOINT"
  fi
done < "$SEED_FILE"
