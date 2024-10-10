# /usr/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
    "name": "Travis",
    "age": 20,
    "country": "Spain",
    "degree": "Undergraduate",
    "status": "Offline",
    "site": "None"
}' http://localhost:8080/insert
