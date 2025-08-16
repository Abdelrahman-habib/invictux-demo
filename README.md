# Network Configuration Checker

A cross-platform desktop application built with Wails v2 that automates network device security audits and configuration compliance checking across enterprise infrastructure. Reduce manual network device security auditing from 4+ hours to 5 minutes with automated compliance reporting.

## Features

### 🔧 Device Management

- **Comprehensive Device Support**: Manage routers, switches, firewalls, and access points from major vendors (Cisco, Juniper, HP, Aruba)
- **Secure Credential Storage**: AES-256 encrypted password storage with proper key management
- **Connectivity Testing**: Real-time device reachability and SSH port accessibility verification
- **Bulk Operations**: Import devices from CSV, run checks on multiple devices simultaneously
- **Advanced Filtering**: Search and filter devices by name, IP, type, vendor, or status

### 🛡️ Security Auditing

- **Automated Security Checks**: Predefined security rules for each vendor with customizable rule sets
- **Real-time Monitoring**: Live progress tracking during security check execution
- **Comprehensive Results**: Categorized results (PASS/FAIL/WARNING/ERROR) with detailed evidence
- **Parallel Processing**: Configurable concurrency limits for efficient bulk checking
- **Retry Logic**: Automatic retry with exponential backoff for transient failures

### 📊 Dashboard & Monitoring

- **Executive Overview**: Real-time dashboard with device counts, status summaries, and critical issues
- **Visual Analytics**: Interactive pie charts showing security status distribution
- **Device Status Grid**: Color-coded device tiles with hover tooltips and quick navigation
- **Auto-refresh**: Automatic data updates every 30 seconds

### 🚨 Issue Management

- **Centralized Issue Tracking**: Filterable list of all security issues with severity levels
- **Detailed Evidence**: Complete check results with remediation suggestions
- **Issue Lifecycle**: Acknowledgment tracking and automatic status updates
- **Export Capabilities**: CSV export for analysis and reporting

### 📈 Reporting System

- **Multiple Report Types**: Executive summaries and detailed technical reports
- **Flexible Export**: PDF for management, CSV for analysis
- **Automated Scheduling**: Recurring report generation with email delivery
- **Historical Archive**: Complete report history with search capabilities
- **Custom Templates**: Customizable report sections and branding

### ⚙️ Configuration Management

- **Application Settings**: Configurable check intervals, timeouts, and concurrency limits
- **Security Configuration**: Session timeouts, password protection, and audit logging
- **Vendor Customization**: Custom command mappings and security check templates
- **Data Retention**: Automatic cleanup of old results and reports

## Technology Stack

- **Frontend**: React 18 + TypeScript + TanStack Router + TanStack Query
- **Backend**: Go with Wails v2.9.2 framework
- **Database**: SQLite with automatic migrations
- **UI Components**: shadcn/ui with Tailwind CSS
- **Forms**: React Hook Form with Zod validation
- **Security**: AES-256 encryption, SSH protocol, secure session management

