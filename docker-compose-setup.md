# Docker Compose Setup Guide

This guide will help you set up and run the Node.js application using Docker Compose.

## Prerequisites

- Docker installed on your system
- Docker Compose installed (usually comes with Docker Desktop)
- Git (to clone the repository)

## Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/wearedood/node.git
   cd node
   ```

2. Build and start the services:
   ```bash
   docker-compose up --build
   ```

3. The application will be available at `http://localhost:3000`

## Configuration

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
NODE_ENV=development
PORT=3000
DATABASE_URL=mongodb://mongo:27017/nodeapp
REDIS_URL=redis://redis:6379
```

### Docker Compose Services

The `docker-compose.yml` file defines the following services:

#### Web Application
- **Service**: `app`
- **Port**: 3000
- **Dependencies**: MongoDB, Redis
- **Volume**: Source code mounted for development

#### Database
- **Service**: `mongo`
- **Port**: 27017
- **Volume**: Persistent data storage

#### Cache
- **Service**: `redis`
- **Port**: 6379
- **Volume**: Persistent cache storage

## Development Workflow

### Starting Services
```bash
# Start all services
docker-compose up

# Start services in background
docker-compose up -d

# Start specific service
docker-compose up app
```

### Stopping Services
```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Viewing Logs
```bash
# View logs for all services
docker-compose logs

# View logs for specific service
docker-compose logs app

# Follow logs in real-time
docker-compose logs -f app
```

### Running Commands
```bash
# Execute command in running container
docker-compose exec app npm install

# Run one-time command
docker-compose run app npm test
```

## Production Deployment

For production deployment, use the production compose file:

```bash
docker-compose -f docker-compose.prod.yml up -d
```

### Production Considerations

1. **Environment Variables**: Set production values in `.env.production`
2. **Volumes**: Use named volumes for data persistence
3. **Networks**: Configure proper network isolation
4. **Health Checks**: Ensure services have health check endpoints
5. **Resource Limits**: Set appropriate CPU and memory limits

## Troubleshooting

### Common Issues

1. **Port Already in Use**
   ```bash
   # Check what's using the port
   lsof -i :3000
   
   # Kill the process or change the port in docker-compose.yml
   ```

2. **Database Connection Issues**
   ```bash
   # Check if MongoDB is running
   docker-compose ps mongo
   
   # View MongoDB logs
   docker-compose logs mongo
   ```

3. **Permission Issues**
   ```bash
   # Fix file permissions
   sudo chown -R $USER:$USER .
   ```

### Rebuilding Services

If you make changes to the Dockerfile or dependencies:

```bash
# Rebuild specific service
docker-compose build app

# Rebuild all services
docker-compose build

# Force rebuild without cache
docker-compose build --no-cache
```

## Useful Commands

```bash
# View running containers
docker-compose ps

# View resource usage
docker-compose top

# Remove unused images and containers
docker system prune

# View service configuration
docker-compose config
```

## File Structure

```
.
├── docker-compose.yml          # Main compose file
├── docker-compose.prod.yml     # Production compose file
├── Dockerfile                  # Application container definition
├── .env                        # Environment variables
├── .env.production            # Production environment variables
└── docs/
    └── docker-compose-setup.md # This documentation
```

## Next Steps

- Review the `docker-compose.yml` file to understand the service configuration
- Customize environment variables for your specific needs
- Set up CI/CD pipeline for automated deployments
- Configure monitoring and logging for production environments

For more information, refer to the [Docker Compose documentation](https://docs.docker.com/compose/).
