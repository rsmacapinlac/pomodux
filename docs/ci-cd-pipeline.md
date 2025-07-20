# CI/CD Pipeline Documentation

> **Note:** This document describes the comprehensive CI/CD pipeline for Pomodux, including continuous integration, automated releases, and multi-platform builds.

## 1.0 Pipeline Overview

### 1.1 Pipeline Architecture

The Pomodux CI/CD pipeline consists of two main workflows:

1. **CI Workflow** (`.github/workflows/ci.yml`)
   - Runs on every push and pull request
   - Ensures code quality and build verification
   - Provides early feedback on issues

2. **Release Workflow** (`.github/workflows/release.yml`)
   - Triggers on git tags (e.g., `v1.0.0`, `v2.1.3`)
   - Creates multi-platform binaries
   - Publishes GitHub releases
   - Runs UAT tests on release builds

### 1.2 Pipeline Integration with Release Management

The CI/CD pipeline integrates seamlessly with the [Release Management Process](release-management.md):

- **Gate 2 (Development Completion)**: CI workflow ensures all tests pass
- **Gate 3 (User Acceptance)**: UAT tests run automatically in CI
- **Gate 4 (Release Approval)**: Release workflow creates final artifacts

## 2.0 CI Workflow Details

### 2.1 Triggers

```yaml
on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
```

The CI workflow runs on:
- Every push to `main` or `develop` branches
- Every pull request targeting `main` or `develop` branches

### 2.2 Jobs

#### 2.2.1 Test Job

**Purpose**: Run unit tests with multiple Go versions

**Features**:
- Matrix testing with Go 1.24.4 and 1.25.0
- Race condition detection
- Coverage reporting
- Integration with Codecov

**Outputs**:
- Test results and coverage metrics
- Coverage report uploaded to Codecov

#### 2.2.2 Lint Job

**Purpose**: Ensure code quality and style consistency

**Features**:
- Uses golangci-lint for comprehensive linting
- Checks for common Go issues
- Enforces coding standards

#### 2.2.3 Build Job

**Purpose**: Verify the application builds successfully

**Features**:
- Builds the application for Linux amd64
- Verifies the binary is executable
- Uploads build artifacts for later use

#### 2.2.4 Security Job

**Purpose**: Identify security vulnerabilities

**Features**:
- Uses gosec for security scanning
- Checks for common security issues
- Reports vulnerabilities in code

#### 2.2.5 UAT Tests Job

**Purpose**: Run User Acceptance Tests

**Features**:
- Installs bats-core testing framework
- Runs automated UAT tests
- Uploads test results as artifacts

## 3.0 Release Workflow Details

### 3.1 Git Tag Strategy

#### 3.1.1 Tag Format

The release workflow uses semantic versioning tags:

```
v{major}.{minor}.{patch}
```

**Examples**:
- `v1.0.0` - Major release
- `v1.1.0` - Minor release with new features
- `v1.1.1` - Patch release with bug fixes

#### 3.1.2 Creating a Release

To create a new release:

1. **Prepare the release**:
   ```bash
   # Ensure all changes are committed
   git add .
   git commit -m "Prepare for release v1.2.0"
   
   # Create and push the tag
   git tag v1.2.0
   git push origin v1.2.0
   ```

2. **Monitor the workflow**:
   - The release workflow will automatically trigger
   - Check GitHub Actions for progress
   - Review the generated release

3. **Verify the release**:
   - Check that all platform binaries are created
   - Verify UAT tests pass
   - Review release notes

#### 3.1.3 Tag Best Practices

- **Use semantic versioning**: Follow `MAJOR.MINOR.PATCH` format
- **Tag from main branch**: Always create release tags from the main branch
- **Include release notes**: Update `RELEASE_NOTES.md` before tagging
- **Test before tagging**: Ensure CI passes before creating a release tag

### 3.2 Multi-Platform Builds

#### 3.2.1 Supported Platforms

The release workflow builds for the following platforms:

