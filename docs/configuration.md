# Configuration Guide

This guide covers all configuration options available in the Network Configuration Checker, helping you customize the application for your specific environment and requirements.

## Configuration Overview

The application uses a hierarchical configuration system:

1. **Default Settings**: Built-in defaults for all options
2. **Configuration File**: Stored settings that persist between sessions
3. **Environment Variables**: Override settings for deployment scenarios
4. **Runtime Settings**: Temporary settings that don't persist

## Configuration File Location

The configuration file is automatically created and stored in:

- **Windows**: `%APPDATA%/NetworkConfigChecker/config.json`
- **macOS**: `~/Library/Application Support/NetworkConfigChecker/config.json`
- **Linux**: `~/.config/NetworkConfigChecker/config.json`

## Application Settings

### General Configuration

```json
{
  "application": {
    "theme": "system",
    "language": "en",
    "autoSave": true,
    "confirmDeletions": true,
    "showTooltips": true,
    "defaultPageSize": 20
  }
}
```

#### Options:

- **theme**: `"light"`, `"dark"`, or `"system"` (follows OS theme)
- **language**: Language code (currently supports `"en"`)
- **autoSave**: Automatically save changes without confirmation
- **confirmDeletions**: Show confirmation dialogs for delete operations
- **showTooltips**: Display helpful tooltips throughout the interface
- **defaultPageSize**: Number of items per page in lists (10, 20, 50, 100)

### Security Configuration

```json
{
  "security": {
    "passwordProtection": false,
    "sessionTimeout": 3600,
    "maxLoginAttempts": 3,
    "lockoutDuration": 900,
    "auditLogging": true,
    "encryptionAlgorithm": "AES-256-GCM",
    "keyDerivationIterations": 100000
  }
}
```

#### Options:

- **passwordProtection**: Enable application-level password protection
- **sessionTimeout**: Session timeout in seconds (0 = no timeout)
- **maxLoginAttempts**: Maximum failed login attempts before lockout
- **lockoutDuration**: Lockout duration in seconds after max attempts
- **auditLogging**: Enable comprehensive audit logging
- **encryptionAlgorithm**: Encryption algorithm for stored credentials
- **keyDerivationIterations**: PBKDF2 iterations for key derivation

### Connection Settings

```json
{
  "connections": {
    "sshTimeout": 30,
    "maxRetries": 3,
    "retryBackoff": "exponential",
    "maxConcurrentConnections": 10,
    "connectionPoolSize": 5,
    "keepAliveInterval": 60,
    "hostKeyVerification": true,
    "preferredAuthMethods": ["password", "publickey"]
  }
}
```

#### Options:

- **sshTimeout**: SSH connection timeout in seconds
- **maxRetries**: Maximum retry attempts for failed connections
- **retryBackoff**: Retry strategy (`"linear"`, `"exponential"`)
- **maxConcurrentConnections**: Maximum simultaneous device connections
- **connectionPoolSize**: Number of connections to keep in pool
- **keepAliveInterval**: SSH keep-alive interval in seconds
- **hostKeyVerification**: Verify SSH host keys (recommended for production)
- **preferredAuthMethods**: Ordered list of SSH authentication methods

### Data Management

```json
{
  "dataManagement": {
    "retentionPolicies": {
      "checkResults": "90d",
      "auditLogs": "365d",
      "reports": "180d",
      "tempFiles": "24h"
    },
    "autoCleanup": true,
    "cleanupSchedule": "daily",
    "backupEnabled": false,
    "backupLocation": "",
    "maxBackups": 7
  }
}
```

#### Options:

- **retentionPolicies**: How long to keep different types of data
- **autoCleanup**: Automatically clean up old data based on retention policies
- **cleanupSchedule**: When to run cleanup (`"daily"`, `"weekly"`, `"monthly"`)
- **backupEnabled**: Enable automatic database backups
- **backupLocation**: Directory for backup files (empty = default location)
- **maxBackups**: Maximum number of backup files to keep

## Device Configuration

### Default Device Settings

```json
{
  "deviceDefaults": {
    "sshPort": 22,
    "connectionTimeout": 30,
    "commandTimeout": 10,
    "enableSNMP": false,
    "snmpVersion": "2c",
    "snmpPort": 161,
    "tags": []
  }
}
```

### Vendor-Specific Settings

```json
{
  "vendors": {
    "cisco": {
      "enableMode": true,
      "enablePassword": "",
      "promptPattern": ".*[>#]\\s*$",
      "commandDelay": 100,
      "pageSize": 24
    },
    "juniper": {
      "cliMode": "operational",
      "promptPattern": ".*[>%]\\s*$",
      "commandDelay": 50,
      "pageSize": 0
    }
  }
}
```

## Security Check Configuration

### Check Engine Settings

```json
{
  "securityChecks": {
    "enabled": true,
    "parallelExecution": true,
    "maxWorkers": 5,
    "checkTimeout": 300,
    "retryFailedChecks": true,
    "severityLevels": ["critical", "high", "medium", "low", "info"],
    "autoAcknowledgeResolved": true
  }
}
```

### Custom Rules Configuration

```json
{
  "customRules": {
    "enabled": true,
    "rulesDirectory": "rules",
    "autoLoadRules": true,
    "validateRules": true,
    "allowUserRules": true
  }
}
```

