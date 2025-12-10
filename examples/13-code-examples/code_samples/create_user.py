# Python example for creating a user
import requests
import json

url = "http://localhost:8080/api/users"
payload = {
    "name": "John Doe"
}
headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, data=json.dumps(payload), headers=headers)
print(response.json())
