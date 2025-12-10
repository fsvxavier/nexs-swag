# Get User Endpoint

Retrieves user information by ID.

## Path Parameters

- `id` (integer, required): The user ID

## Response

Returns a user object containing:
- `id`: User ID
- `name`: User's full name

## Caching

This endpoint implements caching:
- Cache TTL: 5 minutes
- Cache-Control header is set appropriately
- Use `If-None-Match` for conditional requests

## Performance

- Average response time: <50ms
- Supports pagination for large datasets
- Implements database indexing

## Security

- Requires authentication
- Users can only access their own data unless admin
- Sensitive fields are filtered based on permissions
