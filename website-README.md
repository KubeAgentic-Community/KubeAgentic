# KubeAgentic Website

This directory contains the GitHub Pages website for KubeAgentic, built with Jekyll and hosted at `https://KubeAgentic-Community.github.io/kubeagentic`.

## 🌐 Website Structure

```
kubeagentic/                    # Repository root
├── _config.yml                 # Jekyll configuration
├── Gemfile                     # Ruby dependencies
├── index.md                    # Homepage
├── docs/                       # Documentation pages
│   └── index.md               # Main documentation
├── examples.md                 # Examples and use cases
├── local-testing.md           # Local testing guide
├── api-reference.md           # API reference
└── website-README.md          # This file
```

## 🚀 Quick Start

### Local Development

1. **Install Ruby and Jekyll**:
   ```bash
   # macOS
   brew install ruby
   gem install bundler jekyll
   
   # Ubuntu/Debian
   sudo apt-get install ruby-full build-essential zlib1g-dev
   gem install bundler jekyll
   ```

2. **Install Dependencies**:
   ```bash
   bundle install
   ```

3. **Start Local Server**:
   ```bash
   bundle exec jekyll serve
   ```

4. **View Website**:
   Open http://localhost:4000 in your browser

### GitHub Pages Deployment

The website is automatically deployed when you push to the main branch.

1. **Enable GitHub Pages**:
   - Go to repository Settings → Pages
   - Source: Deploy from a branch
   - Branch: main / (root)

2. **Configure Custom Domain** (optional):
   ```bash
   echo "kubeagentic.example.com" > CNAME
   ```

3. **Update Repository URL**:
   Edit `_config.yml`:
   ```yaml
   url: "https://sudeshmu.github.io"
   repository: KubeAgentic-Community/kubeagentic
   github_username: sudeshmu
   ```

## 📁 Content Management

### Adding New Documentation

1. **Create Markdown File**:
   ```bash
   # For main docs
   touch docs/new-guide.md
   
   # For standalone pages
   touch new-page.md
   ```

2. **Add Front Matter**:
   ```yaml
   ---
   layout: page
   title: Page Title
   permalink: /path/to/page/
   ---
   
   # Your content here
   ```

3. **Update Navigation**:
   Edit `_config.yml`:
   ```yaml
   header_pages:
     - index.md
     - docs/index.md
     - new-page.md  # Add your page
   ```

### Updating Examples

1. **Edit Examples Page**:
   ```bash
   vim examples.md
   ```

2. **Add New Example**:
   ```markdown
   ## New Example Title
   
   Description of the example...
   
   ```yaml
   # YAML configuration
   ```
   ```

### API Reference Updates

The API reference should be updated when:
- New fields are added to the Agent CRD
- Field validation rules change
- New API versions are released

```bash
vim api-reference.md
```

## 🎨 Customization

### Theme Customization

The website uses the default Minima theme. To customize:

1. **Override Layouts**:
   ```bash
   mkdir -p _layouts
   cp $(bundle show minima)/_layouts/default.html _layouts/
   # Edit _layouts/default.html
   ```

2. **Custom CSS**:
   ```bash
   mkdir -p assets/css
   cat > assets/css/style.scss << 'EOF'
   ---
   ---
   
   @import "minima";
   
   /* Custom styles */
   .highlight {
     background-color: #f0f8ff;
   }
   EOF
   ```

3. **Custom JavaScript**:
   ```bash
   mkdir -p assets/js
   # Add your JavaScript files
   ```

### Color Scheme

Edit `_config.yml` to add custom variables:
```yaml
minima:
  skin: dark  # or auto, classic
```

## 🔧 Configuration

### Jekyll Configuration

Key settings in `_config.yml`:

```yaml
# Site settings
title: KubeAgentic
description: Deploy and manage AI agents on Kubernetes
baseurl: ""
url: "https://your-username.github.io"

# Build settings
markdown: kramdown
highlighter: rouge
theme: minima

# Plugins
plugins:
  - jekyll-feed
  - jekyll-sitemap
  - jekyll-seo-tag

# Collections (for organizing content)
collections:
  docs:
    output: true
    permalink: /:collection/:name/
```

### GitHub Pages Settings

For optimal GitHub Pages performance:

```yaml
# In _config.yml
exclude:
  - vendor/
  - Gemfile
  - Gemfile.lock
  - node_modules/
  - "*.go"
  - go.mod
  - go.sum
  - Dockerfile*
  - Makefile
  - .git/
  - .gitignore
  - bin/
  - scripts/
```

## 📝 Content Guidelines

### Writing Style

- **Concise**: Keep explanations clear and brief
- **Practical**: Include working examples and code snippets
- **Structured**: Use consistent heading hierarchy
- **Accessible**: Explain technical concepts clearly

### Code Examples

Always include complete, working examples:

