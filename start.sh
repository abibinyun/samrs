#!/bin/bash

echo "🚀 Starting SAMRS Project..."
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env created. Please edit it if needed."
    echo ""
fi

# Start services
echo "🐳 Starting Docker services..."
docker compose up -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check services
echo ""
echo "🔍 Checking services status..."
docker compose ps

echo ""
echo "✅ SAMRS Project is running!"
echo ""
echo "📍 Access URLs:"
echo "   Frontend:  http://localhost:8000"
echo "   Backend:   http://localhost:8090"
echo "   BFF:       http://localhost:3000"
echo ""
echo "👤 Default Login:"
echo "   Username: admin"
echo "   Password: password123"
echo "   Tenant:   rs-pusat"
echo ""
echo "📚 Run 'make help' for more commands"
