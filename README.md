# Network Monitor Backend

A Go-based backend service for real-time network device monitoring and status tracking with email notifications.

## 🚀 Features

- **Real-time Ping Monitoring**: Continuously monitors network devices every 10 seconds
- **Device Management**: Add, remove, and track network devices with custom names and locations
- **Email Notifications**: Automatic email alerts when device status changes (online/offline)
- **Redis Storage**: Fast in-memory database for device status and logs
- **RESTful API**: Clean HTTP endpoints for frontend integration
- **Comprehensive Logging**: Device-specific and global status change logs

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │───▶│   Backend API   │───▶│   Redis DB      │
│   (Next.js)     │    │   (Go Server)   │    │   (Storage)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │  Email Service  │
                       │   (Resend API)  │
                       └─────────────────┘
```

## 📁 Project Structure

```
backend/
├── main.go           # Server entry point & routing
├── .env             # Environment variables
├── go.mod           # Go dependencies
├── api/
│   └── handlers.go  # HTTP request handlers
├── config/
│   └── ip.go        # IP validation utilities
├── db/
│   └── redis.go     # Redis connection & setup
├── mail/
│   └── mail.go      # Email notification service
├── models/
│   └── device.go    # Device data structures
├── ping/
│   └── ping.go      # Network ping functionality
└── store/
    └── device.go    # Device storage operations
```

## 🔧 Setup & Installation

### Prerequisites
- Go 1.23+
- Redis server
- Resend API account (for email notifications)

### 1. Clone & Install
```bash
git clone <repository-url>
cd backend
go mod tidy
```

### 2. Environment Configuration
Create a `.env` file:
```env
RESEND_API_KEY=your_resend_api_key_here
NOTIFICATION_EMAILS=admin@company.com,alerts@company.com
```

### 3. Start Redis
```bash
# On Windows (if Redis is installed)
redis-server

# On Docker
docker run -d -p 6379:6379 redis:alpine
```

### 4. Run the Server
```bash
go run .
# OR build and run
go build
./backend.exe
```

Server starts on `http://localhost:8080`

## 📡 API Endpoints

### Device Management
- `POST /add?ip={ip}&location={location}&name={name}` - Add new device
- `POST /delete?ip={ip}` - Remove device
- `GET /status` - Get all device statuses

### Logging
- `GET /logs` - Get global status change logs
- `GET /device-logs?ip={ip}` - Get device-specific logs

## 💾 Data Storage

### Redis Keys
- `device:{ip}:online` - Device online status (true/false)
- `device:{ip}:location` - Device location/area
- `device:{ip}:name` - Device display name
- `device:{ip}:lastseen` - Last seen timestamp
- `global_logs` - Global status change logs
- `device:{ip}:logs` - Device-specific logs

### Device Areas
Devices are organized by location areas:
- PGCIL, Sophos, Hop Bung, SSC Build
- Plant Area, IT Dept, Admin Build
- Sankalp #2, Township, ET-Hostel, RLI Office

## 📧 Email Notifications

Automatic email alerts are sent when:
- Device goes **offline** (unreachable)
- Device comes **online** (responds to ping)

### Email Configuration
- **Service**: Resend API
- **Recipients**: Configured via `NOTIFICATION_EMAILS` environment variable
- **Content**: Device IP, name, location, and status change details

## 🔄 Monitoring Process

1. **Periodic Ping**: Every 10 seconds, all registered devices are pinged
2. **Status Detection**: Compare current status with previous status
3. **Change Logging**: Log any status changes with timestamp
4. **Email Alerts**: Send notifications for status changes
5. **Data Persistence**: Update Redis with new status and logs

## 🛠️ Key Components

### Ping Service (`ping/ping.go`)
- Executes system ping commands
- Parses ping responses for reachability
- Handles cross-platform ping differences

### Storage Layer (`store/device.go`)
- Redis operations for device data
- Status change detection and logging
- Batch operations for efficiency

### Email Service (`mail/mail.go`)
- Network connectivity checks
- Resend API integration
- Graceful error handling

### API Handlers (`api/handlers.go`)
- HTTP request processing
- JSON response formatting
- Error handling and validation

## 🚦 Status Codes

### Device Status
- `"true"` - Device is online (ping successful)
- `"false"` - Device is offline (ping failed/timeout)

### HTTP Responses
- `200` - Success
- `400` - Bad Request (missing parameters)
- `500` - Internal Server Error

## 📊 Logging Format

### Global Logs
```json
{
  "timestamp": "2025-07-04T16:30:00+05:30",
  "device_ip": "192.168.1.1",
  "device_name": "Router",
  "location": "IT Dept",
  "status": "offline",
  "message": "Device Router (192.168.1.1) in IT Dept went offline"
}
```

### Device-Specific Logs
Same format but filtered by device IP.

## 🔒 Security Considerations

- Environment variables for sensitive data
- Input validation for IP addresses
- Rate limiting considerations for ping operations
- Network timeout configurations

## 🐛 Troubleshooting

### Common Issues

**Redis Connection Failed**
```bash
# Check if Redis is running
redis-cli ping
# Should return: PONG
```

**Email Notifications Not Working**
- Verify `RESEND_API_KEY` in `.env`
- Check `NOTIFICATION_EMAILS` format
- Ensure network connectivity to `api.resend.com`

**Ping Not Working**
- Verify devices are actually reachable
- Check firewall/network settings
- Ensure ping command is available on system

## 📈 Performance

- **Ping Interval**: 10 seconds (configurable)
- **Concurrent Pings**: Handled via goroutines
- **Redis Operations**: Optimized for speed
- **Memory Usage**: Minimal overhead with Redis

## 🔄 Future Enhancements

- Custom ping intervals per device
- SNMP monitoring support
- Web dashboard authentication
- Database migration from Redis to PostgreSQL
- Docker containerization
- Kubernetes deployment support

---

**Developed for NTPC Network Infrastructure Monitoring**