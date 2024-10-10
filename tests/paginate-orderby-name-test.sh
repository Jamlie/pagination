# /usr/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
    "pageSize": 15,
    "page": 1,
    "orderBy": {
        "name": "asc"
    }
}' http://localhost:8080/paginate
