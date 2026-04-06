# Auth Service

A lightweight authentication and identity management service for the Flipopay backend. It handles tenant onboarding, project (client) registration, and user management.

## Overview

The service is built with **Go** and **Gin**, backed by **Supabase (PostgreSQL)**. It follows a multi-tenant architecture where each tenant can have multiple projects, and each project can have its own set of users.

```
Tenant → Projects (Clients) → Users
```

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **Database:** PostgreSQL (Supabase)
- **Auth:** bcrypt password hashing

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/tenant/create` | Register a new tenant |
| POST | `/api/v1/client/create` | Create a project under a tenant |
| POST | `/api/v1/user/create` | Create a user under a project |
| GET | `/health` | Health check with DB status |

## Getting Started

1. Clone the repo and set up the `.env` file:
   ```
   DATABASE_URL=postgresql://postgres:<password>@<host>:5432/postgres
   PORT=8000
   ```

2. Run the service:
   ```bash
   go run main.go
   ```

Tables are created automatically on startup.
