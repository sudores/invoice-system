# Invoice system

__Note:__ name should be chosen as there's no adequate name so far

This should be simple invoice system to create invoices in a small group of
people like family, friends, etc.
It should allow to create invoice with due date, status, notes and amount.

## Requirements

```
sudo xbps-install -Sy protoc
go install 'google.golang.org/protobuf/cmd/protoc-gen-go@latest'
go install 'google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest'
go install 'github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest'
go install 'github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest'
```

