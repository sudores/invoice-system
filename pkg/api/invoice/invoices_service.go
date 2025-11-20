package invoice

import (
	context "context"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	invoiceRepo "github.com/sudores/invoice-system/pkg/repo/invoice"
)

type InvoiceRepo interface {
	CreateInvoice(ctx context.Context, inv *invoiceRepo.Invoice) (*invoiceRepo.Invoice, error)
	DeleteInvoice(ctx context.Context, id uuid.UUID) error
	UpdateInvoice(ctx context.Context, id uuid.UUID, inv *invoiceRepo.Invoice) (*invoiceRepo.Invoice, error)
	GetInvoice(ctx context.Context, id uuid.UUID) (*invoiceRepo.Invoice, error)
}

type InvoicesGrpcService struct {
	UnimplementedInvoiceServiceServer
	log  *zerolog.Logger
	repo InvoiceRepo
}

func NewInvoicesGrpcService(log *zerolog.Logger, repo InvoiceRepo) InvoicesGrpcService {
	return InvoicesGrpcService{
		log:  log,
		repo: repo,
	}
}

func (inv InvoicesGrpcService) CreateInvoice(ctx context.Context, req *CreateInvoiceRequest) (*InvoiceResponse, error) {
	invoice, err := grpcCreateInvoiceReqToInvoice(req)
	if err != nil {
		return nil, err
	}
	created, err := inv.repo.CreateInvoice(ctx, invoice)
	if err != nil {
		return &InvoiceResponse{
			InvoiceId: "",
			Message:   "Failed to create invoice " + err.Error(),
		}, err
	}
	return &InvoiceResponse{
		InvoiceId: created.InvoiceId.String(),
		Message:   "Invoice successfully created!",
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

func grpcCreateInvoiceReqToInvoice(req *CreateInvoiceRequest) (*invoiceRepo.Invoice, error) {
	senderId, err := uuid.Parse(req.SenderId)
	if err != nil {
		return nil, err
	}
	recipientId, err := uuid.Parse(req.RecipientId)
	if err != nil {
		return nil, err
	}
	return &invoiceRepo.Invoice{
		SenderId:    senderId,
		RecipientId: recipientId,
		Items:       grpcInvoiceItemsToInvoiceItems(req.Items),
		Description: req.Description,
		DueDate:     req.DueDate.AsTime(),
		CreatedAt:   time.Now(),
	}, nil
}

func grpcInvoiceItemsToInvoiceItems(items []*InvoiceItem) []*invoiceRepo.InvoiceItem {
	retItems := []*invoiceRepo.InvoiceItem{}
	for _, k := range items {
		retItems = append(retItems, &invoiceRepo.InvoiceItem{
			Title:       k.Title,
			Amount:      k.Amount,
			ReferenceId: k.ReferenceId,
		})
	}
	return retItems
}
