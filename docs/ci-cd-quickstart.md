# CI/CD Pipeline Quick Start Guide

> **Quick reference for using the Pomodux CI/CD pipeline**

## 🚀 Getting Started

### Prerequisites

- Go 1.24.4 or later
- Git
- GitHub repository with Actions enabled

### Setup (One-time)

1. **Enable GitHub Actions**:
   - Go to your repository Settings → Actions → General
   - Enable "Allow all actions and reusable workflows"

2. **Install development tools**:
   ```bash
   make tools
   ```

3. **Setup development environment**:
   ```bash
   make setup
   ```

## 🔄 Daily Development Workflow

### 1. Feature Development

```bash
# Create feature branch
git checkout -b feature/new-feature

# Develop with TDD
# ... write tests and implement features ...

# Run CI checks locally
make ci-check

# Commit and push
git add .
git commit -m "Add new feature"
git push origin feature/new-feature
```

### 2. Pull Request Process

1. Create pull request on GitHub
2. CI workflow runs automatically
3. Address any issues found
4. Merge when all checks pass

## 🏷️ Creating a Release

### Option 1: Using the Release Script (Recommended)

```bash
# Ensure you're on main branch
git checkout main

# Create release (replace 1.2.3 with your version)
./scripts/release.sh 1.2.3
```

### Option 2: Using Make

```bash
# Create release (replace 1.2.3 with your version)
make create-release VERSION=1.2.3
```

### Option 3: Manual Process

```bash
# 1. Update version in internal/cli/version.go
# 2. Commit changes
git add .
git commit -m "Release v1.2.3"

# 3. Create and push tag
git tag v1.2.3
git push origin v1.2.3
```

## 📋 What Happens During Release

1. **Pre-release checks**:
   - Validates version format
   - Ensures you're on main branch
   - Checks working directory is clean
   - Runs tests and linting

2. **GitHub Actions workflow**:
   - Builds for 6 platforms (Linux/macOS/Windows, amd64/arm64)
   - Creates compressed archives
   - Generates checksums
   - Publishes GitHub release
   - Runs UAT tests

3. **Release artifacts**:
   - Platform-specific binaries
   - Compressed archives (.tar.gz, .zip)
   - SHA256 checksums
   - Release notes

## 🔍 Monitoring

### Check Workflow Status

- **GitHub Actions**: https://github.com/[username]/pomodux/actions
- **Releases**: https://github.com/[username]/pomodux/releases

### Local Verification

```bash
# Run all CI checks locally
make ci-check

# Test multi-platform builds
make build-all

# Run UAT tests
./tests/uat/automated/run-tests.sh
```

## 🛠️ Troubleshooting

### Common Issues

**CI fails locally but passes on GitHub**:
```bash
# Clean and rebuild
make clean
make install
make ci-check
```

**Release workflow fails**:
```bash
# Check tag format
git tag -l | grep v1.2.3

# Delete failed tag and retry
git tag -d v1.2.3
git push origin --delete v1.2.3
./scripts/release.sh 1.2.3
```

**Build fails on specific platform**:
```bash
# Test specific platform build
GOOS=linux GOARCH=amd64 go build -o bin/pomodux cmd/pomodux/main.go
```

### Getting Help

1. Check GitHub Actions logs for specific errors
2. Run commands locally to reproduce issues
3. Review the full [CI/CD Pipeline Documentation](ci-cd-pipeline.md)

## 📚 Next Steps

- Read the [full CI/CD documentation](ci-cd-pipeline.md)
- Review [release management process](release-management.md)
- Explore [development setup guide](development-setup.md)

---

**Need help?** Check the full documentation or create an issue on GitHub. 