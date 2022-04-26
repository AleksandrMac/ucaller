package ucaller_test

import (
	"fmt"
	"testing"

	u "github.com/AleksandrMac/ucaller"
	mock_ucaller "github.com/AleksandrMac/ucaller/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetService(t *testing.T) {
	type mockBehavior func(uc *mock_ucaller.MockRequester, s *u.Service)

	testTable := []struct {
		name                    string
		mockBehavior            mockBehavior
		expectedResponseService *u.ResponseService
	}{
		{
			name: "Ok",
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service) {
				uc.EXPECT().Get(fmt.Sprintf("/getService?service_id=%d&key=%s", InputData.ID, InputData.SecretKey)).Return(
					[]byte(`{`+
						`"status":true,`+
						`"service_status":1,`+
						`"name":"PhoneConfirm",`+
						`"creation_time":1632703297,`+
						`"last_request":1634773693,`+
						`"owner":"example@mail.ru",`+
						`"use_direction":"example.com",`+
						`"now_test":true,`+
						`"test_info":`+
						`{`+
						`	"test_requests":98,`+
						`	"verified_phone":79998881111`+
						`}`), nil)
			},
			expectedResponseService: &u.ResponseService{
				Status:        true,
				ServiceStatus: 1,
				Name:          "PhoneConfirm",
				CreationTime:  1632703297,
				LastRequest:   1634773693,
				Owner:         "example@mail.ru",
				UseDirection:  "example.com",
				NowTest:       true,
				TestInfo: u.TestInfo{
					TestRequests:  98,
					VerifiedPhone: 79998881111,
				},
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
			service.AddCalls(UID)
			testCase.mockBehavior(requester, service)
			resp, err := service.GetService()
			if err != nil {
				return
			}
			assert.Equal(t, testCase.expectedResponseService.Status, resp.Status)
			assert.Equal(t, testCase.expectedResponseService.ServiceStatus, resp.ServiceStatus)
			assert.Equal(t, testCase.expectedResponseService.Name, resp.Name)
			assert.Equal(t, testCase.expectedResponseService.CreationTime, resp.CreationTime)
			assert.Equal(t, testCase.expectedResponseService.LastRequest, resp.LastRequest)
			assert.Equal(t, testCase.expectedResponseService.Owner, resp.Owner)
			assert.Equal(t, testCase.expectedResponseService.UseDirection, resp.UseDirection)
			assert.Equal(t, testCase.expectedResponseService.NowTest, resp.NowTest)
			assert.Equal(t, testCase.expectedResponseService.TestInfo.TestRequests, resp.TestInfo.TestRequests)
			assert.Equal(t, testCase.expectedResponseService.TestInfo.VerifiedPhone, resp.TestInfo.VerifiedPhone)
		})
	}
}