| OS | Architecture | Binary Name | Archive Format |
|----|--------------|-------------|----------------|
| Linux | amd64 | `pomodux-linux-amd64` | `.tar.gz` |
| Linux | arm64 | `pomodux-linux-arm64` | `.tar.gz` |
| macOS | amd64 | `pomodux-darwin-amd64` | `.tar.gz` |
| macOS | arm64 | `pomodux-darwin-arm64` | `.tar.gz` |
| Windows | amd64 | `pomodux-windows-amd64.exe` | `.zip` |
| Windows | arm64 | `pomodux-windows-arm64.exe` | `.zip` |

#### 3.2.2 Build Configuration

**Build flags used**:
```bash
-ldflags="-s -w -X main.Version=$VERSION"
```

- `-s`: Strip debug information (smaller binaries)
- `-w`: Strip DWARF symbol table (smaller binaries)
- `-X main.Version=$VERSION`: Inject version information

#### 3.2.3 Version Injection

The build process injects version information into the binary:

```go
// Version will be set during build time via ldflags
var Version = "dev"
```

During release builds, this is set to the actual version (e.g., "1.2.0").

### 3.3 Release Artifacts

#### 3.3.1 Generated Files

For each release, the following artifacts are created:

1. **Platform-specific binaries**:
   - 6 binary files (one per platform/architecture)

2. **Compressed archives**:
   - 4 `.tar.gz` files for Linux and macOS
   - 2 `.zip` files for Windows

3. **Checksums file**:
   - `pomodux-{version}-checksums.txt` with SHA256 hashes

4. **Release notes**:
   - Extracted from `RELEASE_NOTES.md` or generated from git log

#### 3.3.2 GitHub Release

The workflow automatically creates a GitHub release with:
- Release name and description
- All platform binaries as downloadable assets
- Checksums for verification
- Release notes

## 4.0 Workflow Integration

### 4.1 Integration with Approval Gates

#### 4.1.1 Gate 2: Development Completion

**CI Requirements**:
- All tests must pass
- Code coverage meets requirements
- Linting passes without errors
- Security scan shows no critical issues
- Build verification successful

**Workflow Integration**:
```yaml
# CI workflow ensures these requirements are met
test:
  # Runs unit tests with coverage
lint:
  # Ensures code quality
security:
  # Identifies security issues
build:
  # Verifies successful build
```

#### 4.1.2 Gate 3: User Acceptance

**UAT Requirements**:
- All UAT tests pass
- User acceptance criteria met
- No critical issues identified

**Workflow Integration**:
```yaml
# Both CI and Release workflows run UAT tests
uat-tests:
  # Runs automated UAT tests
  # Uploads test results as artifacts
```

#### 4.1.3 Gate 4: Release Approval

**Release Requirements**:
- All release artifacts created successfully
- Multi-platform builds verified
- Release notes complete
- Installation procedures tested

**Workflow Integration**:
```yaml
# Release workflow handles all release requirements
release:
  # Creates multi-platform binaries
  # Generates release artifacts
  # Publishes GitHub release
  # Runs final UAT verification
```

### 4.2 Quality Gates

#### 4.2.1 Required Checks

Before a release can be approved, the following must pass:

1. **CI Workflow**:
   - ✅ All tests pass
   - ✅ Code coverage ≥ 80%
   - ✅ Linting passes
   - ✅ Security scan clean
   - ✅ Build verification

2. **Release Workflow**:
   - ✅ Multi-platform builds successful
   - ✅ UAT tests pass on release build
   - ✅ Release artifacts created
   - ✅ GitHub release published

#### 4.2.2 Failure Handling

**CI Failures**:
- Fix issues in development branch
- Re-run CI workflow
- Ensure all checks pass before proceeding

**Release Failures**:
- Investigate build or test failures
- Fix issues in main branch
- Delete failed tag and re-tag
- Re-run release workflow

## 5.0 Local Development

### 5.1 Testing the Pipeline Locally

#### 5.1.1 Running CI Checks Locally

```bash
# Install development tools
make tools

# Run all CI checks
make test
make lint
make build
make security

# Run UAT tests
./tests/uat/automated/run-tests.sh
```

#### 5.1.2 Testing Multi-Platform Builds

```bash
# Build for all platforms locally
make build-all

# Verify binaries
file bin/pomodux-*
./bin/pomodux-linux-amd64 --version
```

#### 5.1.3 Testing Release Process

```bash
# Create a test tag
git tag v0.0.1-test
git push origin v0.0.1-test

# Monitor the workflow
# Delete test tag when done
git tag -d v0.0.1-test
git push origin --delete v0.0.1-test
```

