# Create User Endpoint

Creates a new user in the system.

## Request Body

The request body should contain user information:
- `name`: User's full name (required)
- `id`: Will be auto-generated if not provided

## Validation Rules

- Name must be at least 3 characters
- Name cannot contain special characters
- Duplicate names are not allowed

## Example Request

```json
{
  "name": "John Doe"
}
```

## Example Response

```json
{
  "id": 123,
  "name": "John Doe"
}
```

## Error Responses

- `400 Bad Request`: Invalid input data
- `409 Conflict`: User already exists
- `500 Internal Server Error`: Server error

## Rate Limiting

This endpoint is rate limited to:
- 10 requests per minute for unauthenticated users
- 100 requests per minute for authenticated users
