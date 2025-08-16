# Device Management Guide

This guide covers all aspects of managing network devices in the Network Configuration Checker, from adding individual devices to bulk operations and advanced management features.

## Overview

The Device Management system allows you to:

- Add, edit, and delete network devices
- Import devices from CSV files
- Test device connectivity
- Organize devices with tags and filters
- Perform bulk operations on multiple devices
- Monitor device status and health

## Adding Devices

### Single Device Addition

1. **Navigate to Devices Page**

   - Click **"Devices"** in the main navigation
   - Click the **"Add Device"** button

2. **Fill Device Information**

   ```
   Device Name: Core-Switch-01
   IP Address: 192.168.1.10
   Device Type: Switch
   Vendor: Cisco
   Username: admin
   Password: [secure password]
   SSH Port: 22 (default)
   SNMP Community: public (optional)
   Tags: core, datacenter, critical
   ```

3. **Test Connectivity**

   - Click **"Test Connectivity"** before saving
   - Verify successful connection
   - Address any connection issues

4. **Save Device**
   - Click **"Save Device"** to add to inventory
   - Device appears in the device grid with status indicator

### Device Information Fields

#### Required Fields

- **Device Name**: Unique identifier for the device
- **IP Address**: IPv4 address of the device
- **Device Type**: Router, Switch, Firewall, or Access Point
- **Vendor**: Cisco, Juniper, HP, or Aruba
- **Username**: SSH username for device access
- **Password**: SSH password (encrypted when stored)

#### Optional Fields

- **SSH Port**: Default is 22, change if using non-standard port
- **SNMP Community**: For SNMP-based monitoring (if enabled)
- **Tags**: Comma-separated tags for organization

### Bulk Device Import

For adding multiple devices efficiently:

1. **Download CSV Template**

   - Go to **Devices** → **Import**
   - Click **"Download Template"**

2. **Prepare CSV File**

   ```csv
   name,ip_address,device_type,vendor,username,password,ssh_port,snmp_community,tags
   Router-01,192.168.1.1,router,cisco,admin,password123,22,public,"core,critical"
   Switch-01,192.168.1.10,switch,cisco,admin,password123,22,public,"access,floor1"
   Firewall-01,192.168.1.254,firewall,cisco,admin,password123,22,,"security,perimeter"
   ```

3. **Import Devices**

   - Click **"Choose File"** and select your CSV
   - Review the preview of devices to be imported
   - Click **"Import Devices"**
   - Monitor import progress and address any errors

4. **Verify Import**
   - Check that all devices appear in the device grid
   - Test connectivity for critical devices
   - Review any import warnings or errors

## Device Management Interface

### Device Grid View

The main device interface shows devices as cards with:

#### Device Card Information

- **Device Name**: Primary identifier
- **IP Address**: Network address
- **Status Indicator**: Online (green), Offline (red), Unknown (gray)
- **Device Type & Vendor**: Icons and labels
- **Last Checked**: Timestamp of last security check
- **Quick Actions**: Edit, Delete, Run Check buttons

#### Status Indicators

- 🟢 **Online**: Device is reachable and responsive
- 🔴 **Offline**: Device is not reachable
- ⚪ **Unknown**: Status not yet determined
- 🟡 **Warning**: Device reachable but has issues

### Search and Filtering

#### Search Functionality

- **Global Search**: Search across device names, IP addresses, and tags
- **Real-time Results**: Results update as you type
- **Search Persistence**: Search terms persist across page navigation

#### Filter Options

- **Device Type**: Filter by Router, Switch, Firewall, Access Point
- **Vendor**: Filter by Cisco, Juniper, HP, Aruba
- **Status**: Filter by Online, Offline, Unknown
- **Tags**: Filter by specific tags

#### Advanced Filtering

1. Click the **Filter** button
2. Select multiple filter criteria
3. Use **"Apply Filters"** to update results
4. Use **"Clear Filters"** to reset

### Sorting Options

Sort devices by:

