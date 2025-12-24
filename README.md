docker run -d \
  -p 8080:8080 \
  -e DB_HOST=localhost \
  -e DB_PORT=5432 \
  -e DB_USER=alie \
  -e DB_PASSWORD=12345678 \
  -e DB_NAME=pastebin_db \
  -e SERVER_PORT=8080 \
  pastebin-backend
