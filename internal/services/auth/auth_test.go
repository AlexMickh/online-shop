package auth_service

import (
	"context"
	"errors"
	"testing"

	"github.com/AlexMickh/coledzh-shop-backend/internal/models"
	auth_service_mocks "github.com/AlexMickh/coledzh-shop-backend/internal/services/auth/__mocks__"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx      context.Context
		login    string
		email    string
		password string
	}

	m := auth_service_mocks.NewMockStorage(t)

	m.EXPECT().SaveUser(
		mock.AnythingOfType("context.backgroundCtx"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return("some id", nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				storage: m,
			},
			args: args{
				ctx:      context.Background(),
				login:    "sas123",
				email:    "sas@gmail.com",
				password: "sas123",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			got, err := s.Register(tt.args.ctx, tt.args.login, tt.args.email, tt.args.password)
			if err != nil {
				if errors.Is(err, tt.wantErr) {
					t.Errorf("Service.Register() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got == "" {
				t.Error("Service.Register() = id is empty")
			}
		})
	}
}

func TestService_RegisterAdmin(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx      context.Context
		login    string
		email    string
		password string
	}

	m := auth_service_mocks.NewMockStorage(t)

	m.EXPECT().SaveAdmin(
		mock.AnythingOfType("context.backgroundCtx"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
	).Return(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				storage: m,
			},
			args: args{
				ctx:      context.Background(),
				login:    "sas123",
				email:    "sas@gmail.com",
				password: "sas123",
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			if err := s.RegisterAdmin(tt.args.ctx, tt.args.login, tt.args.email, tt.args.password); err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Service.RegisterAdmin() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestService_Login(t *testing.T) {
	type fields struct {
		storage      Storage
		sessionStore SessionStore
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}

	ms := auth_service_mocks.NewMockStorage(t)
	mc := auth_service_mocks.NewMockSessionStore(t)

	tests := []struct {
		name           string
		fields         fields
		args           args
		wantErr        error
		wantStorageErr error
		wantCashErr    error
	}{
		{
			name: "good case",
			fields: fields{
				storage:      ms,
				sessionStore: mc,
			},
			args: args{
				ctx:      context.Background(),
				email:    "sas@gmail.com",
				password: "test123",
			},
			wantErr:        nil,
			wantStorageErr: nil,
			wantCashErr:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(tt.args.password), bcrypt.MinCost)

			ms.EXPECT().UserByEmail(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(models.User{
				ID:              "test",
				Login:           "test",
				Email:           "test",
				Password:        string(hashPassword),
				Role:            "test",
				IsEmailVerified: true,
			}, tt.wantStorageErr)

			mc.EXPECT().SaveSession(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("models.User"),
			).Return(tt.wantCashErr)

			s := &Service{
				storage:      tt.fields.storage,
				sessionStore: tt.fields.sessionStore,
			}
			got, err := s.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Service.Login() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got == "" {
				t.Error("Service.Login() = id is empty")
			}
		})
	}
}
