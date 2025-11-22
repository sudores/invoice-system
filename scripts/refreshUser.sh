#!/bin/sh
data=$(grpcurl -plaintext \
  -proto user.proto \
  -import-path ./proto/googleapis \
  -import-path ./pkg/api/user \
  -d '{
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjY0NDQ0MzEsImlhdCI6MTc2Mzg1MjQzMSwic3ViIjoiZDY2MzZjZGMtYmEzMy00NzkyLTkwMTMtYjU4Nzg3MjgyOWRhIn0.3Gju31OCCAxM7fmxXbj5iZASoilrfDhblGcjkR0sdp4",
      }' \
  127.0.0.1:50051 \
  user.UserService/Refresh)

echo $data
