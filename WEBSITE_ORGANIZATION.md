# Website Organization

## Current Status

The website files are currently in the root directory for immediate functionality, but we have a `website/` folder prepared for future organization.

## Structure

### Current (Root Directory)
```
├── _config.yml          # Jekyll configuration
├── _layouts/            # Custom page layouts
├── _includes/           # Reusable HTML components
├── assets/              # Static assets
├── index.md             # Homepage
├── direct-framework.md  # Direct Framework guide
├── langgraph-framework.md # LangGraph Framework guide
├── api-reference.md     # API documentation
├── examples.md          # Usage examples
├── local-testing.md     # Local testing guide
├── CNAME                # Custom domain configuration
└── Gemfile              # Ruby dependencies
```

### Future Organization (website/ folder)
```
website/
├── _config.yml          # Jekyll configuration
├── _layouts/            # Custom page layouts
├── _includes/           # Reusable HTML components
├── assets/              # Static assets
├── index.md             # Homepage
├── direct-framework.md  # Direct Framework guide
├── langgraph-framework.md # LangGraph Framework guide
├── api-reference.md     # API documentation
├── examples.md          # Usage examples
├── local-testing.md     # Local testing guide
├── CNAME                # Custom domain configuration
├── Gemfile              # Ruby dependencies
└── README.md            # Website documentation
```

## Migration to website/ folder

To move the website to the `website/` folder:

1. **Enable GitHub Actions**: Configure GitHub Pages to use GitHub Actions instead of the default Jekyll build
2. **Update Settings**: In repository settings, change Pages source to "GitHub Actions"
3. **Move Files**: Move all website files to the `website/` folder
4. **Test Deployment**: Verify the GitHub Actions workflow deploys correctly

## Benefits of Organization

- **Cleaner Structure**: Separates website from project code
- **Better Maintenance**: Easier to manage website-specific files
- **Contributor Friendly**: Clear separation of concerns
- **Deployment Control**: More control over website deployment process

## GitHub Actions Workflow

The `.github/workflows/deploy-website.yml` file is ready for when we migrate to the `website/` folder structure.
