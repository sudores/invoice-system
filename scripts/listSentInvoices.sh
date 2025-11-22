#!/bin/sh
. ./scripts/loginUser.sh
grpcurl -plaintext \
  -H "Authorization: Bearer ${JWT_TOKEN}" \
  -proto ./pkg/api/invoice/invoice.proto \
  -import-path ./proto/googleapis \
  -import-path ./pkg/api/invoice \
  127.0.0.1:50051 \
  invoice.InvoiceService/ListSentInvoices

