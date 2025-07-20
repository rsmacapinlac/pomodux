#!/bin/bash

# Pomodux Release Script
# This script helps automate the release process

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to validate version format
validate_version() {
    local version=$1
    if [[ ! $version =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        print_error "Invalid version format. Use semantic versioning (e.g., 1.2.3)"
        exit 1
    fi
}

# Function to check if we're on main branch
check_branch() {
    local current_branch=$(git branch --show-current)
    if [[ $current_branch != "main" ]]; then
        print_error "You must be on the main branch to create a release"
        print_error "Current branch: $current_branch"
        exit 1
    fi
}

# Function to check if working directory is clean
check_clean_working_dir() {
    if [[ -n $(git status --porcelain) ]]; then
        print_error "Working directory is not clean. Please commit or stash changes."
        git status --short
        exit 1
    fi
}

# Function to check if tag already exists
check_tag_exists() {
    local version=$1
    if git tag -l | grep -q "^v$version$"; then
        print_error "Tag v$version already exists"
        exit 1
    fi
}

# Function to run pre-release checks
run_pre_release_checks() {
    print_status "Running pre-release checks..."
    
    # Check if required commands exist
    if ! command_exists go; then
        print_error "Go is not installed"
        exit 1
    fi
    
    if ! command_exists git; then
        print_error "Git is not installed"
        exit 1
    fi
    
    # Check branch
    check_branch
    
    # Check working directory
    check_clean_working_dir
    
    # Pull latest changes
    print_status "Pulling latest changes from remote..."
    git pull origin main
    
    # Run tests
    print_status "Running tests..."
    if ! make test > /dev/null 2>&1; then
        print_error "Tests failed"
        exit 1
    fi
    
    # Run linting
    print_status "Running linting..."
    if ! make lint > /dev/null 2>&1; then
        print_error "Linting failed"
        exit 1
    fi
    
    # Build verification
    print_status "Building application..."
    if ! make build > /dev/null 2>&1; then
        print_error "Build failed"
        exit 1
    fi
    
    print_success "Pre-release checks passed"
}

# Function to update version in files
update_version() {
    local version=$1
    print_status "Updating version to $version..."
    
    # Update version in internal/cli/version.go
    if [[ -f "internal/cli/version.go" ]]; then
        sed -i.bak "s/Version   = \"[^\"]*\"/Version   = \"$version\"/" internal/cli/version.go
        rm internal/cli/version.go.bak
    fi
    
    print_success "Version updated in source files"
}

# Function to create release commit
create_release_commit() {
    local version=$1
    print_status "Creating release commit..."
    
    git add .
    git commit -m "Release v$version"
    
    print_success "Release commit created"
}

# Function to create and push tag
create_tag() {
    local version=$1
    print_status "Creating tag v$version..."
    
    git tag "v$version"
    git push origin "v$version"
    
    print_success "Tag v$version created and pushed"
}

# Function to show next steps
show_next_steps() {
    local version=$1
    print_success "Release v$version initiated!"
    echo
    print_status "Next steps:"
    echo "1. Monitor the GitHub Actions workflow:"
    echo "   https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/actions"
    echo
    echo "2. Check the release page once the workflow completes:"
    echo "   https://github.com/$(git config --get remote.origin.url | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/releases"
    echo
    echo "3. Verify that all platform binaries are available"
    echo "4. Test the release on different platforms"
    echo "5. Update documentation if needed"
}

# Function to show usage
show_usage() {
    echo "Usage: $0 <version>"
    echo
    echo "Creates a new release with the specified version"
    echo
    echo "Arguments:"
    echo "  version    Semantic version (e.g., 1.2.3)"
    echo
    echo "Examples:"
    echo "  $0 1.0.0    # Create major release"
    echo "  $0 1.1.0    # Create minor release"
    echo "  $0 1.1.1    # Create patch release"
    echo
    echo "The script will:"
    echo "1. Validate the version format"
    echo "2. Check that you're on the main branch"
    echo "3. Ensure working directory is clean"
    echo "4. Pull latest changes"
    echo "5. Run tests and linting"
    echo "6. Update version in source files"
    echo "7. Create a release commit"
    echo "8. Create and push the version tag"
    echo "9. Trigger the GitHub Actions release workflow"
}

# Main script logic
main() {
    # Check if help is requested
    if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
        show_usage
        exit 0
    fi
    
    # Check if version argument is provided
    if [[ $# -eq 0 ]]; then
        show_usage
        exit 1
    fi
    
    local version=$1
    
    # Show what we're about to do
    echo "=========================================="
    echo "Pomodux Release Script"
    echo "=========================================="
    echo "Version: $version"
    echo "Repository: $(git config --get remote.origin.url)"
    echo "Branch: $(git branch --show-current)"
    echo "=========================================="
    echo
    
    # Validate version
    validate_version "$version"
    
    # Check if tag exists
    check_tag_exists "$version"
    
    # Run pre-release checks
    run_pre_release_checks
    
    # Update version
    update_version "$version"
    
    # Create release commit
    create_release_commit "$version"
    
    # Create and push tag
    create_tag "$version"
    
    # Show next steps
    show_next_steps "$version"
}

# Run main function with all arguments
main "$@" 