## Reporting Configuration

### Report Generation

```json
{
  "reporting": {
    "defaultFormat": "pdf",
    "includeCharts": true,
    "includeLogo": true,
    "logoPath": "",
    "companyName": "",
    "reportFooter": "",
    "dateFormat": "YYYY-MM-DD",
    "timezone": "UTC"
  }
}
```

### Email Configuration

```json
{
  "email": {
    "enabled": false,
    "smtpServer": "",
    "smtpPort": 587,
    "username": "",
    "password": "",
    "useTLS": true,
    "fromAddress": "",
    "fromName": "Network Config Checker",
    "replyTo": ""
  }
}
```

## Logging Configuration

### Log Settings

```json
{
  "logging": {
    "level": "info",
    "format": "json",
    "output": "file",
    "logFile": "app.log",
    "maxFileSize": "100MB",
    "maxFiles": 10,
    "compress": true,
    "includeStackTrace": false
  }
}
```

#### Log Levels:

- **debug**: Detailed debugging information
- **info**: General information messages
- **warn**: Warning messages for potential issues
- **error**: Error messages for failures
- **fatal**: Critical errors that cause application shutdown

## Performance Tuning

### Memory Settings

```json
{
  "performance": {
    "maxMemoryUsage": "1GB",
    "cacheSize": "100MB",
    "enableCaching": true,
    "cacheTimeout": 300,
    "gcInterval": 60
  }
}
```

### Database Optimization

```json
{
  "database": {
    "pragmas": {
      "journal_mode": "WAL",
      "synchronous": "NORMAL",
      "cache_size": 10000,
      "temp_store": "MEMORY"
    },
    "connectionPool": {
      "maxOpenConnections": 10,
      "maxIdleConnections": 5,
      "connectionMaxLifetime": "1h"
    }
  }
}
```

## Environment Variables

Override configuration settings using environment variables:

```bash
# Security settings
export NCC_PASSWORD_PROTECTION=true
export NCC_SESSION_TIMEOUT=7200

# Connection settings
export NCC_SSH_TIMEOUT=45
export NCC_MAX_CONCURRENT=15

# Database settings
export NCC_DB_PATH="/custom/path/database.db"

# Logging settings
export NCC_LOG_LEVEL=debug
export NCC_LOG_FILE="/var/log/network-config-checker.log"
```

## Configuration Validation

The application validates all configuration settings on startup:

### Validation Rules

- **Numeric values**: Must be within acceptable ranges
- **File paths**: Must be accessible and writable
- **Network settings**: Must be valid IP addresses and ports
- **Email settings**: Must be valid email addresses and SMTP configuration

### Error Handling

- **Invalid settings**: Application uses defaults and logs warnings
- **Missing files**: Application creates default configuration
- **Permission errors**: Application reports errors and suggests fixes

## Configuration Management

### Backup Configuration

```bash
# Create backup
cp ~/.config/NetworkConfigChecker/config.json config-backup.json

# Restore from backup
cp config-backup.json ~/.config/NetworkConfigChecker/config.json
```

### Reset to Defaults

1. Stop the application
2. Delete or rename the configuration file
3. Restart the application (creates new default configuration)

### Export/Import Settings

Use the Settings interface:

1. Go to **Settings** → **Advanced**
2. Click **"Export Configuration"** to save current settings
3. Click **"Import Configuration"** to load saved settings

## Advanced Configuration

### Custom Themes

Create custom themes by modifying CSS variables:

```json
{
  "theme": {
    "custom": true,
    "colors": {
      "primary": "#007acc",
      "secondary": "#6c757d",
      "success": "#28a745",
      "warning": "#ffc107",
      "danger": "#dc3545"
    }
  }
}
```

### Plugin Configuration

Enable and configure plugins:

```json
{
  "plugins": {
    "enabled": true,
    "pluginDirectory": "plugins",
    "autoLoad": true,
    "allowThirdParty": false,
    "sandboxed": true
  }
}
```

## Troubleshooting Configuration

### Common Issues

#### Configuration Not Loading

- Check file permissions
- Verify JSON syntax
- Review application logs

#### Settings Not Persisting

- Ensure write permissions to config directory
- Check disk space
- Verify application has proper permissions

#### Performance Issues

- Reduce concurrent connections
- Increase timeouts for slow networks
- Enable caching for better performance

### Diagnostic Commands

Check configuration status:

```bash
# View current configuration
network-config-checker --show-config

# Validate configuration
network-config-checker --validate-config

# Reset to defaults
network-config-checker --reset-config
```

## Best Practices

### Security

- Enable password protection in production
- Use strong encryption settings
- Enable audit logging
- Regularly review and rotate credentials

### Performance

- Tune concurrent connections based on network capacity
- Set appropriate timeouts for your environment
- Enable caching for frequently accessed data
- Monitor resource usage and adjust limits

### Maintenance

- Regularly backup configuration
- Review and update retention policies
- Monitor log files for errors
- Keep configuration documentation updated

## Next Steps

After configuring the application:

1. **Test Configuration**: Verify all settings work as expected
2. **Monitor Performance**: Watch for any performance issues
3. **Review Security**: Ensure security settings meet your requirements
4. **Document Changes**: Keep track of configuration changes
5. **Train Users**: Ensure users understand any configuration-specific behaviors
