# Syntra 🚀

**Intelligent Synthetic Data Generation for API Testing**

Syntra is a comprehensive full-stack SaaS platform that automatically generates synthetic data and test cases for your backend API endpoints. It seamlessly integrates with your localhost environment, runs iterative tests, and provides detailed analytics on your API performance and reliability.

## ✨ Features

- **🎯 Smart Synthetic Data Generation**: AI-powered data generation tailored to your API schemas
- **⚡ Real-time API Testing**: Automated test case execution against your endpoints
- **🔄 Localhost Integration**: Direct connection to your development environment
- **📊 Comprehensive Analytics**: Detailed test results, performance metrics, and failure analysis
- **🎨 Modern UI/UX**: Clean, intuitive interface built with Next.js
- **⚙️ High Performance Backend**: Lightning-fast Go Fiber backend for optimal performance
- **📈 Test Iteration Management**: Track and compare test runs over time
- **🛡️ Data Privacy**: All testing happens locally - your data never leaves your environment

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   Your API      │
│   (Next.js)     │◄──►│  (Go Fiber)     │◄──►│  (Localhost)    │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Web Interface  │    │ Data Generator  │    │  Test Results   │
│  Test Config    │    │ Test Executor   │    │   Storage       │
│   Analytics     │    │   API Client    │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- **Node.js** (v18 or higher)
- **Go** (v1.19 or higher)
- **Docker** (optional, for containerized deployment)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/syntra.git
   cd syntra
   ```

2. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   cp .env.example .env.local
   # Configure your environment variables
   npm run dev
   ```

3. **Backend Setup**
   ```bash
   cd backend
   go mod tidy
   cp .env.example .env
   # Configure your environment variables
   go run main.go
   ```

4. **Access the Application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Docker Setup (Alternative)

```bash
# Build and run with Docker Compose
docker-compose up --build

# Or run individual containers
docker build -t syntra-frontend ./frontend
docker build -t syntra-backend ./backend
docker run -p 3000:3000 syntra-frontend
docker run -p 8080:8080 syntra-backend
```

## 📖 Usage

### 1. Configure Your API
- Navigate to the **Settings** page
- Add your API base URL (e.g., `http://localhost:4000`)
- Import your OpenAPI/Swagger specification or manually define endpoints

### 2. Generate Test Data
- Select the endpoints you want to test
- Configure data generation parameters:
  - Data types and constraints
  - Volume of test cases
  - Complexity levels
- Click **Generate Synthetic Data**

### 3. Run Tests
- Review generated test cases
- Customize test scenarios if needed
- Execute tests against your API
- Monitor real-time progress

### 4. Analyze Results
- View comprehensive test reports
- Identify failing endpoints and error patterns
- Export results for further analysis
- Track performance trends over time

## 🛠️ Configuration

### Frontend Configuration (.env.local)
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_APP_NAME=Syntra
NEXT_PUBLIC_VERSION=1.0.0
```

### Backend Configuration (.env)
```env
PORT=8080
ENVIRONMENT=development
DATABASE_URL=sqlite://syntra.db
CORS_ORIGINS=http://localhost:3000
LOG_LEVEL=info
```

## 🔧 API Reference

### Core Endpoints

#### Generate Test Data
```http
POST /api/v1/generate
Content-Type: application/json

{
  "endpoint": "/api/users",
  "method": "POST",
  "schema": {...},
  "count": 100,
  "complexity": "medium"
}
```

#### Execute Tests
```http
POST /api/v1/test/execute
Content-Type: application/json

{
  "target_url": "http://localhost:4000",
  "test_cases": [...],
  "config": {...}
}
```

#### Get Test Results
```http
GET /api/v1/test/results/{test_run_id}
```

For complete API documentation, visit `/docs` when running the backend server.

## 📊 Supported Data Types

- **Primitives**: String, Integer, Float, Boolean
- **Complex Types**: Objects, Arrays, Nested structures
- **Formats**: Email, Phone, URL, UUID, Date/Time
- **Custom Patterns**: Regex-based generation
- **Relationships**: Foreign keys, references
- **Constraints**: Min/Max values, length limits, enums

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- **Frontend**: Follow Next.js best practices, use TypeScript
- **Backend**: Follow Go conventions, write tests for new features
- **Testing**: Ensure all tests pass before submitting PR
- **Documentation**: Update README and API docs for new features

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support & Community

- **Documentation**: [docs.syntra.dev](https://docs.syntra.dev)
- **Issues**: [GitHub Issues](https://github.com/yourusername/syntra/issues)
- **Discord**: [Join our community](https://discord.gg/syntra)
- **Email**: support@syntra.dev

## 🗺️ Roadmap

- [ ] **v1.1**: GraphQL endpoint support
- [ ] **v1.2**: Advanced data relationships and constraints
- [ ] **v1.3**: CI/CD pipeline integration
- [ ] **v1.4**: Cloud deployment options
- [ ] **v1.5**: Performance benchmarking tools
- [ ] **v2.0**: Multi-language SDK support

## 📸 Screenshots

### Dashboard Overview
*Coming soon - Add screenshots of your main dashboard*

### Test Configuration
*Coming soon - Add screenshots of test setup interface*

### Results Analytics
*Coming soon - Add screenshots of analytics dashboard*

---

**Built with ❤️ by the SyntheticAPI Team**

*Making API testing simple, reliable, and intelligent.*