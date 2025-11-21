#!/bin/sh
grpcurl -plaintext \
  -proto user.proto \
  -import-path ./proto/googleapis \
  -import-path ./pkg/api/user \
  -d '{
        "email": "test@example.com",
        "password": "StrongPassword123!"
      }' \
  127.0.0.1:50051 \
  user.UserService/Signup