### 5.2 Development Workflow

#### 5.2.1 Feature Development

1. **Create feature branch**:
   ```bash
   git checkout -b feature/new-feature
   ```

2. **Develop with TDD**:
   - Write tests first
   - Implement features
   - Ensure all tests pass

3. **Push and create PR**:
   ```bash
   git push origin feature/new-feature
   # Create pull request on GitHub
   ```

4. **CI verification**:
   - CI workflow runs automatically
   - Address any issues found
   - Ensure all checks pass

#### 5.2.2 Release Preparation

1. **Update release notes**:
   ```bash
   # Edit RELEASE_NOTES.md
   # Add new version section
   ```

2. **Final testing**:
   ```bash
   make test
   make build
   ./tests/uat/automated/run-tests.sh
   ```

3. **Create release tag**:
   ```bash
   git checkout main
   git pull origin main
   git tag v1.2.0
   git push origin v1.2.0
   ```

## 6.0 Monitoring and Maintenance

### 6.1 Pipeline Monitoring

#### 6.1.1 GitHub Actions Dashboard

Monitor pipeline health through:
- GitHub Actions tab in repository
- Workflow run history
- Job status and logs

#### 6.1.2 Key Metrics

Track the following metrics:
- **Build success rate**: Percentage of successful builds
- **Test coverage**: Maintain ≥ 80% overall coverage
- **UAT pass rate**: Percentage of successful UAT runs
- **Release frequency**: Time between releases

### 6.2 Maintenance Tasks

#### 6.2.1 Regular Updates

- **Update Go version**: Keep Go version current
- **Update dependencies**: Regular `go mod tidy` and updates
- **Update GitHub Actions**: Keep actions current
- **Review security alerts**: Address GitHub security alerts

#### 6.2.2 Artifact Cleanup

- **CI artifacts**: Automatically cleaned up after 30 days
- **Release artifacts**: Automatically cleaned up after 90 days
- **Test reports**: Automatically cleaned up after 30 days

### 6.3 Troubleshooting

#### 6.3.1 Common Issues

**Build Failures**:
- Check Go version compatibility
- Verify all dependencies are available
- Review build logs for specific errors

**Test Failures**:
- Run tests locally to reproduce
- Check for environment-specific issues
- Review test coverage requirements

**Release Failures**:
- Verify tag format is correct
- Check that all CI checks pass
- Review release workflow logs

#### 6.3.2 Debugging Steps

1. **Reproduce locally**:
   ```bash
   # Run the same commands locally
   go test -v ./...
   go build -v cmd/pomodux/main.go
   ```

2. **Check environment**:
   ```bash
   # Verify Go version
   go version
   
   # Check dependencies
   go mod verify
   ```

3. **Review logs**:
   - Check GitHub Actions logs
   - Look for specific error messages
   - Verify environment variables

## 7.0 Security Considerations

### 7.1 Security Scanning

The pipeline includes multiple security checks:

- **gosec**: Static analysis for security issues
- **Dependency scanning**: GitHub's automated dependency scanning
- **Code review**: Manual security review in pull requests

### 7.2 Secrets Management

- **No secrets in workflows**: All secrets are managed through GitHub Secrets
- **Minimal permissions**: Workflows use minimal required permissions
- **Audit trail**: All workflow runs are logged and auditable

### 7.3 Binary Security

- **Stripped binaries**: Release binaries are stripped of debug information
- **Checksums**: All binaries include SHA256 checksums for verification
- **Reproducible builds**: Build process is deterministic and reproducible

## 8.0 Future Enhancements

### 8.1 Planned Improvements

- **Docker images**: Add Docker image builds for containerized deployment
- **Homebrew formula**: Automate Homebrew formula updates
- **Chocolatey package**: Add Windows package manager support
- **Signing**: Add GPG signing for release binaries
- **Performance testing**: Add performance benchmarks to CI

### 8.2 Integration Opportunities

- **Slack notifications**: Add notifications for build status
- **JIRA integration**: Link releases to JIRA tickets
- **Metrics dashboard**: Create dashboard for pipeline metrics
- **Automated changelog**: Generate changelog from commit messages

---

**Last Updated:** 2025-01-27 