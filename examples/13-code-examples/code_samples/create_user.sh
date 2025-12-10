#!/bin/bash
# Bash example for creating a user

curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe"
  }'