```yaml
# ✅ Good: Complete example
apiVersion: ai.example.com/v1
kind: Agent
metadata:
  name: example-agent
spec:
  provider: openai
  model: gpt-4
  systemPrompt: "You are a helpful assistant."
  apiSecretRef:
    name: openai-secret
    key: api-key
```

```yaml
# ❌ Bad: Incomplete example
spec:
  provider: openai
  # Missing required fields
```

### Screenshots

When adding screenshots:

1. **Optimize Images**:
   ```bash
   mkdir -p assets/images
   # Add optimized images
   ```

2. **Reference in Markdown**:
   ```markdown
   ![Description](assets/images/screenshot.png)
   ```

## 🚦 Testing

### Local Testing

1. **Content Validation**:
   ```bash
   # Check for broken links
   bundle exec jekyll build
   bundle exec htmlproofer ./_site --check-html --check-opengraph
   ```

2. **Performance Testing**:
   ```bash
   # Test build time
   time bundle exec jekyll build
   
   # Test site speed
   lighthouse http://localhost:4000
   ```

### Automated Testing

GitHub Actions workflow (`.github/workflows/pages.yml`):

```yaml
name: Deploy GitHub Pages
on:
  push:
    branches: [main]
jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    - uses: actions/setup-ruby@v1
      with:
        ruby-version: 3.0
    - run: |
        gem install bundler
        bundle install
        bundle exec jekyll build
    - uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./_site
```

## 🔍 SEO Optimization

### Meta Tags

The website automatically includes SEO tags via `jekyll-seo-tag`. Customize in front matter:

```yaml
---
title: Page Title
description: Page description for search engines
image: /assets/images/social-preview.png
---
```

### Sitemap

Automatically generated via `jekyll-sitemap` plugin at `/sitemap.xml`.

### Analytics

Add Google Analytics (optional):

```yaml
# In _config.yml
google_analytics: UA-XXXXXXXX-X
```

## 📊 Maintenance

### Regular Updates

1. **Content Review**: Monthly review of documentation for accuracy
2. **Dependency Updates**: Keep Jekyll and gems updated
3. **Link Checking**: Verify external links are still valid
4. **Performance Monitoring**: Check site speed and responsiveness

### Updating Dependencies

```bash
# Update Gemfile.lock
bundle update

# Update Jekyll
bundle update jekyll

# Update all gems
bundle update --all
```

### Monitoring

- **GitHub Pages Status**: Check GitHub repository Insights → Traffic
- **Analytics**: Monitor page views and user behavior
- **Performance**: Regular Lighthouse audits

## 🐛 Troubleshooting

### Common Issues

**Build failures**:
```bash
# Clear Jekyll cache
bundle exec jekyll clean

# Rebuild with verbose output
bundle exec jekyll build --verbose
```

**Plugin errors**:
```bash
# Check plugin compatibility
bundle exec jekyll doctor

# Update plugins
bundle update jekyll-feed jekyll-sitemap jekyll-seo-tag
```

**Styling issues**:
```bash
# Force regenerate CSS
rm -rf _site
bundle exec jekyll build
```

### GitHub Pages Limitations

- **Repository size**: Max 1GB
- **File size**: Max 100MB per file
- **Build time**: Max 10 minutes
- **Monthly bandwidth**: 100GB

### Alternative Hosting

If GitHub Pages limitations are reached:

1. **Netlify**:
   ```bash
   # netlify.toml
   [build]
     command = "bundle exec jekyll build"
     publish = "_site"
   ```

2. **Vercel**:
   ```bash
   # vercel.json
   {
     "buildCommand": "bundle exec jekyll build",
     "outputDirectory": "_site"
   }
   ```

## 🤝 Contributing

### Content Contributions

1. **Fork** the repository
2. **Create** a feature branch: `git checkout -b docs/new-feature`
3. **Edit** documentation files
4. **Test** locally: `bundle exec jekyll serve`
5. **Submit** a pull request

### Style Guide

- Use `###` for section headings within pages
- Include code examples for all features
- Add front matter to all pages
- Use consistent terminology
- Optimize images before adding

## 📚 Resources

### Jekyll Documentation
- [Jekyll Docs](https://jekyllrb.com/docs/)
- [Liquid Templating](https://shopify.github.io/liquid/)
- [Kramdown Syntax](https://kramdown.gettalong.org/syntax.html)

### GitHub Pages
- [GitHub Pages Docs](https://docs.github.com/en/pages)
- [Supported Themes](https://pages.github.com/themes/)
- [Custom Domains](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site)

### Minima Theme
- [Minima Repository](https://github.com/jekyll/minima)
- [Theme Customization](https://github.com/jekyll/minima#customization)

---

For technical questions about the website, please open an issue in the [KubeAgentic repository](https://github.com/KubeAgentic-Community/kubeagentic/issues).