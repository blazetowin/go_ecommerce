 curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 14", "description":"Applenin son modeli", "price":39999.99, "in_stock":true}'

curl http://localhost:8080/api/products

curl -X POST http://localhost:8080/api/chat \
  -H "Content-Type: application/json" \
  -d '{"prompt": "iphone 14 satın almak istiyorum"}'

curl http://localhost:8080/api/orders 

sqlite3 ecommerce.db "SELECT name, stock FROM products WHERE name = 'iPhone 14';"


