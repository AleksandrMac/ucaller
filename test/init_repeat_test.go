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
	"github.com/stretchr/testify/require"
)

func TestInitRepeat(t *testing.T) {
	type mockBehavior func(uc *mock_ucaller.MockRequester, s *u.Service, uid u.ID)

	testTable := []struct {
		name                       string
		inputUid                   u.ID
		mockBehavior               mockBehavior
		expectedResponseInitRepeat *u.ResponseInitRepeat
		expectedSizeServiceCalls   int
		expectedError              error
	}{
		{
			name:     "Ok",
			inputUid: UID,
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, uid u.ID) {
				uc.EXPECT().Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&ucaller_id=%d", InputData.ID, InputData.SecretKey, uid)).Return(
					[]byte(`{`+
						`"status":true`+
						`,"ucaller_id":`+strconv.Itoa(int(uid))+
						`,"phone":"`+phoneMask+`"`+
						`,"code":"`+strconv.Itoa(int(code))+`"`+
						`,"client":"`+client+`"`+
						`,"unique_request_id":"`+unique+`"`+
						`,"exists":true`+
						`,"free_repeated":true}`),
					nil)
			},
			expectedResponseInitRepeat: &u.ResponseInitRepeat{
				Status:          true,
				ID:              UID,
				Phone:           phoneMask,
				Code:            fmt.Sprintf("%d", code),
				Client:          client,
				UniqueRequestID: uuid.MustParse(unique),
				Exists:          true,
				FreeRepeated:    true,
			},
			expectedError:            nil,
			expectedSizeServiceCalls: 1,
		},
		{
			name:     "Error",
			inputUid: UID,
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, uid u.ID) {
				uc.EXPECT().Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&ucaller_id=%d", InputData.ID, InputData.SecretKey, uid)).Return(
					[]byte(`{`+
						`"status": false`+
						`,"error": "This uCaller ID is already repeated"`+
						`,"code":"13"}`),
					nil)
			},
			expectedResponseInitRepeat: &u.ResponseInitRepeat{
				Status: false,
				Code:   "13",
				Error:  "This uCaller ID is already repeated",
			},
			expectedError: errors.New("This uCaller ID is already repeated"),
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
			service.AddCalls(UID)
			testCase.mockBehavior(requester, service, UID)
			resp, err := service.InitRepeat(UID)
			if err != nil {
				require.Equal(t, testCase.expectedError, err)
				require.Equal(t, testCase.expectedResponseInitRepeat.Status, resp.Status)
				require.Equal(t, testCase.expectedResponseInitRepeat.Error, resp.Error)
				require.Equal(t, testCase.expectedResponseInitRepeat.Code, resp.Code)
				return
			}
			require.Equal(t, testCase.expectedResponseInitRepeat.Status, resp.Status)
			require.Equal(t, testCase.expectedResponseInitRepeat.Client, resp.Client)
			require.Equal(t, testCase.expectedResponseInitRepeat.Code, resp.Code)
			require.Equal(t, testCase.expectedResponseInitRepeat.Exists, resp.Exists)
			require.Equal(t, testCase.expectedResponseInitRepeat.Phone, resp.Phone)
			require.Equal(t, testCase.expectedResponseInitRepeat.ID, resp.ID)
			require.Equal(t, testCase.expectedResponseInitRepeat.UniqueRequestID, resp.UniqueRequestID)
			require.Equal(t, testCase.expectedResponseInitRepeat.FreeRepeated, resp.FreeRepeated)
		})
	}
}
