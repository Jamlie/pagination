# /usr/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
    "pageSize": 10,
    "page": 1,
    "orderBy": {
        "age": "desc"
    }
}' http://localhost:8080/paginate
