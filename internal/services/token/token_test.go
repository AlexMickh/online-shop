package token_service

import (
	"context"
	"errors"
	"testing"

	"github.com/AlexMickh/coledzh-shop-backend/internal/consts"
	"github.com/AlexMickh/coledzh-shop-backend/internal/errs"
	token_service_mocks "github.com/AlexMickh/coledzh-shop-backend/internal/services/token/__mocks__"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateToken(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		ctx       context.Context
		userId    string
		tokenType string
	}

	m := token_service_mocks.NewMockStorage(t)

	m.EXPECT().SaveToken(
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
				ctx:       context.Background(),
				userId:    "id",
				tokenType: consts.TokenTypeEmailVerify,
			},
			wantErr: nil,
		},
		{
			name: "wrong token type",
			fields: fields{
				storage: m,
			},
			args: args{
				ctx:       context.Background(),
				userId:    "id",
				tokenType: "wrong",
			},
			wantErr: errs.ErrWrongTokenType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			got, err := s.CreateToken(tt.args.ctx, tt.args.userId, tt.args.tokenType)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Service.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got == "" && tt.wantErr == nil {
				t.Error("Service.CreateToken() = token is empty")
			}
		})
	}
}
