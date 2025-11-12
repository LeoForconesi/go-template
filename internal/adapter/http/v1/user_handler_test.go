package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/LeonardoForconesi/go-template/internal/adapter/http/v1"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	useruc "github.com/LeonardoForconesi/go-template/internal/usecase/user"
)

type mockCreator struct{ mock.Mock }

func (m *mockCreator) Execute(ctx context.Context, in useruc.CreateInput) (useruc.CreateOutput, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(useruc.CreateOutput), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var mc mockCreator                 // ← sin literal compuesto
	h := &v1.UserHandlers{Create: &mc} // ← pasás &mc

	router := gin.New()
	api := router.Group("/api/v1")
	h.Register(api)

	input := useruc.CreateInput{Name: "Leo", Email: "leo@email.com", Phone: "123"}
	expected := useruc.CreateOutput{
		ID:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		Phone:     input.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mc.On("Execute", mock.Anything, input).Return(expected, nil)

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var got useruc.CreateOutput
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	require.Equal(t, expected.Email, got.Email)

	mc.AssertExpectations(t)
}
