package ucaller_test

import (
	"fmt"
	"testing"

	u "github.com/AleksandrMac/ucaller"
	mock_ucaller "github.com/AleksandrMac/ucaller/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetBalance(t *testing.T) {
	type mockBehavior func(uc *mock_ucaller.MockRequester, s *u.Service)

	testTable := []struct {
		name                    string
		mockBehavior            mockBehavior
		expectedResponseBalance *u.ResponseBalance
	}{
		{
			name: "Ok",
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service) {
				uc.EXPECT().Get(fmt.Sprintf("/getBalance?service_id=%d&key=%s", InputData.ID, InputData.SecretKey)).Return(
					[]byte(`{`+
						`"status":true,`+
						`"rub_balance":100,`+
						`"bonus_balance":30,`+
						`"tariff":"start",`+
						`"tariff_name":"\u041d\u043e\u0432\u0438\u0447\u043e\u043a"`+ // Новичок
						`}`), nil)
			},
			expectedResponseBalance: &u.ResponseBalance{
				Status:     true,
				Rub:        100,
				Bonus:      30,
				Tariff:     "start",
				TariffName: "Новичок",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			requester := mock_ucaller.NewMockRequester(c)
			service, err := u.New(&InputData, requester)
			if err != nil {
				fmt.Println(err)
				return
			}
			service.Calls[UID] = append(service.Calls[UID], UID)
			testCase.mockBehavior(requester, service)
			resp, err := service.GetBalance()
			assert.Equal(t, err, nil)
			assert.Equal(t, testCase.expectedResponseBalance.Status, resp.Status)
			assert.Equal(t, testCase.expectedResponseBalance.Rub, resp.Rub)
			assert.Equal(t, testCase.expectedResponseBalance.Bonus, resp.Bonus)
			assert.Equal(t, testCase.expectedResponseBalance.Tariff, resp.Tariff)
			assert.Equal(t, testCase.expectedResponseBalance.TariffName, resp.TariffName)
		})
	}
}
