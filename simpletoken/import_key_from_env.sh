#!/bin/bash

set -e

# .env ファイルを読み込む
if [ ! -f .env ]; then
  echo "❌ .env ファイルが見つかりません"
  exit 1
fi

# .env から PRIVATE_KEY を読み取り（bash-safe）
PRIVATE_KEY=$(grep "^PRIVATE_KEY=" .env | cut -d '=' -f2 | tr -d '"')

if [ -z "$PRIVATE_KEY" ]; then
  echo "❌ PRIVATE_KEY が定義されていません"
  exit 1
fi

# 先頭の "0x" を除去
CLEANED_KEY=$(echo "$PRIVATE_KEY" | sed 's/^0x//')

if [ ${#CLEANED_KEY} -ne 64 ]; then
  echo "❌ 秘密鍵の長さが不正です（64 hex文字）: 現在 ${#CLEANED_KEY} 文字"
  exit 1
fi

# 一時ファイルに秘密鍵を改行なしで保存
TEMP_FILE=$(mktemp)
# TEMP_FILE="temp_private_key.txt"
echo -n "$CLEANED_KEY" > "$TEMP_FILE"

# Geth アカウントをインポート
echo "🔐 Geth に秘密鍵をインポート中..."
geth  --datadir ../devchain account import "$TEMP_FILE"

# クリーンアップ
rm -f "$TEMP_FILE"
