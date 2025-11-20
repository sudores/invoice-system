#!/bin/sh
grpcurl -plaintext \
  -proto ./pkg/api/user/user.proto \
  -d '{
        "email": "test@example.com",
        "password": "StrongPassword123!"
      }' \
  127.0.0.1:50051 \
  user.UserService/RegisterUser
