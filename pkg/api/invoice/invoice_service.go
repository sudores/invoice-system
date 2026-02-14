package invoice

import (
	context "context"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	invoiceRepo "github.com/sudores/invoice-system/pkg/repo/invoice"
	grpc "google.golang.org/grpc"
)

type InvoiceRepo interface {
	Create(ctx context.Context, u *invoiceRepo.Invoice) (*invoiceRepo.Invoice, error)
	Update(ctx context.Context, id uuid.UUID, u *invoiceRepo.Invoice) (*invoiceRepo.Invoice, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Get(ctx context.Context, id uuid.UUID) (*invoiceRepo.Invoice, error)
	GetByUser(ctx context.Context, uuid uuid.UUID) (*invoiceRepo.Invoice, error)
}

type InvoiceGrpcService struct {
	UnimplementedInvoiceServiceServer
	log  *zerolog.Logger
	repo InvoiceRepo
}

func (ugs InvoiceGrpcService) Descriptor() *grpc.ServiceDesc {
	return &InvoiceService_ServiceDesc
}

func (igs *InvoiceGrpcService) RegisterHttp(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterInvoiceServiceHandlerServer(ctx, mux, igs)
}

func NewInvoiceGrpcService(log *zerolog.Logger, repo InvoiceRepo) *InvoiceGrpcService {
	return &InvoiceGrpcService{
		log:  log,
		repo: repo,
	}
}

func (igs InvoiceGrpcService) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	panic("not implemented") // TODO: Implement
}

func (igs InvoiceGrpcService) Update(ctx context.Context, _ *UpdateReq) (*UpdateResp, error) {
	panic("not implemented") // TODO: Implement
}

func (igs InvoiceGrpcService) Get(ctx context.Context, _ *GetReq) (*GetResp, error) {
	panic("not implemented") // TODO: Implement
}

func (igs InvoiceGrpcService) List(ctx context.Context, _ *ListReq) (*ListResp, error) {
	panic("not implemented") // TODO: Implement
}

func (igs InvoiceGrpcService) ChangeStatus(ctx context.Context, _ *ChangeStatusReq) (*ChangeStatusResp, error) {
	panic("not implemented") // TODO: Implement
}

func (igs InvoiceGrpcService) Delete(ctx context.Context, _ *DeleteReq) (*DeleteResp, error) {
	panic("not implemented") // TODO: Implement
}
