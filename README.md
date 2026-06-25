# Auth Service

A lightweight authentication and identity management service for the opensource backend. It supports tenant onboarding, user pool creation, app client registration, and authenticated user management.

## What this project does

This service provides:

- Tenant creation and management
- Multiple user pools per tenant
- App client (project) registration with generated client ID and secret
- Authenticated user operations via app client credentials
- Email/password and OTP login flows
- Health checks for service and database connectivity

The architecture is:

```
Tenant → User Pools → App Clients → Users
```

A tenant can create one or more user pools and one or more app clients. Each app client receives a `clientId` and `clientSecret` and uses those credentials to access authenticated user endpoints.

## How to install and run it

1. Clone the repository:
   ```bash
   git clone https://github.com/raam-techie/new-auth-service.git
   cd new-auth-service
   ```

2. Configure your environment variables. Create a `.env` file or export the values:
   ```bash
   DATABASE_URL=postgresql://postgres:<password>@<host>:5432/postgres
   PORT=8000
   ```

3. Run the service:
   ```bash
   go run main.go
   ```

4. Verify the service is running:
   ```bash
   curl http://localhost:8000/health
   ```

> The database tables are created automatically on startup if they do not already exist.

## How to use it

Available API endpoints:

- `POST /api/v1/tenant/create` — Register a new tenant
- `POST /api/v1/user-pool/create` — Create a new user pool for a tenant
- `GET /api/v1/user-pool/fetch` — Fetch a user pool or all pools for a tenant
- `POST /api/v1/app-client/create?tenantId=<tenant-id>` — Create an app client for a tenant
- `GET /api/v1/app-client/fetch?tenantId=<tenant-id>` — Get app clients for a tenant
- `POST /api/v1/user/create` — Create a user with app client credentials
- `POST /api/v1/user/emailpassword` — Login a user with email and password
- `POST /api/v1/user/otp/email/send` — Send an OTP email
- `POST /api/v1/user/otp/email/verify` — Verify OTP and receive tokens
- `PUT /api/v1/user/:userId/attributes` — Update a user’s attributes
- `GET /api/v1/user/fetch` — Fetch users by ID or tenant

### Authenticated requests

The `POST /api/v1/user/*` endpoints require app client authentication with headers:

- `X-Client-ID: <client-id>`
- `X-Client-Secret: <client-secret>`

### Example usage

Create a tenant:

```bash
curl -X POST http://localhost:8000/api/v1/tenant/create \
  -H "Content-Type: application/json" \
  -d '{"name":"example-tenant","email":"owner@example.com"}'
```

Create a user pool:

```bash
curl -X POST http://localhost:8000/api/v1/user-pool/create \
  -H "Content-Type: application/json" \
  -d '{"tenantId":"<tenant-id>","poolName":"example-pool","signInMethods":["email"],"appType":"web","description":"Example pool"}'
```

Create an app client for the tenant:

```bash
curl -X POST "http://localhost:8000/api/v1/app-client/create?tenantId=<tenant-id>" \
  -H "Content-Type: application/json" \
  -d '{"clientName":"example-client","appType":"web"}'
```

Create a user using the app client credentials:

```bash
curl -X POST http://localhost:8000/api/v1/user/create \
  -H "Content-Type: application/json" \
  -H "X-Client-ID: <client-id>" \
  -H "X-Client-Secret: <client-secret>" \
  -d '{"name":"Example User","email":"user@example.com","password":"password123","role":"user"}'
```

## How to contribute

1. Fork the repository.
2. Create a feature branch: `git checkout -b feature/my-change`
3. Make your changes and add tests if needed.
4. Commit your work: `git commit -m "Add feature description"`
5. Push to your fork and open a pull request.

Contributions are welcome for:

- bug fixes
- API improvements
- security enhancements
- documentation updates

## Tech Stack

- Language: Go
- Framework: Gin
- Database: PostgreSQL (Supabase)
- Password hashing: bcrypt
