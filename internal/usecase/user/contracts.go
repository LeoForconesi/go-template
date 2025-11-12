package useruc

import "context"

type CreateUsecase interface {
	Execute(ctx context.Context, in CreateInput) (CreateOutput, error)
}

type GetUsecase interface {
	Execute(ctx context.Context, in GetInput) (GetOutput, error)
}

type ListUsecase interface {
	Execute(ctx context.Context, in ListInput) (ListOutput, error)
}

type UpdateUsecase interface {
	Execute(ctx context.Context, in UpdateInput) (UpdateOutput, error)
}

type DeleteUsecase interface {
	Execute(ctx context.Context, in DeleteInput) error
}

type NotifyUsecase interface {
	Execute(ctx context.Context, in NotifyInput) error
}
