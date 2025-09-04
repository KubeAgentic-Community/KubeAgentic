# KubeAgentic Website

This folder contains all the website-related files for the KubeAgentic project.

## Structure

```
website/
├── _config.yml          # Jekyll configuration
├── _layouts/            # Custom page layouts
│   └── page.html        # Main page layout
├── _includes/           # Reusable HTML components
│   └── head-custom.html # Custom head elements
├── assets/              # Static assets
│   ├── logo.png         # Main logo
│   └── favicon.ico      # Browser favicon
├── index.md             # Homepage
├── direct-framework.md  # Direct Framework guide
├── langgraph-framework.md # LangGraph Framework guide
├── api-reference.md     # API documentation
├── examples.md          # Usage examples
├── local-testing.md     # Local testing guide
├── CNAME                # Custom domain configuration
├── Gemfile              # Ruby dependencies
└── README.md            # This file
```

## Development

The website is built using Jekyll and deployed via GitHub Actions.

### Local Development

1. Install Ruby and Bundler
2. Navigate to the website folder:
   ```bash
   cd website
   ```
3. Install dependencies:
   ```bash
   bundle install
   ```
4. Start the development server:
   ```bash
   bundle exec jekyll serve
   ```
5. Open http://localhost:4000 in your browser

### Deployment

The website is automatically deployed to GitHub Pages when changes are pushed to the main branch. The deployment is handled by the GitHub Actions workflow in `.github/workflows/deploy-website.yml`.

## Customization

- **Styling**: All custom CSS is in the `_layouts/page.html` file
- **Content**: Edit the `.md` files to update content
- **Assets**: Add images and other assets to the `assets/` folder
- **Layout**: Modify `_layouts/page.html` to change the page structure

## Features

- Responsive design that works on all devices
- Custom gradient hero sections
- Dark theme code blocks
- Breadcrumb navigation
- SEO optimization
- Custom favicon and meta tags
- Professional typography and spacing
