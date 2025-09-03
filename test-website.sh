#!/bin/bash

echo "🌐 KubeAgentic Website Test"
echo "=========================="

# Check if all required files exist
echo "📁 Checking website structure..."

check_file() {
    if [ -f "$1" ]; then
        echo "✅ $1"
    else
        echo "❌ $1 (missing)"
        return 1
    fi
}

# Core Jekyll files
check_file "_config.yml"
check_file "Gemfile"
check_file "index.md"

# Documentation pages
check_file "docs/index.md"
check_file "examples.md"
check_file "local-testing.md"
check_file "api-reference.md"
check_file "website-README.md"

echo ""
echo "📝 Validating content structure..."

# Check front matter in key pages
validate_frontmatter() {
    local file="$1"
    if [ -f "$file" ]; then
        if head -1 "$file" | grep -q "^---$"; then
            echo "✅ $file has front matter"
        else
            echo "⚠️  $file missing front matter"
        fi
    fi
}

validate_frontmatter "index.md"
validate_frontmatter "docs/index.md"
validate_frontmatter "examples.md"
validate_frontmatter "local-testing.md"
validate_frontmatter "api-reference.md"

echo ""
echo "🔗 Checking internal links..."

# Check for common markdown issues
check_markdown() {
    local file="$1"
    if [ -f "$file" ]; then
        # Check for broken relative links
        local broken_links=$(grep -o '\[.*\](.*\.md)' "$file" | grep -v "github.com" | while read link; do
            local target=$(echo "$link" | sed 's/.*(\(.*\))/\1/')
            if [ ! -f "$target" ]; then
                echo "$file: broken link to $target"
            fi
        done)
        
        if [ -z "$broken_links" ]; then
            echo "✅ $file links OK"
        else
            echo "⚠️  $file has broken links:"
            echo "$broken_links"
        fi
    fi
}

check_markdown "index.md"
check_markdown "docs/index.md"
check_markdown "examples.md"

echo ""
echo "⚙️  GitHub Pages Configuration Check..."

# Check _config.yml content
if [ -f "_config.yml" ]; then
    echo "✅ Jekyll configuration found"
    
    # Check key settings
    if grep -q "title:" "_config.yml"; then
        echo "✅ Site title configured"
    else
        echo "⚠️  Site title not configured"
    fi
    
    if grep -q "github-pages" "Gemfile"; then
        echo "✅ GitHub Pages gem configured"
    else
        echo "⚠️  GitHub Pages gem not found in Gemfile"
    fi
else
    echo "❌ Jekyll configuration missing"
fi

echo ""
echo "📊 Content Summary..."

# Count pages and estimate content
total_files=$(find . -name "*.md" -not -path "./.git/*" | wc -l)
total_lines=$(find . -name "*.md" -not -path "./.git/*" -exec wc -l {} + | tail -1 | awk '{print $1}')

echo "📄 Total pages: $total_files"
echo "📝 Total lines: $total_lines"

echo ""
echo "🚀 Next Steps for Local Testing:"
echo "1. Install Ruby 3.0+ (recommended: rbenv or RVM)"
echo "2. Run: gem install bundler jekyll"
echo "3. Run: bundle install"
echo "4. Run: bundle exec jekyll serve"
echo "5. Open: http://localhost:4000"

echo ""
echo "🌐 GitHub Pages Deployment:"
echo "1. Push to GitHub"
echo "2. Go to Settings → Pages"
echo "3. Select 'Deploy from branch: main'"
echo "4. Wait for deployment"

echo ""
echo "✨ Website structure validation complete!"