#!/bin/bash

set -e

echo "🚀 Setting up Go Modular Monolith..."

# Check if required tools are installed
check_tool() {
    if ! command -v $1 &> /dev/null; then
        echo "❌ $1 is not installed. Please install it first."
        exit 1
    fi
}

echo "📋 Checking required tools..."
check_tool go
check_tool docker
check_tool task

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your configuration."
fi

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy

# Build the application
echo "🔨 Building application..."
task app:build

# Start database services
echo "🗄️ Starting database services..."
task db:start

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 10

# Run migrations
echo "🔧 Running database migrations..."
task db:migrate-up

echo ""
echo "🎉 Setup completed successfully!"
echo ""
echo "📖 Next steps:"
echo "  1. Update your .env file with proper configuration"
echo "  2. Run 'task app:run' to start the application"
echo "  3. Visit http://localhost:8080/health to check if the app is running"
echo ""
echo "📚 Useful commands:"
echo "  task                    # Show all available tasks"
echo "  task app:run           # Run the application"
echo "  task db:migrate-up     # Run database migrations"
echo "  task test:unit         # Run unit tests"
echo "  task docker:run        # Run with Docker Compose"
echo ""
echo "🔗 Documentation:"
echo "  Health Check: http://localhost:8080/health"
echo "  API Base URL: http://localhost:8080/api/v1"
echo "  Jaeger UI: http://localhost:16686"
echo ""