package useruc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/LeonardoForconesi/go-template/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	duser "github.com/LeonardoForconesi/go-template/internal/domain/user"
	usecase "github.com/LeonardoForconesi/go-template/internal/usecase/user"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) GetByID(ctx context.Context, id uuid.UUID) (*duser.User, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) List(ctx context.Context, page, size int) ([]duser.User, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) Update(ctx context.Context, u *duser.User) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) Delete(ctx context.Context, id uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) GetByEmail(ctx context.Context, email string) (*duser.User, error) {
	args := m.Called(ctx, email)
	if v := args.Get(0); v != nil {
		return v.(*duser.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRepo) Create(ctx context.Context, u *duser.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

type mockPublisher struct{ mock.Mock }

func (m *mockPublisher) PublishUserCreated(ctx context.Context, evt duser.CreatedEvent) error {
	args := m.Called(ctx, evt)
	return args.Error(0)
}

func TestCreator_Execute_Success(t *testing.T) {
	repo := new(mockRepo)
	pub := new(mockPublisher)
	uc := &usecase.Creator{Repo: repo, Publisher: pub}

	in := usecase.CreateInput{
		Name:  "Leo",
		Email: "leo@email.com",
		Phone: "123456",
	}

	repo.On("GetByEmail", mock.Anything, in.Email).Return(nil, errors.New("not found"))
	repo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
	pub.On("PublishUserCreated", mock.Anything, mock.AnythingOfType("user.CreatedEvent")).Return(nil)

	out, err := uc.Execute(context.Background(), in)
	require.NoError(t, err)
	require.Equal(t, in.Name, out.Name)
	require.Equal(t, in.Email, out.Email)
	require.NotEqual(t, uuid.Nil, out.ID)

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestCreator_Execute_Duplicate(t *testing.T) {
	repo := new(mockRepo)
	pub := new(mockPublisher)
	uc := &usecase.Creator{Repo: repo, Publisher: pub}

	existing := &duser.User{
		ID:    uuid.New(),
		Name:  "Existing",
		Email: "leo@email.com",
	}

	repo.On("GetByEmail", mock.Anything, existing.Email).Return(existing, nil)

	in := usecase.CreateInput{Name: "Leo", Email: existing.Email, Phone: "123"}

	_, err := uc.Execute(context.Background(), in)
	require.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestCreator_Execute_InvalidUser(t *testing.T) {
	repo := new(mockRepo)
	pub := new(mockPublisher)
	uc := &usecase.Creator{Repo: repo, Publisher: pub}

	in := usecase.CreateInput{Email: "sin_nombre@email.com"} // sin Name

	_, err := uc.Execute(context.Background(), in)
	require.ErrorIs(t, err, domain.ErrInvalidArgument)
}

func TestCreator_Execute_InvalidDomain(t *testing.T) {
	repo := &mockRepo{}
	pub := &mockPublisher{}
	uc := &usecase.Creator{Repo: repo, Publisher: pub}

	in := usecase.CreateInput{Name: "", Email: "leo@email.com"}

	out, err := uc.Execute(context.Background(), in)
	require.ErrorIs(t, err, domain.ErrInvalidArgument)
	require.Equal(t, uuid.Nil, out.ID)

	repo.AssertNotCalled(t, "Create")
	pub.AssertNotCalled(t, "PublishUserCreated")
}

func TestCreator_Execute_CreateError(t *testing.T) {
	repo := &mockRepo{}
	pub := &mockPublisher{}
	uc := &usecase.Creator{Repo: repo, Publisher: pub}

	in := usecase.CreateInput{Name: "Leo", Email: "leo@email.com", Phone: "123"}

	repo.On("GetByEmail", mock.Anything, in.Email).Return(nil, errors.New("not found"))
	repo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(errors.New("db fail"))

	_, err := uc.Execute(context.Background(), in)
	require.EqualError(t, err, "db fail")

	pub.AssertNotCalled(t, "PublishUserCreated")
}
