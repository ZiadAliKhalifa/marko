# Marko API Collection for Bruno

This is a Bruno API collection for testing the Marko travel tracking application API.

## Setup

1. Open Bruno
2. Click "Open Collection"
3. Navigate to this folder (`/Users/ziadali/Desktop/gigs/marko/Marko/`)
4. Select the folder to open the collection

## Environment Variables

The collection uses environment variables that are defined in `environments/local.bru`:

- `baseUrl`: The base URL for the API (default: `http://localhost:8080`)
- `token`: The authentication token (default: `test-token`)
- `groupId`: The group ID for group-related requests (initially empty)

## Available Endpoints

1. **Health Check** - `GET /healthz`
2. **Create Group** - `POST /api/v1/groups`
3. **List User Groups** - `GET /api/v1/groups`
4. **Join Group** - `POST /api/v1/groups/:id/join`
5. **Get Group Members** - `GET /api/v1/groups/:id/members`
6. **Update Location** - `POST /api/v1/locations` (two variations: arrived/left)
7. **List Notifications** - `GET /api/v1/notifications`

## Usage Workflow

1. **Test the API health**: Run the "Health Check" request
2. **Create a group**: Run "Create Group" to create a new group
3. **Copy the group ID**: From the response, copy the `id` field
4. **Update the environment**: Set the `groupId` variable in the environment
5. **Test other endpoints**: Now you can test group-related endpoints

## Authentication

All API endpoints (except Health Check) require authentication using the Bearer token specified in the `token` environment variable.

## Notes

- The collection is configured for the development environment
- Make sure the Marko server is running on `http://localhost:8080`
- The test token `test-token` should work with the development server