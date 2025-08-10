#!/bin/bash

# Script to generate release page from HTML template
# Usage: ./generate-release-page.sh <version> <output_dir>

set -e

VERSION="$1"
OUTPUT_DIR="$2"
TEMPLATE_FILE=".github/release-page.html"
OUTPUT_FILE="$OUTPUT_DIR/index.html"

if [ -z "$VERSION" ] || [ -z "$OUTPUT_DIR" ]; then
    echo "Usage: $0 <version> <output_dir>"
    echo "Example: $0 v1.0.0 pages/"
    exit 1
fi

if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "❌ Template file not found: $TEMPLATE_FILE"
    exit 1
fi

echo "🎨 Generating release page for $VERSION..."

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Copy template and replace variables
cp "$TEMPLATE_FILE" "$OUTPUT_FILE"

# Replace placeholders with actual values (cross-platform compatible)
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/{{VERSION}}/$VERSION/g" "$OUTPUT_FILE"
    sed -i '' "s/{{BUILD_DATE}}/$BUILD_DATE/g" "$OUTPUT_FILE"
else
    # Linux
    sed -i "s/{{VERSION}}/$VERSION/g" "$OUTPUT_FILE"
    sed -i "s/{{BUILD_DATE}}/$BUILD_DATE/g" "$OUTPUT_FILE"
fi

# Clean up backup file
rm -f "$OUTPUT_FILE.bak"

echo "✅ Release page generated: $OUTPUT_FILE"
echo "🎯 Template used: $TEMPLATE_FILE"
echo "📦 Version: $VERSION"
echo "📅 Build Date: $BUILD_DATE"
