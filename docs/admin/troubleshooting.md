# Troubleshooting Guide

This guide helps administrators diagnose and resolve common issues with the Network Configuration Checker.

## Common Issues

### Installation Problems

#### Application Won't Start

**Windows: "Application failed to start"**

- **Cause**: Missing Visual C++ Redistributables
- **Solution**: Install Microsoft Visual C++ Redistributable packages
- **Download**: https://aka.ms/vs/17/release/vc_redist.x64.exe

**macOS: "App is damaged and can't be opened"**

- **Cause**: Gatekeeper security restrictions
- **Solution**:
  ```bash
  sudo xattr -rd com.apple.quarantine /Applications/NetworkConfigChecker.app
  ```

**Linux: "Permission denied"**

- **Cause**: Executable permissions not set
- **Solution**:
  ```bash
  chmod +x network-config-checker
  ```

#### Database Initialization Fails

**Symptoms**: Application crashes on startup with database errors

**Causes and Solutions**:

1. **Insufficient Permissions**

   - Check write permissions to data directory
   - Run as administrator/sudo if necessary

2. **Disk Space**

   - Verify available disk space (minimum 100MB)
   - Clean up temporary files

3. **Corrupted Database**
   - Delete database file to force recreation
   - Restore from backup if available

### Connection Issues

#### SSH Connection Failures

**"Connection refused" Errors**

- **Check SSH Service**: Verify SSH is enabled on target device
- **Port Configuration**: Confirm SSH port (default 22) is correct
- **Firewall Rules**: Ensure SSH port is open in firewalls
- **Network Connectivity**: Test basic network connectivity with ping

**Authentication Failures**

- **Verify Credentials**: Double-check username and password
- **Account Status**: Ensure account is not locked or disabled
- **Privilege Levels**: Confirm account has sufficient privileges
- **Key-based Auth**: If using SSH keys, verify key format and permissions

**Connection Timeouts**

- **Network Latency**: Increase timeout values for slow networks
- **Device Load**: Check if device is overloaded
- **Concurrent Connections**: Reduce parallel connection limits
- **MTU Issues**: Check for MTU/fragmentation problems

#### Device Connectivity Issues

**Device Shows as Offline**

1. **Network Reachability**

   ```bash
   ping [device-ip]
   telnet [device-ip] 22
   ```

2. **DNS Resolution**

   - Use IP addresses instead of hostnames
   - Verify DNS configuration

3. **Routing Issues**
   - Check routing tables
   - Verify VLAN configuration
   - Test from different network segments

### Performance Issues

#### Slow Application Performance

**High Memory Usage**

- **Symptoms**: Application uses excessive RAM
- **Solutions**:
  - Reduce concurrent connection limits
  - Implement data retention policies
  - Restart application periodically
  - Increase system RAM if needed

**Slow Database Operations**

- **Symptoms**: Long delays when loading data
- **Solutions**:
  - Rebuild database indexes
  - Implement database maintenance schedule
  - Consider database optimization
  - Monitor disk I/O performance

**Network Bottlenecks**

- **Symptoms**: Slow security checks, timeouts
- **Solutions**:
  - Reduce parallel check execution
  - Implement bandwidth throttling
  - Schedule checks during off-peak hours
  - Optimize network infrastructure

#### Security Check Performance

**Checks Taking Too Long**

1. **Optimize Concurrency**

   - Reduce max concurrent connections
   - Stagger check execution
   - Monitor device CPU usage

2. **Rule Optimization**

   - Review complex rules
   - Combine related checks
   - Cache frequently used data

3. **Network Optimization**
   - Use wired connections when possible
   - Minimize network hops
   - Optimize routing

### Security Check Issues

#### Checks Failing Consistently

**Rule Pattern Mismatches**

- **Symptoms**: False positives or missed issues
- **Solutions**:
  - Review rule patterns
  - Test against known configurations
  - Update patterns for vendor variations
  - Add exception handling

