package category_service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	category_service_mocks "github.com/AlexMickh/coledzh-shop-backend/internal/services/category/__mocks__"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateCategory(t *testing.T) {
	type fields struct {
		repository Repository
		cash       Cash
	}
	type args struct {
		ctx  context.Context
		name string
	}

	mDb := category_service_mocks.NewMockRepository(t)
	mCash := category_service_mocks.NewMockCash(t)

	tests := []struct {
		name            string
		fields          fields
		args            args
		wantDbMockErr   error
		wantCashMockErr error
		wantErr         error
	}{
		{
			name: "good case",
			fields: fields{
				repository: mDb,
				cash:       mCash,
			},
			args: args{
				ctx:  context.Background(),
				name: "fvewwecvn",
			},
			wantDbMockErr:   nil,
			wantCashMockErr: nil,
			wantErr:         nil,
		},
		{
			name: "failed to save in db case",
			fields: fields{
				repository: mDb,
				cash:       mCash,
			},
			args: args{
				ctx:  context.Background(),
				name: "fvewwecvn",
			},
			wantDbMockErr:   errors.New("failed to save"),
			wantCashMockErr: nil,
			wantErr:         errors.New("services.category.CreateCategory: failed to save"),
		},
		{
			name: "failed to save in cash case",
			fields: fields{
				repository: mDb,
				cash:       mCash,
			},
			args: args{
				ctx:  context.Background(),
				name: "fvewwecvn",
			},
			wantDbMockErr:   nil,
			wantCashMockErr: errors.New("failed to save"),
			wantErr:         errors.New("services.category.CreateCategory: failed to save"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mDb.EXPECT().SaveCategory(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(tt.wantDbMockErr)

			mCash.EXPECT().SaveCategory(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("models.Category"),
			).Return(tt.wantCashMockErr)

			fmt.Println(mCash)

			s := &Service{
				repository: tt.fields.repository,
				cash:       tt.fields.cash,
			}
			got, err := s.CreateCategory(tt.args.ctx, tt.args.name)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("Service.CreateCategory() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if got == "" {
				t.Errorf("Service.CreateCategory() = id is empty")
			}
		})
	}
}
