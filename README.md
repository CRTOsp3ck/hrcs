# HR Claims Management System

A comprehensive HR Claims Management application built with Go backend and Vue.js frontend, enabling employees to submit and track claims while providing administrators with powerful management and approval capabilities.

![HR Claims Management](https://img.shields.io/badge/Status-Production%20Ready-brightgreen)
![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![Vue Version](https://img.shields.io/badge/Vue-3.3+-green)
![License](https://img.shields.io/badge/License-MIT-yellow)

## ✨ Features

### 👥 **User Management**
- **Role-based Access Control**: Admin and Normal user roles
- **User Groups**: Organize employees into departments/teams
- **Profile Management**: User authentication and profile updates

### 💰 **Claims Management**
- **Multi-Status Workflow**: Draft → Submitted → Approved/Rejected → Payment → Paid
- **Claim Types**: Configurable categories (Travel, Medical, Office Supplies, etc.)
- **Rich Claims**: Title, description, amount, attachments support
- **Real-time Tracking**: Live status updates and notifications

### 🔄 **Approval Workflows**
- **Multi-level Approval**: Configure approval levels per user group
- **Flexible Permissions**: Define who can approve/reject at each level
- **Approval History**: Complete audit trail with comments
- **Automated Routing**: Claims automatically routed to appropriate approvers

### 📊 **Admin Dashboard**
- **Analytics Overview**: Claims statistics and metrics
- **User Management**: Promote/demote users, manage roles
- **System Configuration**: Claim types, user groups, approval workflows
- **Claims Oversight**: View and manage all claims in the system

### 🎨 **Modern UI/UX**
- **Responsive Design**: Works seamlessly on desktop and mobile
- **Professional Interface**: Clean, modern design
- **Intuitive Navigation**: Easy-to-use interface for all user types
- **Real-time Updates**: Dynamic status changes and notifications

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+**
- **Node.js 18+**
- **PostgreSQL 12+**
- **Docker & Docker Compose** (optional, for easy PostgreSQL setup)

### 1. Clone Repository
```bash
git clone <repository-url>
cd hrcs
```

### 2. Start Database
```bash
# Using Docker (recommended)
docker-compose up -d postgres

# Or configure your own PostgreSQL and update DATABASE_URL in .env
```

### 3. Setup Environment
```bash
# Copy environment file
cp .env.example .env

# Edit .env file with your database credentials if needed
```

### 4. Run Automated Setup
```bash
# Option 1: Use setup script
./setup.sh

# Option 2: Use Makefile
make setup
```

### 5. Seed Database
```bash
# Seed with sample data
make seed

# Or clear and reseed
make seed-clear
```

### 6. Start Application
```bash
# Start both backend and frontend
make dev

# Or start individually
make dev-backend  # Backend: http://localhost:8000
make dev-frontend # Frontend: http://localhost:3000
```

## 🔑 Default Login Credentials

After seeding, you can log in with these accounts:

### Admin Users
| Email | Password | Role |
|-------|----------|------|
| `admin@hrcs.com` | `password123` | Super Admin |
| `hr.manager@hrcs.com` | `password123` | HR Manager |
| `finance.manager@hrcs.com` | `password123` | Finance Manager |
| `dept.head@hrcs.com` | `password123` | Department Head |

### Normal Users
| Email | Password | Department |
|-------|----------|------------|
| `john.doe@hrcs.com` | `password123` | Engineering |
| `jane.smith@hrcs.com` | `password123` | Sales |
| `bob.wilson@hrcs.com` | `password123` | Engineering |
| `alice.brown@hrcs.com` | `password123` | Marketing |
| (and 4 more users...) | `password123` | Various |

## 📁 Project Structure

```
hrcs/
├── backend/                 # Go backend application
│   ├── cmd/                # Command-line tools
│   │   └── seed/           # Database seeding tool
│   ├── config/             # Configuration management
│   ├── database/           # Database connection and migrations
│   ├── handlers/           # HTTP request handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/             # Database models
│   ├── routes/             # API route definitions
│   ├── seeds/              # Database seeding logic
│   └── utils/              # Utility functions
├── frontend/               # Vue.js frontend application
│   ├── src/
│   │   ├── api/           # API service layer
│   │   ├── components/    # Reusable Vue components
│   │   ├── router/        # Vue Router configuration
│   │   ├── stores/        # Pinia stores (state management)
│   │   ├── types/         # TypeScript type definitions
│   │   └── views/         # Page components
│   └── public/            # Static assets
├── scripts/               # Utility scripts
├── docker-compose.yml     # Docker services configuration
├── Makefile              # Build and development commands
└── README.md             # This file
```

## 🛠️ Development

### Backend Development
```bash
cd backend

# Install dependencies
go mod download

# Run with hot reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Or run normally
go run main.go

# Run tests
go test ./...

# Build
go build -o bin/hrcs-backend main.go
```

### Frontend Development
```bash
cd frontend

# Install dependencies
npm install

# Start dev server with hot reload
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

### Database Operations
```bash
# Seed database
make seed

# Clear and reseed
make seed-clear

# Manual seeding
cd backend && go run cmd/seed/main.go

# Clear data before seeding
cd backend && go run cmd/seed/main.go -clear
```

## 📊 API Documentation

### Authentication Endpoints
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `GET /api/profile` - Get user profile

### Claims Endpoints
- `GET /api/claims` - List claims (personal for normal users, all for admins)
- `POST /api/claims` - Create new claim
- `GET /api/claims/{id}` - Get claim details
- `PUT /api/claims/{id}` - Update claim
- `DELETE /api/claims/{id}` - Cancel claim
- `POST /api/claims/{id}/submit` - Submit claim for approval
- `POST /api/claims/{id}/approve` - Approve/reject claim (admin only)

### Admin Endpoints
- `GET /api/users` - List all users (admin only)
- `PUT /api/users/{id}/role` - Update user role (admin only)
- `GET /api/claim-types` - List claim types
- `POST /api/claim-types` - Create claim type (admin only)
- `PUT /api/claim-types/{id}` - Update claim type (admin only)
- `DELETE /api/claim-types/{id}` - Delete claim type (admin only)
- `GET /api/user-groups` - List user groups
- `POST /api/user-groups` - Create user group (admin only)
- `PUT /api/user-groups/{id}` - Update user group (admin only)
- `DELETE /api/user-groups/{id}` - Delete user group (admin only)
- `GET /api/approval-levels` - List approval levels
- `POST /api/approval-levels` - Create approval level (admin only)
- `DELETE /api/approval-levels/{id}` - Delete approval level (admin only)

## 🔧 Configuration

### Environment Variables
Create a `.env` file in the project root:

```env
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost/hrcs?sslmode=disable

# JWT Secret (change in production!)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Server Port
PORT=8000
```

### Database Configuration
The application uses PostgreSQL. Configure your database connection in the `.env` file or use the provided Docker Compose setup.

## 🚀 Deployment

### Production Build
```bash
# Build backend
make build

# The binary will be created at: bin/hrcs-backend
```

### Docker Deployment
```bash
# Build and start all services
docker-compose up -d

# Scale services if needed
docker-compose up -d --scale backend=2
```

### Environment Setup
1. Set production environment variables
2. Use a strong JWT secret
3. Configure production database
4. Enable HTTPS in production
5. Set up proper logging and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Troubleshooting

### Common Issues

**Database Connection Error**
```bash
# Make sure PostgreSQL is running
docker-compose up -d postgres

# Check connection
psql -h localhost -U postgres -d hrcs
```

**Frontend Build Errors**
```bash
# Clear node modules and reinstall
cd frontend
rm -rf node_modules package-lock.json
npm install
```

**Backend Module Issues**
```bash
# Clean and download modules
cd backend
go clean -modcache
go mod download
```

### Getting Help
- Check the [Issues](../../issues) page for known problems
- Create a new issue if you encounter a bug
- Check logs in `backend/logs/` for detailed error information

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/) and [Vue.js](https://vuejs.org/)
- UI inspired by modern design principles
- Database powered by [PostgreSQL](https://www.postgresql.org/)
- Authentication using JWT tokens
