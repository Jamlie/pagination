# /usr/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
    "pageSize": 5,
    "page": 1,
    "orderBy": {
        "id": "desc"
    }
}' http://localhost:8080/paginate
