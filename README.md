curl -X POST http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '{
  "prompt": "benim siparişlerim",
  "user_id": 42
}'

curl -X POST http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '{                                                      
  "prompt": "2  tane Xiaomi Mi Band 8 sipariş ver",
  "user_id": 42
}'

curl -X POST http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '{"prompt": "son 3 Xiaomi Mi Band 8 siparişini göster"}'

curl -X POST http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '{"prompt": "son 5 siparişi göster"}'  

curl -X POST http://localhost:8080/api/chat \
-H "Content-Type: application/json" \
-d '{"prompt": "ürünlerin stok durumu nedir?"}' 

curl -X GET http://localhost:8080/api/orders

curl -X POST http://localhost:8080/api/products \
-H "Content-Type: application/json" \
-d '{                                          
  "name": "iPhone 14",     
  "description": "Apple akıllı telefon",      
  "price": 45999.99,
  "stock": 3
}'
