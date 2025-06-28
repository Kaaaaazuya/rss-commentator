# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

RSS Commentator is a serverless application that automatically fetches RSS feeds, summarizes technical articles using DeepSeek API, and notifies users. The system is built with Go Lambda functions, DynamoDB storage, and AWS CDK for infrastructure.

## Architecture

This is a serverless multi-service application with the following structure:

- **3 Lambda Functions** (Go): RSS Fetcher, Summarizer, Summary Notifier
- **Infrastructure**: AWS CDK v2 (TypeScript) in `infra/`
- **Database**: DynamoDB with tables for articles, tags, and article_tags
- **Shared Go Module**: Common models and utilities in `go/shared/`

### Key Directories
- `go/lambda/`: Contains 3 Lambda functions (rss-fetcher, summarizer, summary-notifier)
- `go/shared/`: Shared Go modules for models, database, and configuration
- `infra/`: AWS CDK infrastructure code
- `docker/`: Docker Compose files for local development
- `docs/`: Design documentation

## Common Development Commands

### Build and Deploy
```bash
# Build all Lambda functions
task build-all-lambdas

# Build all Docker images
task build-all-docker

# Deploy to AWS (with CDK)
task deploy

# Deploy with hotswap (faster Lambda updates)
task deploy:hotswap

# Deploy reusing existing DynamoDB tables
task deploy:reuse-tables
```

### Local Development
```bash
# Start local development environment (DynamoDB Local)
task dev:start

# Stop local development environment
task dev:stop

# Run Lambda functions locally
task lambda:local

# Invoke specific Lambda functions locally
task lambda:invoke:rss-fetcher
task lambda:invoke:summarizer
task lambda:invoke:notifier
```

### Infrastructure Management
```bash
# CDK commands (run from infra/ directory)
cd infra
npm run build
npm run synth    # Generate CloudFormation
npm run diff     # Show deployment differences
npm run deploy   # Deploy infrastructure
npm run destroy  # Destroy infrastructure

# Or use root-level task commands
task cdk:synth
task cdk:diff
task cdk:destroy
```

### ECR Management
```bash
# Create ECR repositories
task ecr:create:all

# Login to ECR
task ecr:login

# Push all images to ECR
task push-all
```

### Testing and Linting
```bash
# CDK infrastructure
cd infra
npm run test
npm run lint
npm run lint:fix
```

## Go Lambda Development

Each Lambda function in `go/lambda/` has its own:
- `Taskfile.yaml` for local tasks
- `Dockerfile` for containerization
- `go.mod` for dependencies
- `cmd/main.go` as entry point

### Shared Go Module
The `go/shared/` directory contains:
- `models/`: Data models (Article, Tag, ArticleTag)
- `db/`: DynamoDB client and operations
- `config/`: Configuration management
- `repos/`: Repository pattern implementations

## Database Schema

### DynamoDB Tables
- **articles**: Stores summarized articles (PK: id, attributes: url, title, summary, tags, created_at)
- **tags**: Tag information (PK: tag_name, attributes: created_at)
- **article_tags**: Many-to-many relationship with scoring (PK: tag_name, SK: article_id, attributes: score)

## Docker Development

### Local Environment
```bash
# Start all services (DynamoDB Local, Lambda containers)
docker compose -f docker/compose.yml up -d

# Lambda-specific development
docker compose -f docker/compose-lambda.yml up -d
```

### DynamoDB Local Management
```bash
# Create tables
docker/dynamodb/create-tables.sh

# Insert seed data
docker/dynamodb/insert-tags.sh

# Delete tables
docker/dynamodb/delete-tables.sh
```

## Task Runner Usage

This project uses [Task](https://taskfile.dev/) as the primary build tool. Key commands:

- `task --list`: Show all available tasks
- `task help`: Display help information
- Individual Lambda functions have their own Taskfile.yaml in their respective directories

## Environment Variables

Lambda functions expect standard AWS environment variables:
- `AWS_REGION`: Default ap-northeast-1
- `AWS_ACCOUNT_ID`: Retrieved automatically via AWS CLI
- DeepSeek API configuration (for summarizer)
- DynamoDB table names (configured via CDK)

## Development Notes

- All Lambda functions are containerized and use Docker for both local development and AWS deployment
- The CDK stack manages all AWS resources including DynamoDB tables, Lambda functions, and IAM roles
- Local development uses DynamoDB Local instead of actual AWS DynamoDB
- The project uses UUIDv7 for article IDs to maintain time-based ordering