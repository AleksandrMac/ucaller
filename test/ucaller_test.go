package ucaller_test

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	mock_ucaller "github.com/AleksandrMac/ucaller/mocks"
	u "github.com/AleksandrMac/ucaller/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var UID u.ID = 103000
var InputData u.InputData = u.InputData{
	SecretKey:        "zxcvbasdfgqwertZXCVBASDFGQWERT23",
	ID:               id,
	FreeRepeatTime:   5 * time.Minute,
	FreeRepeatNumber: 2,
}

var (
	id        uint32 = 1616
	secretKey string = "zxcvbasdfgqwertZXCVBASDFGQWERT23"
	phone     uint64 = 79991234567
	code      uint16 = 7777
	client    string = "TestClient"
	unique    string = "f32d7ab0-2695-44ee-a20c-a34262a06b90"
)

// тестирование отработки очистки мапы от истекших  ucaller_id
func TestUCallerClearCalls(t *testing.T) {
	type mockBehavior func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall, uids []u.ID)

	testTable := []struct {
		name          string
		uids          []u.ID
		inputData     *u.InputData
		inputInitCall *u.InitCall
		sleepTime     time.Duration
		mockBehavior  mockBehavior
		expectedError []error
	}{
		{
			name: "Ok",
			inputInitCall: &u.InitCall{
				Phone:  &phone,
				Code:   &code,
				Client: &client,
				Unique: &unique,
			},
			inputData: &u.InputData{
				SecretKey:        secretKey,
				ID:               id,
				FreeRepeatTime:   3 * time.Millisecond,
				FreeRepeatNumber: 2,
			},
			uids: []u.ID{
				103000,
				103000,
				103000,
			},
			sleepTime: time.Millisecond,
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall, uids []u.ID) {
				url, _ := u.GetInitCallURL(s, ic)
				uc.EXPECT().Get(url).Return(
					[]byte(`{`+
						`"status": true`+
						`,"ucaller_id":`+strconv.Itoa(int(uids[0]))+
						`,"phone":`+strconv.Itoa(int(phone))+
						`,"code":`+strconv.Itoa(int(code))+
						`,"client":"`+client+`"`+
						`,"unique_request_id":"`+unique+`"`+
						`,"exists": true}`),
					nil)
				for i := 1; i < len(uids); i++ {
					uc.EXPECT().Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&ucaller_id=%d", InputData.ID, InputData.SecretKey, uids[i])).Return(
						[]byte(`{`+
							`"status":true`+
							`,"ucaller_id":`+strconv.Itoa(int(uids[i]))+
							`,"phone":`+strconv.Itoa(int(phone))+
							`,"code":`+strconv.Itoa(int(code))+
							`,"client":"`+client+`"`+
							`,"unique_request_id":"`+unique+`"`+
							`,"exists":true`+
							`,"free_repeated":true}`),
						nil)
				}
			},
			expectedError: []error{nil, nil, nil},
		},
		{
			name: "ErrorOutOfTime",
			inputInitCall: &u.InitCall{
				Phone:  &phone,
				Code:   &code,
				Client: &client,
				Unique: &unique,
			},
			inputData: &u.InputData{
				SecretKey:        secretKey,
				ID:               id,
				FreeRepeatTime:   2 * time.Millisecond,
				FreeRepeatNumber: 2,
			},
			uids: []u.ID{
				103001,
				103001,
				103001,
			},
			sleepTime: time.Millisecond,
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall, uids []u.ID) {
				url, _ := u.GetInitCallURL(s, ic)
				uc.EXPECT().Get(url).Return(
					[]byte(`{`+
						`"status": true`+
						`,"ucaller_id":`+strconv.Itoa(int(uids[0]))+
						`,"phone":`+strconv.Itoa(int(phone))+
						`,"code":`+strconv.Itoa(int(code))+
						`,"client":"`+client+`"`+
						`,"unique_request_id":"`+unique+`"`+
						`,"exists": true}`),
					nil)
				for i := 1; i < len(uids)-1; i++ {
					uc.EXPECT().Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&ucaller_id=%d", InputData.ID, InputData.SecretKey, uids[i])).Return(
						[]byte(`{`+
							`"status":true`+
							`,"ucaller_id":`+strconv.Itoa(int(uids[i]))+
							`,"phone":`+strconv.Itoa(int(phone))+
							`,"code":`+strconv.Itoa(int(code))+
							`,"client":"`+client+`"`+
							`,"unique_request_id":"`+unique+`"`+
							`,"exists":true`+
							`,"free_repeated":true}`),
						nil)
				}
			},
			expectedError: []error{
				nil,
				nil,
				fmt.Errorf("ucaller_id не найден или истек")},
		},
		{
			name: "ErrorOutOfRange",
			inputInitCall: &u.InitCall{
				Phone:  &phone,
				Code:   &code,
				Client: &client,
				Unique: &unique,
			},
			inputData: &u.InputData{
				SecretKey:        secretKey,
				ID:               id,
				FreeRepeatTime:   3 * time.Millisecond,
				FreeRepeatNumber: 1,
			},
			uids: []u.ID{
				103002,
				103002,
				103002,
			},
			sleepTime: time.Millisecond,
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, ic *u.InitCall, uids []u.ID) {
				url, _ := u.GetInitCallURL(s, ic)
				uc.EXPECT().Get(url).Return(
					[]byte(`{`+
						`"status": true`+
						`,"ucaller_id":`+strconv.Itoa(int(uids[0]))+
						`,"phone":`+strconv.Itoa(int(phone))+
						`,"code":`+strconv.Itoa(int(code))+
						`,"client":"`+client+`"`+
						`,"unique_request_id":"`+unique+`"`+
						`,"exists": true}`),
					nil)
				for i := 1; i < len(uids)-1; i++ {
					uc.EXPECT().Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&ucaller_id=%d", InputData.ID, InputData.SecretKey, uids[i])).Return(
						[]byte(`{`+
							`"status":true`+
							`,"ucaller_id":`+strconv.Itoa(int(uids[i]))+
							`,"phone":`+strconv.Itoa(int(phone))+
							`,"code":`+strconv.Itoa(int(code))+
							`,"client":"`+client+`"`+
							`,"unique_request_id":"`+unique+`"`+
							`,"exists":true`+
							`,"free_repeated":true}`),
						nil)
				}
			},
			expectedError: []error{
				nil,
				nil,
				fmt.Errorf("превышено количество бесплатных повторов")},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			requester := mock_ucaller.NewMockRequester(c)
			service, err := u.New(testCase.inputData, requester)
			if err != nil {
				fmt.Println(err)
				return
			}
			testCase.mockBehavior(requester, service, testCase.inputInitCall, testCase.uids)

			_, err = service.InitCall(testCase.inputInitCall)
			assert.Equal(t, testCase.expectedError[0], err)

			for i := 1; i < len(testCase.uids); i++ {
				time.Sleep(testCase.sleepTime)
				_, err = service.InitRepeat(testCase.uids[i])
				assert.Equal(t, testCase.expectedError[i], err, i)
			}

			time.Sleep(testCase.inputData.FreeRepeatTime)
			var mux sync.RWMutex
			mux.RLock()
			ids := service.Calls(testCase.uids[0])
			mux.RUnlock()
			var u []u.ID
			assert.Equal(t, u, ids)
		})
	}
}
