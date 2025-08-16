# Quick Start Guide

Get up and running with the Network Configuration Checker in just 5 minutes! This guide will walk you through the essential steps to start monitoring your network devices.

## Step 1: Launch the Application

After installation, launch the Network Configuration Checker:

- **Windows**: Start Menu → Network Configuration Checker
- **macOS**: Applications → Network Configuration Checker
- **Linux**: Run `network-config-checker` or find it in your applications menu

## Step 2: Initial Setup (2 minutes)

### Set Application Password (Recommended)

1. Click the **Settings** icon in the top navigation
2. Go to **Security** tab
3. Enable **"Password Protection"**
4. Set a strong password and confirm
5. Click **Save**

### Configure Basic Settings

1. In Settings, go to **Connections** tab
2. Set **SSH Timeout** to `30` seconds (adjust for your network)
3. Set **Max Concurrent Connections** to `5` (start conservative)
4. Click **Save**

## Step 3: Add Your First Device (2 minutes)

1. Navigate to the **Devices** page
2. Click the **"Add Device"** button
3. Fill in the device information:

   ```
   Device Name: Core-Router-01
   IP Address: 192.168.1.1
   Device Type: Router
   Vendor: Cisco
   Username: admin
   Password: [your device password]
   SSH Port: 22
   ```

4. Click **"Test Connectivity"** to verify the connection
5. If successful, click **"Save Device"**

Your device will appear in the device grid with a status indicator.

## Step 4: Run Your First Security Check (1 minute)

1. Find your newly added device in the device grid
2. Click the **three-dot menu** on the device card
3. Select **"Run Security Check"**
4. Watch the progress indicator as the check runs
5. View results when complete

## Step 5: Explore the Dashboard

1. Navigate to the **Dashboard** page
2. You'll see:
   - **Device Overview**: Total devices and their status
   - **Security Status**: Pie chart of check results
   - **Device Grid**: Visual representation of all devices
   - **Recent Activity**: Latest security check results

## What's Next?

Now that you have the basics working, explore these features:

### Add More Devices

- **Bulk Import**: Use CSV import for multiple devices
- **Device Groups**: Organize devices with tags
- **Bulk Operations**: Run checks on multiple devices simultaneously

### Explore Security Features

- **Issue Management**: View and track security issues
- **Custom Rules**: Create custom security checks
- **Scheduled Checks**: Automate regular security audits

### Generate Reports

- **Executive Summary**: High-level compliance overview
- **Technical Reports**: Detailed findings with remediation steps
- **Scheduled Reports**: Automatic report generation and email delivery

## Common First-Time Tasks

### Adding Multiple Devices

If you have many devices to add:

1. Go to **Devices** → **Import**
2. Download the CSV template
3. Fill in your device information
4. Upload the completed CSV file

### Setting Up Scheduled Checks

1. Go to **Settings** → **Scheduling**
2. Enable **"Automatic Security Checks"**
3. Set frequency (daily, weekly, monthly)
4. Choose which devices to include
5. Set notification preferences

### Configuring Email Reports

1. Go to **Settings** → **Email**
2. Configure SMTP settings:
   ```
   SMTP Server: smtp.yourcompany.com
   Port: 587
   Username: reports@yourcompany.com
   Password: [email password]
   ```
3. Test email configuration
4. Set up report schedules in **Reports** → **Schedule**

## Troubleshooting Quick Issues

### Device Won't Connect

1. **Check Network Connectivity**: Ping the device IP
2. **Verify Credentials**: Ensure username/password are correct
3. **Check SSH Access**: Verify SSH is enabled on the device
4. **Firewall Rules**: Ensure port 22 is accessible

### Security Checks Fail

1. **Check Device Type**: Ensure correct vendor is selected
2. **Verify Permissions**: User needs sufficient privileges
3. **Review Logs**: Check application logs for detailed errors

### Application Runs Slowly

1. **Reduce Concurrency**: Lower max concurrent connections
2. **Check Resources**: Monitor CPU and memory usage
3. **Network Latency**: Increase timeout values for slow networks

## Getting Help

### Built-in Help

- **Tooltips**: Hover over UI elements for quick help
- **Form Validation**: Real-time feedback on form inputs
- **Status Indicators**: Color-coded status throughout the app

### Documentation

- **User Guides**: Detailed guides for each feature
- **Admin Guides**: Advanced configuration and troubleshooting
- **API Reference**: For integration and automation

### Support Channels

- **GitHub Issues**: Report bugs and request features
- **Community Forum**: Ask questions and share tips
- **Email Support**: Direct support for critical issues

## Sample Workflow

Here's a typical workflow after initial setup:

### Daily Monitoring

1. **Check Dashboard**: Review overall status
2. **Review Alerts**: Address any critical issues
3. **Monitor Trends**: Look for patterns in security status

### Weekly Tasks

1. **Run Bulk Checks**: Execute security checks on all devices
2. **Review Issues**: Triage and assign security issues
3. **Generate Reports**: Create weekly compliance reports

### Monthly Activities

1. **Add New Devices**: Update inventory as network grows
2. **Review Rules**: Update security check rules as needed
3. **Archive Data**: Clean up old results and reports

## Pro Tips

### Efficiency Tips

- **Use Tags**: Organize devices by location, criticality, or team
- **Keyboard Shortcuts**: Learn shortcuts for common actions
- **Bulk Operations**: Select multiple devices for batch operations
- **Filters**: Use search and filters to quickly find devices

### Security Best Practices

- **Regular Updates**: Keep the application updated
- **Strong Passwords**: Use unique, strong passwords for each device
- **Access Control**: Limit who has access to the application
- **Audit Logs**: Regularly review audit logs for suspicious activity

### Performance Optimization

- **Stagger Checks**: Don't run all checks simultaneously
- **Monitor Resources**: Keep an eye on system resource usage
- **Clean Up Data**: Regularly clean up old results and logs
- **Network Optimization**: Use wired connections when possible

## Next Steps

Ready to dive deeper? Check out these guides:

1. **[Device Management](user-guide/device-management.md)**: Advanced device management features
2. **[Security Checks](user-guide/security-checks.md)**: Comprehensive security auditing
3. **[Reports](user-guide/reports.md)**: Advanced reporting and analytics
4. **[Configuration](configuration.md)**: Fine-tune the application for your environment

Welcome to the Network Configuration Checker! You're now ready to start automating your network security audits.
