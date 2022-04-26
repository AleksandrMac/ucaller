package ucaller_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	u "github.com/AleksandrMac/ucaller"
	mock_ucaller "github.com/AleksandrMac/ucaller/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInitCall(t *testing.T) {
	type mockBehavior func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall)

	testTable := []struct {
		name                     string
		inputInitCall            *u.InitCall
		mockBehavior             mockBehavior
		expectedResponseInitCall *u.ResponseInitCall
		expectedError            error
	}{
		{
			name: "Ok",
			inputInitCall: &u.InitCall{
				Phone:  &phone,
				Code:   &code,
				Client: &client,
				Unique: &unique,
			},
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall) {
				url, _ := u.GetInitCallURL(s, ic)
				uc.EXPECT().Get(url).Return(
					[]byte(`{`+
						`"status": true`+
						`,"ucaller_id": 103000`+
						`,"phone":"`+phoneMask+`"`+
						`,"code":"`+strconv.Itoa(int(code))+`"`+
						`,"client":"`+client+`"`+
						`,"unique_request_id":"`+unique+`"`+
						`,"exists": true}`),
					nil)
			},
			expectedResponseInitCall: &u.ResponseInitCall{
				Status:          true,
				ID:              103000,
				Phone:           phoneMask,
				Code:            fmt.Sprintf("%d", code),
				Client:          client,
				UniqueRequestID: uuid.MustParse(unique),
				Exists:          true,
			},
			expectedError: nil,
		},
		{
			name: "Error",
			inputInitCall: &u.InitCall{
				Phone:  &phone,
				Code:   &code,
				Client: &client,
				Unique: &unique,
			},
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall) {
				url, _ := u.GetInitCallURL(s, ic)
				uc.EXPECT().Get(url).Return(
					[]byte(`{`+
						`"status": false`+
						`,"error": "Ваш IP адрес заблокирован"`+
						`,"code":"0"}`),
					nil)
			},
			expectedResponseInitCall: &u.ResponseInitCall{
				Status: false,
				Code:   "0",
				Error:  "Ваш IP адрес заблокирован",
			},
			expectedError: errors.New("Ваш IP адрес заблокирован"),
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
			testCase.mockBehavior(requester, service, testCase.inputInitCall)
			resp, err := service.InitCall(testCase.inputInitCall)
			if err != nil {
				assert.Equal(t, testCase.expectedError, err)
				assert.Equal(t, testCase.expectedResponseInitCall.Status, resp.Status)
				assert.Equal(t, testCase.expectedResponseInitCall.Error, resp.Error)
				assert.Equal(t, testCase.expectedResponseInitCall.Code, resp.Code)
				return
			}
			assert.Equal(t, testCase.expectedResponseInitCall.Status, resp.Status)
			assert.Equal(t, testCase.expectedResponseInitCall.Client, resp.Client)
			assert.Equal(t, testCase.expectedResponseInitCall.Code, resp.Code)
			assert.Equal(t, testCase.expectedResponseInitCall.Exists, resp.Exists)
			assert.Equal(t, testCase.expectedResponseInitCall.Phone, resp.Phone)
			assert.Equal(t, testCase.expectedResponseInitCall.ID, resp.ID)
			assert.Equal(t, testCase.expectedResponseInitCall.UniqueRequestID, resp.UniqueRequestID)
		})
	}
}
