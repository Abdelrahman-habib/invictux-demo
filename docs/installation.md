# Installation Guide

This guide will walk you through installing the Network Configuration Checker on your system.

## System Requirements

### Minimum Requirements

- **Operating System**: Windows 10, macOS 10.15, or Linux (Ubuntu 18.04+)
- **RAM**: 4 GB minimum, 8 GB recommended
- **Storage**: 500 MB free disk space
- **Network**: TCP/IP connectivity to target network devices

### Recommended Requirements

- **RAM**: 8 GB or more for large device inventories
- **Storage**: 2 GB for reports and historical data
- **CPU**: Multi-core processor for parallel security checks

## Installation Methods

### Option 1: Pre-built Binaries (Recommended)

#### Windows

1. Download the latest MSI installer from the [releases page](https://github.com/your-org/network-config-checker/releases)
2. Run the installer as Administrator
3. Follow the installation wizard
4. Launch from Start Menu or Desktop shortcut

#### macOS

1. Download the DMG file from the [releases page](https://github.com/your-org/network-config-checker/releases)
2. Open the DMG file
3. Drag the application to your Applications folder
4. Launch from Applications or Spotlight

#### Linux

**Ubuntu/Debian:**

```bash
# Download the DEB package
wget https://github.com/your-org/network-config-checker/releases/latest/download/network-config-checker.deb

# Install the package
sudo dpkg -i network-config-checker.deb
sudo apt-get install -f  # Fix any dependency issues

# Launch the application
network-config-checker
```

**RHEL/CentOS/Fedora:**

```bash
# Download the RPM package
wget https://github.com/your-org/network-config-checker/releases/latest/download/network-config-checker.rpm

# Install the package
sudo rpm -i network-config-checker.rpm

# Launch the application
network-config-checker
```

### Option 2: Build from Source

#### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Wails CLI v2.9.2

#### Installation Steps

1. **Install Wails CLI**

   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

2. **Clone the repository**

   ```bash
   git clone https://github.com/your-org/network-config-checker.git
   cd network-config-checker
   ```

3. **Install dependencies**

   ```bash
   # Backend dependencies
   go mod download

   # Frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

4. **Build the application**

   ```bash
   # Development build
   wails build -debug

   # Production build
   wails build
   ```

5. **Run the application**

   ```bash
   # Development mode with hot reload
   wails dev

   # Or run the built executable
   ./build/bin/network-config-checker
   ```

## Post-Installation Setup

### First Launch

1. Launch the Network Configuration Checker
2. The application will create its data directory:
   - **Windows**: `%APPDATA%/NetworkConfigChecker`
   - **macOS**: `~/Library/Application Support/NetworkConfigChecker`
   - **Linux**: `~/.config/NetworkConfigChecker`

### Initial Configuration

1. **Set Application Password** (Recommended)

   - Go to Settings → Security
   - Enable "Password Protection"
   - Set a strong password

2. **Configure Connection Settings**

   - Go to Settings → Connections
   - Set appropriate timeouts for your network
   - Configure retry limits

3. **Set Data Retention Policies**
   - Go to Settings → Data Management
   - Configure how long to keep check results
   - Set up automatic cleanup schedules

## Verification

### Test Installation

1. Launch the application
2. Navigate to the Dashboard
3. Verify all UI elements load correctly
4. Check that the database initializes properly

### Test Device Connectivity

1. Go to Devices → Add Device
2. Add a test device with valid credentials
3. Use the "Test Connectivity" feature
4. Verify successful connection

## Troubleshooting Installation Issues

### Common Issues

#### Windows: "Application failed to start"

- **Cause**: Missing Visual C++ Redistributables
- **Solution**: Install Microsoft Visual C++ Redistributable packages

#### macOS: "App is damaged and can't be opened"

- **Cause**: Gatekeeper security restrictions
- **Solution**:
  ```bash
  sudo xattr -rd com.apple.quarantine /Applications/NetworkConfigChecker.app
  ```

#### Linux: "Permission denied"

- **Cause**: Executable permissions not set
- **Solution**:
  ```bash
  chmod +x network-config-checker
  ```

#### Database initialization fails

- **Cause**: Insufficient permissions or disk space
- **Solution**:
  - Check available disk space
  - Verify write permissions to data directory
  - Run as administrator/sudo if necessary

### Getting Help

If you encounter issues not covered here:

1. Check the [Troubleshooting Guide](admin/troubleshooting.md)
2. Review the application logs:
   - **Windows**: `%APPDATA%/NetworkConfigChecker/logs`
   - **macOS**: `~/Library/Application Support/NetworkConfigChecker/logs`
   - **Linux**: `~/.config/NetworkConfigChecker/logs`
3. Create an issue on [GitHub](https://github.com/your-org/network-config-checker/issues)

## Uninstallation

### Windows

1. Use "Add or Remove Programs" from Windows Settings
2. Select "Network Configuration Checker"
3. Click "Uninstall"

### macOS

1. Drag the application from Applications to Trash
2. Remove application data (optional):
   ```bash
   rm -rf ~/Library/Application\ Support/NetworkConfigChecker
   ```

### Linux

```bash
# Ubuntu/Debian
sudo apt remove network-config-checker

# RHEL/CentOS/Fedora
sudo rpm -e network-config-checker

# Remove application data (optional)
rm -rf ~/.config/NetworkConfigChecker
```

## Next Steps

After successful installation:

1. Read the [Quick Start Guide](quick-start.md) to get up and running
2. Follow the [Configuration Guide](configuration.md) to customize settings
3. Start with the [Device Management Guide](user-guide/device-management.md) to add your first devices