## Quick Start

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Wails CLI v2.9.2

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/your-org/network-config-checker.git
   cd network-config-checker
   ```

2. **Install dependencies**

   ```bash
   # Install Go dependencies
   go mod download

   # Install frontend dependencies
   cd frontend
   npm install
   cd ..
   ```

3. **Build the application**

   ```bash
   # Development build
   wails build -debug

   # Production build
   wails build
   ```

4. **Run the application**

   ```bash
   # Development mode with hot reload
   wails dev

   # Or run the built executable
   ./build/bin/network-config-checker
   ```

## Usage

### Adding Devices

1. Navigate to the **Devices** page
2. Click **Add Device** to open the device form
3. Fill in device details (name, IP, credentials, etc.)
4. Test connectivity before saving
5. Device will appear in the device grid with status indicator

### Running Security Checks

1. Select devices from the device grid
2. Click **Run Security Checks** for bulk operations
3. Monitor progress in real-time
4. View results in the **Security Issues** page

### Generating Reports

1. Go to the **Reports** page
2. Select report type (Executive Summary or Technical Details)
3. Choose date range and device scope
4. Generate and download in PDF or CSV format

### Dashboard Monitoring

- View the **Dashboard** for real-time overview
- Monitor device status with color-coded tiles
- Track critical issues and compliance metrics
- Click device tiles for detailed information

## Configuration

### Application Settings

Access settings through the application menu to configure:

- Check scheduling intervals
- Connection timeouts and retry limits
- Concurrent check limits
- Data retention policies

### Security Settings

- Enable application password protection
- Configure session timeout duration
- Set up audit logging preferences
- Manage credential encryption settings

### Vendor Configuration

- Customize device command mappings
- Define custom security check rules
- Configure vendor-specific templates
- Import/export rule sets

## Security Features

- **Credential Protection**: All passwords encrypted with AES-256
- **Secure Communication**: SSH protocol for device connections
- **Audit Logging**: Complete activity logging for security operations
- **Session Management**: Automatic session timeout and cleanup
- **Input Validation**: Comprehensive input sanitization and validation
- **Host Key Verification**: SSH host key verification in production

## Development

### Project Structure

```
├── frontend/                 # React frontend application
│   ├── src/
│   │   ├── features/        # Feature-based components
│   │   ├── components/      # Shared UI components
│   │   ├── hooks/          # Custom React hooks
│   │   └── services/       # API communication layer
├── internal/               # Go backend modules
│   ├── app/               # Wails application context
│   ├── device/            # Device management
│   ├── checker/           # Security check engine
│   ├── database/          # Database operations
│   └── security/          # Security utilities
├── docs/                  # Documentation
└── scripts/              # Build and deployment scripts
```

### Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Testing

```bash
# Run backend tests
go test ./...

# Run frontend tests
cd frontend
npm test

# Run end-to-end tests
npm run test:e2e
```

## Deployment

### Building for Distribution

```bash
# Windows
./scripts/build-windows.sh

# macOS (Intel)
./scripts/build-macos-intel.sh

# macOS (ARM)
./scripts/build-macos-arm.sh

# Linux
./scripts/build.sh
```

### Installation Packages

- **Windows**: MSI installer with automatic updates
- **macOS**: DMG package with code signing
- **Linux**: DEB/RPM packages for major distributions

## Documentation

Comprehensive documentation is available in the `docs/` folder:

### Getting Started

- [Installation Guide](docs/installation.md) - Complete installation instructions
- [Quick Start Guide](docs/quick-start.md) - Get up and running in 5 minutes
- [Configuration Guide](docs/configuration.md) - Detailed configuration options

### User Guides

- [Device Management](docs/user-guide/device-management.md) - Managing network devices
- [Security Checks](docs/user-guide/security-checks.md) - Running and understanding security audits
- [Dashboard](docs/user-guide/dashboard.md) - Real-time monitoring and overview
- [Issue Management](docs/user-guide/issue-management.md) - Tracking and resolving security issues
- [Reports](docs/user-guide/reports.md) - Generating compliance and technical reports

### Administration

- [Troubleshooting Guide](docs/admin/troubleshooting.md) - Common issues and solutions
- [Security Policy](SECURITY.md) - Security features and vulnerability reporting

## Support

- **Documentation**: See the `docs/` folder for detailed guides
- **Issues**: Report bugs and feature requests on GitHub Issues
- **Security**: Report security vulnerabilities to security@yourorg.com

## Roadmap

- [ ] Support for additional vendors (Fortinet, Palo Alto, etc.)
- [ ] REST API for integration with external systems
- [ ] Advanced analytics and trend analysis
- [ ] Mobile companion app for alerts
- [ ] Cloud-based device discovery
- [ ] Integration with SIEM systems

## License

MIT License - see LICENSE file for details
