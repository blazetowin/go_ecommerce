 curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 14", "description":"Applenin son modeli", "price":39999.99, "in_stock":true}'

curl http://localhost:8080/api/products

