# /usr/bin/bash

curl -X POST -H "Content-Type: application/json" -d '{
    "pageSize": 5,
    "page": 1,
    "filters": {
        "countries": ["Ireland", "Spain"]
    }
}' http://localhost:8080/paginate