- **Name** (A-Z or Z-A)
- **IP Address** (ascending or descending)
- **Status** (Online first or Offline first)
- **Last Checked** (newest or oldest first)
- **Device Type** (alphabetical)

## Device Operations

### Individual Device Actions

#### Edit Device

1. Click the **three-dot menu** on a device card
2. Select **"Edit Device"**
3. Modify any device properties
4. Test connectivity if credentials changed
5. Click **"Save Changes"**

#### Delete Device

1. Click the **three-dot menu** on a device card
2. Select **"Delete Device"**
3. Confirm deletion in the dialog
4. Device and all associated data are removed

#### Test Connectivity

1. Click the **three-dot menu** on a device card
2. Select **"Test Connectivity"**
3. Monitor the test progress
4. Review connection results

#### Run Security Check

1. Click the **three-dot menu** on a device card
2. Select **"Run Security Check"**
3. Monitor check progress
4. View results when complete

### Bulk Operations

#### Selecting Multiple Devices

1. Enable **Selection Mode** using the toggle
2. Click device cards to select them
3. Selected devices show with checkmarks
4. Use **"Select All"** or **"Clear Selection"** as needed

#### Bulk Actions

With devices selected:

- **Run Security Checks**: Execute checks on all selected devices
- **Delete Devices**: Remove multiple devices at once
- **Export Selection**: Export selected devices to CSV
- **Apply Tags**: Add or remove tags from selected devices

#### Bulk Security Checks

1. Select multiple devices
2. Click **"Run Security Checks"**
3. Configure check options:
   - **Parallel Execution**: Run checks simultaneously
   - **Max Concurrent**: Limit simultaneous connections
   - **Stop on Error**: Stop if any device fails
4. Monitor progress for all devices
5. Review aggregated results

## Device Organization

### Using Tags

Tags help organize and categorize devices:

#### Adding Tags

- **During Creation**: Add tags in the device form
- **After Creation**: Edit device to add/modify tags
- **Bulk Tagging**: Apply tags to multiple selected devices

#### Tag Examples

- **Location**: `datacenter`, `branch-office`, `floor1`
- **Criticality**: `critical`, `important`, `standard`
- **Function**: `core`, `access`, `distribution`
- **Environment**: `production`, `staging`, `development`

#### Tag Management

- **View All Tags**: See all tags used across devices
- **Rename Tags**: Update tag names globally
- **Delete Tags**: Remove unused tags
- **Tag Statistics**: See how many devices use each tag

### Device Groups

Organize devices into logical groups:

#### Creating Groups

1. Go to **Devices** → **Groups**
2. Click **"Create Group"**
3. Name the group and add description
4. Select devices to include
5. Save the group

#### Group Operations

- **Run Group Checks**: Execute security checks on all group devices
- **Group Reports**: Generate reports for specific groups
- **Group Settings**: Configure group-specific settings

## Connectivity Management

### Connection Testing

#### Manual Testing

- Test individual devices using the context menu
- Bulk test multiple devices simultaneously
- Schedule regular connectivity tests

#### Automatic Testing

- Enable automatic connectivity monitoring
- Configure test intervals (hourly, daily, weekly)
- Set up alerts for connectivity failures

#### Connection Troubleshooting

Common connectivity issues and solutions:

1. **SSH Connection Refused**

   - Verify SSH is enabled on device
   - Check firewall rules
   - Confirm correct port number

2. **Authentication Failed**

   - Verify username and password
   - Check for account lockouts
   - Ensure sufficient privileges

3. **Connection Timeout**

   - Increase timeout values
   - Check network connectivity
   - Verify device is powered on

4. **Host Key Verification Failed**
   - Update known hosts file
   - Disable host key verification (not recommended for production)
   - Manually verify host key fingerprint

### SSH Configuration

#### Per-Device SSH Settings

- **Custom SSH Port**: Override default port 22
- **Connection Timeout**: Device-specific timeout values
- **Authentication Method**: Password or key-based authentication
- **Keep-Alive Settings**: Maintain persistent connections

#### Global SSH Settings

Configure default SSH behavior:

