#!/usr/bin/env bash
#
# Beads (beads) installation script
# Usage: curl -fsSL https://raw.githubusercontent.com/shaneholloman/beads/main/scripts/install.sh | bash
#
# ⚠️ IMPORTANT: This script must be EXECUTED, never SOURCED
# ❌ WRONG: source install.sh (will exit your shell on errors)
# ✅ CORRECT: bash install.sh
# ✅ CORRECT: curl -fsSL ... | bash
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}==>${NC} $1"
}

log_success() {
    echo -e "${GREEN}==>${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}==>${NC} $1"
}

log_error() {
    echo -e "${RED}Error:${NC} $1" >&2
}

# Detect OS and architecture
detect_platform() {
    local os arch

    case "$(uname -s)" in
        Darwin)
            os="darwin"
            ;;
        Linux)
            os="linux"
            ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)
            arch="amd64"
            ;;
        aarch64|arm64)
            arch="arm64"
            ;;
        *)
            log_error "Unsupported architecture: $(uname -m)"
            exit 1
            ;;
    esac

    echo "${os}_${arch}"
}

# Download and install from GitHub releases
install_from_release() {
    log_info "Installing beads from GitHub releases..."

    local platform=$1
    local tmp_dir
    tmp_dir=$(mktemp -d)

    # Get latest release version
    log_info "Fetching latest release..."
    local latest_url="https://api.github.com/repos/shaneholloman/beads/releases/latest"
    local version
    
    if command -v curl &> /dev/null; then
        version=$(curl -fsSL "$latest_url" | grep '"tag_name"' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    elif command -v wget &> /dev/null; then
        version=$(wget -qO- "$latest_url" | grep '"tag_name"' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    else
        log_error "Neither curl nor wget found. Please install one of them."
        return 1
    fi

    if [ -z "$version" ]; then
        log_error "Failed to fetch latest version"
        return 1
    fi

    log_info "Latest version: $version"

    # Download URL
    local archive_name="beads_${version#v}_${platform}.tar.gz"
    local download_url="https://github.com/shaneholloman/beads/releases/download/${version}/${archive_name}"
    
    log_info "Downloading $archive_name..."
    
    cd "$tmp_dir"
    if command -v curl &> /dev/null; then
        if ! curl -fsSL -o "$archive_name" "$download_url"; then
            log_error "Download failed"
            cd - > /dev/null || cd "$HOME"
            rm -rf "$tmp_dir"
            return 1
        fi
    elif command -v wget &> /dev/null; then
        if ! wget -q -O "$archive_name" "$download_url"; then
            log_error "Download failed"
            cd - > /dev/null || cd "$HOME"
            rm -rf "$tmp_dir"
            return 1
        fi
    fi

    # Extract archive
    log_info "Extracting archive..."
        cd - > /dev/null || cd "$HOME"
    if ! tar -xzf "$archive_name"; then
        log_error "Failed to extract archive"
        rm -rf "$tmp_dir"
        return 1
    fi

    # Determine install location
    local install_dir
    if [[ -w /usr/local/bin ]]; then
        install_dir="/usr/local/bin"
    else
        install_dir="$HOME/.local/bin"
        mkdir -p "$install_dir"
    fi

    # Install binary
    log_info "Installing to $install_dir..."
    if [[ -w "$install_dir" ]]; then
        mv beads "$install_dir/"
    else
        sudo mv beads "$install_dir/"
    fi

    log_success "beads installed to $install_dir/beads"

    # Check if install_dir is in PATH
    if [[ ":$PATH:" != *":$install_dir:"* ]]; then
        log_warning "$install_dir is not in your PATH"
        echo ""
        echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
        echo "  export PATH=\"\$PATH:$install_dir\""
        echo ""
    cd - > /dev/null || cd "$HOME"
    fi

    cd - > /dev/null
    rm -rf "$tmp_dir"
    return 0
}

# Check if Go is installed and meets minimum version
check_go() {
    if command -v go &> /dev/null; then
        local go_version=$(go version | awk '{print $3}' | sed 's/go//')
        log_info "Go detected: $(go version)"

        # Extract major and minor version numbers
    local major=$(echo "$go_version" | cut -d. -f1)
    local minor=$(echo "$go_version" | cut -d. -f2)

    # Check if Go version is 1.24 or later
    if [ "$major" -eq 1 ] && [ "$minor" -lt 24 ]; then
        log_error "Go 1.24 or later is required (found: $go_version)"
            echo ""
            echo "Please upgrade Go:"
            echo "  - Download from https://go.dev/dl/"
            echo "  - Or use your package manager to update"
            echo ""
            return 1
        fi

        return 0
    else
        return 1
    fi
}

# Install using go install (fallback)
install_with_go() {
    log_info "Installing beads using 'go install'..."

    if go install github.com/shaneholloman/beads/cmd/beads@latest; then
        log_success "beads installed successfully via go install"

        # Record where we expect the binary to have been installed
        # Prefer GOBIN if set, otherwise GOPATH/bin
        local gobin
        gobin=$(go env GOBIN 2>/dev/null || true)
        if [ -n "$gobin" ]; then
            bin_dir="$gobin"
        else
            bin_dir="$(go env GOPATH)/bin"
        fi
        LAST_INSTALL_PATH="$bin_dir/beads"

        # Check if GOPATH/bin (or GOBIN) is in PATH
        if [[ ":$PATH:" != *":$bin_dir:"* ]]; then
            log_warning "$bin_dir is not in your PATH"
            echo ""
            echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            echo "  export PATH=\"\$PATH:$bin_dir\""
            echo ""
        fi

        return 0
    else
        log_error "go install failed"
        return 1
    fi
}

# Build from source (last resort)
build_from_source() {
    log_info "Building beads from source..."

    local tmp_dir
    tmp_dir=$(mktemp -d)

    cd "$tmp_dir"
    log_info "Cloning repository..."

    if git clone --depth 1 https://github.com/shaneholloman/beads.git; then
        cd beads
        log_info "Building binary..."

        if go build -o beads ./cmd/beads; then
            # Determine install location
            local install_dir
            if [[ -w /usr/local/bin ]]; then
                install_dir="/usr/local/bin"
            else
                install_dir="$HOME/.local/bin"
                mkdir -p "$install_dir"
            fi

            log_info "Installing to $install_dir..."
            if [[ -w "$install_dir" ]]; then
                mv beads "$install_dir/"
            else
                sudo mv beads "$install_dir/"
            fi

            log_success "beads installed to $install_dir/beads"

            # Record where we installed the binary when building from source
            LAST_INSTALL_PATH="$install_dir/beads"

            # Check if install_dir is in PATH
            if [[ ":$PATH:" != *":$install_dir:"* ]]; then
                log_warning "$install_dir is not in your PATH"
                echo ""
                echo "Add this to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
                echo "  export PATH=\"\$PATH:$install_dir\""
        cd - > /dev/null || cd "$HOME"
                echo ""
            fi

            cd - > /dev/null
        cd - > /dev/null || cd "$HOME"
            rm -rf "$tmp_dir"
            return 0
        else
            log_error "Build failed"
    cd - > /dev/null || cd "$HOME"
            cd - > /dev/null
            rm -rf "$tmp_dir"
            return 1
        fi
    else
        log_error "Failed to clone repository"
        rm -rf "$tmp_dir"
        return 1
    fi
}

# Verify installation
verify_installation() {
    # If multiple 'beads' binaries exist on PATH, warn the user before verification
    warn_if_multiple_beads || true

    if command -v beads &> /dev/null; then
        log_success "beads is installed and ready!"
        echo ""
        beads version 2>/dev/null || echo "beads (development build)"
        echo ""
        echo "Get started:"
        echo "  cd your-project"
        echo "  beads init"
        echo "  beads quickstart"
        echo ""
        return 0
    else
        log_error "beads was installed but is not in PATH"
        return 1
    fi
}

# Returns a list of full paths to 'beads' found in PATH (earlier entries first)
get_beads_paths_in_path() {
    local IFS=':'
    local -a entries
    read -ra entries <<< "$PATH"
    local -a found
    local p
    for p in "${entries[@]}"; do
        [ -z "$p" ] && continue
        if [ -x "$p/beads" ]; then
            # Resolve symlink if possible
            if command -v readlink >/dev/null 2>&1; then
                resolved=$(readlink -f "$p/beads" 2>/dev/null || printf '%s' "$p/beads")
            else
                resolved="$p/beads"
            fi
            # avoid duplicates
            skip=0
            for existing in "${found[@]:-}"; do
                if [ "$existing" = "$resolved" ]; then skip=1; break; fi
            done
            if [ $skip -eq 0 ]; then
                found+=("$resolved")
            fi
        fi
    done
    # print results, one per line
    for item in "${found[@]:-}"; do
        printf '%s\n' "$item"
    done
}

warn_if_multiple_beads() {
    mapfile -t beads_paths < <(get_beads_paths_in_path)
    if [ "${#beads_paths[@]}" -le 1 ]; then
        return 0
    fi

    log_warning "Multiple 'beads' executables found on your PATH. An older copy may be executed instead of the one we installed."
    echo "Found the following 'beads' executables (entries earlier in PATH take precedence):"
    local i=1
    for p in "${beads_paths[@]}"; do
        local ver
        if [ -x "$p" ]; then
            ver=$("$p" version 2>/dev/null || true)
        fi
        if [ -z "$ver" ]; then ver="<unknown version>"; fi
        echo "  $i. $p  -> $ver"
        i=$((i+1))
    done

    if [ -n "$LAST_INSTALL_PATH" ]; then
        echo ""
        echo "We installed to: $LAST_INSTALL_PATH"
        # Compare first PATH entry vs installed path
        first="${beads_paths[0]}"
        if [ "$first" != "$LAST_INSTALL_PATH" ]; then
            log_warning "The 'beads' executable that appears first in your PATH is different from the one we installed. To make the newly installed 'beads' the one you get when running 'beads', either:"
            echo "  - Remove or rename the older $first from your PATH, or"
            echo "  - Reorder your PATH so that $(dirname "$LAST_INSTALL_PATH") appears before $(dirname "$first")"
            echo "After updating PATH, restart your shell and run 'beads version' to confirm."
        else
            echo "The installed 'beads' is first in your PATH.";
        fi
    else
        log_warning "We couldn't determine where we installed 'beads' during this run.";
    fi
}

# Main installation flow
main() {
    echo ""
    echo "🔗 Beads (beads) Installer"
    echo ""

    log_info "Detecting platform..."
    local platform
    platform=$(detect_platform)
    log_info "Platform: $platform"

    # Try downloading from GitHub releases first
    if install_from_release "$platform"; then
        verify_installation
        exit 0
    fi

    log_warning "Failed to install from releases, trying alternative methods..."

    # Try go install as fallback
    if check_go; then
        if install_with_go; then
            verify_installation
            exit 0
        fi
    fi

    # Try building from source as last resort
    log_warning "Falling back to building from source..."

    if ! check_go; then
        log_warning "Go is not installed"
        echo ""
        echo "beads requires Go 1.24 or later to build from source. You can:"
        echo "  1. Install Go from https://go.dev/dl/"
        echo "  2. Use your package manager:"
        echo "     - macOS: brew install go"
        echo "     - Ubuntu/Debian: sudo apt install golang"
        echo "     - Other Linux: Check your distro's package manager"
        echo ""
        echo "After installing Go, run this script again."
        exit 1
    fi

    if build_from_source; then
        verify_installation
        exit 0
    fi

    # All methods failed
    log_error "Installation failed"
    echo ""
    echo "Manual installation:"
    echo "  1. Download from https://github.com/shaneholloman/beads/releases/latest"
    echo "  2. Extract and move 'beads' to your PATH"
    echo ""
    echo "Or install from source:"
    echo "  1. Install Go from https://go.dev/dl/"
    echo "  2. Run: go install github.com/shaneholloman/beads/cmd/beads@latest"
    echo ""
    exit 1
}

main "$@"
