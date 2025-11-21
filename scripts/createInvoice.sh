#!/bin/sh
grpcurl -plaintext \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjM3MjYxNzYsImlhdCI6MTc2MzcyMjU3Niwic3ViIjoiYWZlYjUzMTAtMmIyMy00MDRlLWIzMzQtZjZkYTVhZTc1ZmZiIn0.oJvfZLEZPGjNuqJvlo9yWy0cWpJ34LE1jXUeDXYMGTo" \
  -proto ./pkg/api/invoice/invoice.proto \
  -d '{
        "sender_id": "5e441075-6119-4286-9dfe-2785b524e1bc",
        "recipient_id": "071b48d4-9115-401f-a3bd-4c0439d1900e",
        "items": [
          {
            "title": "Service A",
            "amount": 2,
            "reference_id": "46446244-7972-4d8e-ba05-e89c1a9140c3"
          },
          {
            "title": "Service B",
            "amount": 1,
            "reference_id": "c3725b86-bb7f-4f31-af80-3d2ad96754af"
          }
        ],
        "description": "Monthly services",
        "due_date": "2025-12-20T00:00:00Z"
      }' \
  127.0.0.1:50051 \
  invoice.InvoiceService/CreateInvoice

