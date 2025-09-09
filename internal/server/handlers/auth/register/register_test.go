package register

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	register_mocks "github.com/AlexMickh/coledzh-shop-backend/internal/server/handlers/auth/register/__mocks__"
	"github.com/AlexMickh/coledzh-shop-backend/pkg/api"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name           string
		login          string
		email          string
		password       string
		respStatus     int
		regMockError   error
		tokenMockError error
		mailMockError  error
	}{
		{
			name:           "good case",
			login:          "sas123",
			email:          "sas123@gmail.com",
			password:       "sas123",
			respStatus:     http.StatusCreated,
			regMockError:   nil,
			tokenMockError: nil,
			mailMockError:  nil,
		},
		{
			name:           "not valid login",
			login:          "sa",
			email:          "sas123@gmail.com",
			password:       "sas123",
			respStatus:     http.StatusBadRequest,
			regMockError:   nil,
			tokenMockError: nil,
			mailMockError:  nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			regMock := register_mocks.NewMockRegisterer(t)
			regMock.EXPECT().Register(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(uuid.NewString(), tt.regMockError).Maybe()

			tokenMock := register_mocks.NewMockTokenCreator(t)
			tokenMock.EXPECT().CreateToken(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(uuid.NewString(), tt.tokenMockError).Maybe()

			mailMock := register_mocks.NewMockVerificationSender(t)
			mailMock.EXPECT().SendVerification(
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(tt.mailMockError).Maybe()

			handler := api.ErrorWrapper(New(validator.New(), regMock, tokenMock, mailMock))

			input := fmt.Sprintf(`{"login": "%s", "email": "%s", "password": "%s"}`, tt.login, tt.email, tt.password)

			req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respStatus, rr.Code)

			var resp Response
			err = json.NewDecoder(rr.Body).Decode(&resp)
			require.NoError(t, err)

			if tt.respStatus == http.StatusCreated {
				_, err = uuid.Parse(resp.ID)
				require.NoError(t, err)
			}
		})
	}
}