**Command Execution Failures**

- **Symptoms**: Commands fail on devices
- **Solutions**:
  - Verify command syntax for device type
  - Check user privilege levels
  - Test commands manually
  - Update command templates

**Incomplete Results**

- **Symptoms**: Some checks don't run
- **Solutions**:
  - Check rule dependencies
  - Verify device compatibility
  - Review error logs
  - Update rule definitions

### Data and Configuration Issues

#### Configuration Not Persisting

**Settings Reset on Restart**

- **Cause**: Configuration file permissions or corruption
- **Solutions**:
  - Check file permissions on config directory
  - Verify disk space for config writes
  - Restore from configuration backup
  - Reset to default configuration

**Device Data Loss**

- **Cause**: Database corruption or migration issues
- **Solutions**:
  - Restore from database backup
  - Re-import device data from CSV
  - Check database integrity
  - Implement regular backups

#### Report Generation Issues

**PDF Generation Fails**

- **Symptoms**: Reports fail to generate or are corrupted
- **Solutions**:
  - Check available disk space
  - Verify report template integrity
  - Update PDF generation libraries
  - Test with smaller data sets

**Email Delivery Problems**

- **Symptoms**: Scheduled reports not delivered
- **Solutions**:
  - Verify SMTP configuration
  - Check email server connectivity
  - Test with different email providers
  - Review email logs

## Diagnostic Tools

### Log Analysis

#### Application Logs

**Log Locations**:

- **Windows**: `%APPDATA%/NetworkConfigChecker/logs`
- **macOS**: `~/Library/Application Support/NetworkConfigChecker/logs`
- **Linux**: `~/.config/NetworkConfigChecker/logs`

**Log Levels**:

- **ERROR**: Critical errors requiring attention
- **WARN**: Warnings that may indicate problems
- **INFO**: General information messages
- **DEBUG**: Detailed debugging information

#### Log Analysis Commands

```bash
# View recent errors
grep "ERROR" app.log | tail -20

# Monitor logs in real-time
tail -f app.log

# Search for specific issues
grep -i "connection" app.log | grep "failed"

# Analyze performance issues
grep "timeout\|slow\|performance" app.log
```

### System Diagnostics

#### Resource Monitoring

**Memory Usage**:

```bash
# Windows
tasklist /fi "imagename eq network-config-checker.exe"

# macOS/Linux
ps aux | grep network-config-checker
top -p $(pgrep network-config-checker)
```

**Disk Usage**:

```bash
# Check data directory size
du -sh ~/.config/NetworkConfigChecker

# Check available space
df -h
```

**Network Connectivity**:

```bash
# Test device connectivity
ping -c 4 [device-ip]
telnet [device-ip] 22
nmap -p 22 [device-ip]
```

### Database Diagnostics

#### Database Integrity

```bash
# Check database file
file ~/.config/NetworkConfigChecker/database.db

# Basic integrity check (if SQLite tools available)
sqlite3 database.db "PRAGMA integrity_check;"

# Check database size
ls -lh ~/.config/NetworkConfigChecker/database.db
```

#### Database Maintenance

```bash
# Backup database
cp database.db database.db.backup

# Vacuum database (optimize)
sqlite3 database.db "VACUUM;"

# Analyze database statistics
sqlite3 database.db "ANALYZE;"
```

## Advanced Troubleshooting

### Debug Mode

#### Enabling Debug Logging

1. **Through Settings**:

   - Go to Settings → Advanced
   - Set Log Level to "Debug"
   - Restart application

2. **Environment Variable**:

   ```bash
   export NCC_LOG_LEVEL=debug
   network-config-checker
   ```

3. **Command Line**:
   ```bash
   network-config-checker --debug
   ```

#### Debug Information Collection

**System Information**:

- Operating system and version
- Application version and build
- Available memory and disk space
- Network configuration

**Application State**:

