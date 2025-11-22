#!/bin/sh
grpcurl -plaintext \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjM3MjYxNzYsImlhdCI6MTc2MzcyMjU3Niwic3ViIjoiYWZlYjUzMTAtMmIyMy00MDRlLWIzMzQtZjZkYTVhZTc1ZmZiIn0.oJvfZLEZPGjNuqJvlo9yWy0cWpJ34LE1jXUeDXYMGTo" \
  -proto ./pkg/api/invoice/invoice.proto \
  -import-path ./proto/googleapis \
  -import-path ./pkg/api/invoice \
  127.0.0.1:50051 \
  invoice.InvoiceService/ListReceivedInvoices

