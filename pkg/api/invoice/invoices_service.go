package invoice

import (
	context "context"

	"github.com/rs/zerolog"
)

type InvoicesGrpcService struct {
	UnimplementedInvoiceServiceServer
	log *zerolog.Logger
}

func NewInvoicesGrpcService(log *zerolog.Logger) InvoicesGrpcService {
	return InvoicesGrpcService{
		log: log,
	}
}

func (inv InvoicesGrpcService) CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*InvoiceResponse, error) {
	return &InvoiceResponse{
		InvoiceId: "12345",
		Message:   "Invoice Created!",
	}, nil
}

func (inv InvoicesGrpcService) UpdateInvoice(ctx context.Context, req *UpdateInvoiceRequest) (*InvoiceResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (i InvoicesGrpcService) DeleteInvoice(ctx context.Context, req *DeleteInvoiceRequest) (*InvoiceResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (i InvoicesGrpcService) ListSentInvoices(ctx context.Context, req *ListInvoicesRequest) (*ListInvoicesResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (i InvoicesGrpcService) ListReceivedInvoices(ctx context.Context, req *ListInvoicesRequest) (*ListInvoicesResponse, error) {
	panic("not implemented") // TODO: Implement
}