- Configuration file contents
- Database schema version
- Active connections
- Recent error messages

### Network Debugging

#### SSH Connection Testing

```bash
# Test SSH connectivity manually
ssh -v username@device-ip

# Test with specific port
ssh -p 2222 username@device-ip

# Test with timeout
timeout 30 ssh username@device-ip
```

#### Network Trace Analysis

```bash
# Capture network traffic (requires admin privileges)
tcpdump -i any host [device-ip] and port 22

# Windows equivalent
netsh trace start capture=yes provider=Microsoft-Windows-TCPIP
```

### Performance Profiling

#### Application Performance

1. **CPU Profiling**:

   - Monitor CPU usage during operations
   - Identify performance bottlenecks
   - Optimize resource-intensive operations

2. **Memory Profiling**:

   - Track memory allocation patterns
   - Identify memory leaks
   - Optimize memory usage

3. **I/O Profiling**:
   - Monitor disk and network I/O
   - Identify I/O bottlenecks
   - Optimize data access patterns

## Recovery Procedures

### Data Recovery

#### Database Recovery

1. **From Backup**:

   ```bash
   # Stop application
   # Restore database
   cp database.db.backup database.db
   # Restart application
   ```

2. **Rebuild from Scratch**:
   ```bash
   # Stop application
   # Remove corrupted database
   rm database.db
   # Restart application (creates new database)
   # Re-import device data
   ```

#### Configuration Recovery

1. **Reset to Defaults**:

   ```bash
   # Stop application
   # Remove configuration
   rm config.json
   # Restart application
   ```

2. **Restore from Backup**:
   ```bash
   # Stop application
   # Restore configuration
   cp config.json.backup config.json
   # Restart application
   ```

### System Recovery

#### Complete Reinstallation

1. **Backup Data**:

   ```bash
   # Backup important data
   cp -r ~/.config/NetworkConfigChecker ~/NetworkConfigChecker.backup
   ```

2. **Uninstall Application**:

   - Use system uninstaller
   - Remove application directories
   - Clean registry entries (Windows)

3. **Clean Installation**:
   - Download latest version
   - Install with default settings
   - Restore data from backup

## Prevention Strategies

### Regular Maintenance

#### Scheduled Tasks

1. **Database Maintenance**:

   - Weekly database vacuum
   - Monthly integrity checks
   - Quarterly full backups

2. **Log Management**:

   - Daily log rotation
   - Weekly log cleanup
   - Monthly log analysis

3. **Performance Monitoring**:
   - Daily resource usage checks
   - Weekly performance reports
   - Monthly optimization reviews

#### Monitoring Setup

1. **Application Monitoring**:

   - Set up health checks
   - Monitor key metrics
   - Configure alerting

2. **System Monitoring**:
   - Monitor system resources
   - Track application performance
   - Set up automated notifications

### Best Practices

#### Operational Excellence

1. **Documentation**:

   - Maintain troubleshooting logs
   - Document configuration changes
   - Keep vendor contact information

2. **Training**:

   - Regular training on troubleshooting procedures
   - Knowledge sharing sessions
   - Incident response drills

3. **Continuous Improvement**:
   - Regular review of issues and solutions
   - Process optimization
   - Tool and procedure updates

## Getting Additional Help

### Support Channels

1. **Documentation**: Check all available documentation first
2. **Community Forums**: Search community discussions
3. **GitHub Issues**: Report bugs and feature requests
4. **Professional Support**: Contact support team for critical issues

### Information to Provide

When seeking help, include:

1. **System Information**:

   - Operating system and version
   - Application version
   - Hardware specifications

2. **Problem Description**:

   - Detailed problem description
   - Steps to reproduce
   - Expected vs actual behavior

3. **Diagnostic Information**:

   - Relevant log entries
   - Configuration details
   - Error messages

4. **Environment Details**:
   - Network configuration
   - Device types and vendors
   - Security policies
