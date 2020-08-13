package transfer

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/zzell/transfer/cfg"
	"github.com/zzell/transfer/db/repo"
	"github.com/zzell/transfer/model"

	"github.com/zzell/transfer/currency/mock"
	mock2 "github.com/zzell/transfer/db/repo/mock"
)

func TestHandler_Handle(t *testing.T) {
	var (
		ctrl            = gomock.NewController(t)
		converterMock   = mock.NewMockConverter(ctrl)
		walletsRepoMock = mock2.NewMockWalletsRepo(ctrl)
		repoMock        = repo.Repository{WalletsRepo: walletsRepoMock}
		config          = cfg.Config{CommissionPercent: 1.5}
		target          = NewHandler(&config, repoMock, converterMock)
	)

	rand.Seed(time.Now().Unix())

	t.Run("positive", func(t *testing.T) {
		payload := model.TransferPayload{
			From:   rand.Int(),
			To:     rand.Int(),
			Amount: 100,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/transfer", bytes.NewReader(b))
		if err != nil {
			t.Fatal(err)
		}

		walletsRepoMock.EXPECT().GetWallet(gomock.Any()).Return(&model.Wallet{
			Currency: model.Currency{Symbol: "btc"},
			Score:    200,
		}, nil).AnyTimes()

		walletsRepoMock.EXPECT().Transfer(payload.From, payload.To, payload.Amount, gomock.Any()).Return(nil)

		http.HandlerFunc(target.Transfer).ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatal()
		}
	})
}