- **Default Timeout**: 30 seconds
- **Max Retries**: 3 attempts
- **Retry Backoff**: Exponential delay between retries
- **Connection Pooling**: Reuse connections when possible

## Device Status Monitoring

### Status Types

#### Online Status

- Device responds to network connectivity tests
- SSH connection successful
- Ready for security checks

#### Offline Status

- Device not reachable via network
- SSH connection fails
- May indicate device or network issues

#### Unknown Status

- Device not yet tested
- Status check in progress
- Insufficient information to determine status

### Status History

Track device status over time:

- **Status Timeline**: Visual representation of status changes
- **Uptime Statistics**: Calculate device availability
- **Downtime Analysis**: Identify patterns in connectivity issues
- **Status Alerts**: Notifications for status changes

### Health Monitoring

Monitor device health indicators:

- **Response Time**: SSH connection establishment time
- **Check Success Rate**: Percentage of successful security checks
- **Error Frequency**: Rate of connection or check failures
- **Performance Trends**: Historical performance data

## Advanced Features

### Device Discovery

Automatically discover devices on your network:

1. **Network Scanning**

   - Configure IP ranges to scan
   - Detect SSH-enabled devices
   - Identify device types and vendors

2. **SNMP Discovery**
   - Use SNMP to gather device information
   - Automatically populate device details
   - Discover device relationships

### Device Templates

Create templates for common device configurations:

1. **Create Template**

   - Define common settings for device types
   - Include default credentials and settings
   - Specify security check preferences

2. **Apply Templates**
   - Use templates when adding new devices
   - Bulk apply templates to existing devices
   - Customize template settings per device

### Integration Features

#### Export Capabilities

- **CSV Export**: Export device lists for external tools
- **JSON Export**: Machine-readable device data
- **Report Integration**: Include device data in reports

#### API Access

- **REST API**: Programmatic access to device data
- **Webhooks**: Notifications for device events
- **Bulk Operations**: API endpoints for bulk device management

## Best Practices

### Device Management

- **Consistent Naming**: Use standardized device naming conventions
- **Regular Updates**: Keep device information current
- **Credential Management**: Use strong, unique passwords
- **Documentation**: Maintain device documentation and diagrams

### Security

- **Least Privilege**: Use accounts with minimal required permissions
- **Credential Rotation**: Regularly update device passwords
- **Access Logging**: Monitor device access patterns
- **Secure Storage**: Ensure device credentials are properly encrypted

### Organization

- **Logical Grouping**: Group devices by function, location, or criticality
- **Consistent Tagging**: Use standardized tag naming conventions
- **Regular Cleanup**: Remove obsolete devices and unused tags
- **Backup Data**: Regularly backup device configurations and data

### Performance

- **Batch Operations**: Use bulk operations for efficiency
- **Connection Limits**: Don't overwhelm devices with too many concurrent connections
- **Monitoring**: Monitor device performance and adjust settings as needed
- **Maintenance Windows**: Schedule intensive operations during maintenance windows

## Troubleshooting

### Common Issues

#### Device Not Appearing After Addition

- Check for duplicate IP addresses
- Verify all required fields are filled
- Review application logs for errors

#### Bulk Import Failures

- Validate CSV format and headers
- Check for special characters in device names
- Ensure IP addresses are unique

#### Connectivity Test Failures

- Verify network connectivity to device
- Check SSH service status on device
- Confirm credentials are correct

#### Performance Issues

- Reduce concurrent connection limits
- Increase timeout values
- Check network bandwidth and latency

### Getting Help

If you encounter issues:

1. Check the [Troubleshooting Guide](../admin/troubleshooting.md)
2. Review application logs
3. Test connectivity manually
4. Contact support with specific error messages

## Next Steps

After mastering device management:

1. Learn about [Security Checks](security-checks.md) to audit your devices
2. Explore the [Dashboard](dashboard.md) for monitoring overview
3. Set up [Reports](reports.md) for compliance documentation
4. Configure [Issue Management](issue-management.md) for tracking problems
