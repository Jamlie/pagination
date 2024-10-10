# /usr/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
    "pageSize": 4,
    "page": 3
}' http://localhost:8080/paginate

