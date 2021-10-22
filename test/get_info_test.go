package ucaller_test

import (
	"fmt"
	"strconv"
	"testing"

	u "github.com/AleksandrMac/ucaller"
	mock_ucaller "github.com/AleksandrMac/ucaller/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	type mockBehavior func(uc *mock_ucaller.MockRequester, s *u.Service, uid u.ID)

	testTable := []struct {
		name                 string
		inputUid             u.ID
		mockBehavior         mockBehavior
		expectedResponseInfo *u.ResponseInfo
		expectedError        error
	}{
		{
			name:     "Ok",
			inputUid: UID,
			mockBehavior: func(uc *mock_ucaller.MockRequester, s *u.Service, uid u.ID) {
				uc.EXPECT().Get(fmt.Sprintf("/getInfo?service_id=%d&key=%s&ucaller_id=%d", InputData.ID, InputData.SecretKey, uid)).Return(
					[]byte(`{`+
						`"status":true`+
						`,"ucaller_id":`+strconv.Itoa(int(UID))+
						`,"init_time":1634708599`+
						`,"call_status":-1`+
						`,"code":"`+strconv.Itoa(int(code))+`"`+
						`,"is_repeated":false`+
						`,"repeatable":false`+
						`,"repeated_ucaller_ids":[`+
						`		103000`+
						`,		103001`+
						`]`+
						`,"unique":"`+unique+`"`+
						`,"client":"`+client+`"`+
						`,"phone":`+strconv.Itoa(int(phone))+
						`,"country_code":"RU"`+
						`,"country_image":"https:\/\/static.ucaller.ru\/flag\/ru.svg"`+
						`,"phone_info":[`+
						`	{`+
						`		"operator":"\u0411\u0438\u043b\u0430\u0439\u043d"`+ // Билайн
						`,		"region":"\u0410\u043c\u0443\u0440\u0441\u043a\u0430\u044f \u043e\u0431\u043b\u0430\u0441\u0442\u044c"`+ // Амурская область
						`,		"mnp":"\u041c\u0422\u0421`+ // МТС
						`	}`+
						`]`+
						`,"cost":0}`),
					nil)
			},
			expectedResponseInfo: &u.ResponseInfo{
				Status:     true,
				ID:         103001,
				InitTime:   1634708599,
				CallStatus: -1,
				Code:       uint64(code),
				IsRepeated: false,
				Repeatable: false,
				RepeatedUIDs: []u.ID{
					103000,
					103001,
				},
				UniqueRequestID: uuid.MustParse(unique),
				Client:          client,
				Phone:           phone,
				CountryCode:     "RU",
				CountryImage:    "https://static.ucaller.ru/flag/ru.svg",
				PhoneInfo: []u.PhoneInfo{
					{
						Operator: "Билайн",
						Region:   "Амурская область",
						MNP:      "МТС",
					},
				},
				Cost: 0,
			},
			expectedError: nil,
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
			testCase.mockBehavior(requester, service, UID)
			resp, err := service.GetInfo(UID)
			if err != nil {
				return
			}
			assert.Equal(t, testCase.expectedResponseInfo.Status, resp.Status)
			assert.Equal(t, testCase.expectedResponseInfo.ID, resp.ID)
			assert.Equal(t, testCase.expectedResponseInfo.InitTime, resp.InitTime)
			assert.Equal(t, testCase.expectedResponseInfo.CallStatus, resp.CallStatus)
			assert.Equal(t, testCase.expectedResponseInfo.Code, resp.Code)
			assert.Equal(t, testCase.expectedResponseInfo.IsRepeated, resp.IsRepeated)
			assert.Equal(t, testCase.expectedResponseInfo.Repeatable, resp.Repeatable)

			assert.Equal(t, len(testCase.expectedResponseInfo.RepeatedUIDs), len(resp.RepeatedUIDs))
			for i := range testCase.expectedResponseInfo.RepeatedUIDs {
				assert.Equal(t, testCase.expectedResponseInfo.RepeatedUIDs[i], resp.RepeatedUIDs[i])
			}
			assert.Equal(t, testCase.expectedResponseInfo.UniqueRequestID, resp.UniqueRequestID)
			assert.Equal(t, testCase.expectedResponseInfo.Client, resp.Client)
			assert.Equal(t, testCase.expectedResponseInfo.Phone, resp.Phone)
			assert.Equal(t, testCase.expectedResponseInfo.CountryCode, resp.CountryCode)
			assert.Equal(t, testCase.expectedResponseInfo.CountryImage, resp.CountryImage)

			assert.Equal(t, len(testCase.expectedResponseInfo.PhoneInfo), len(resp.PhoneInfo))
			for i := range testCase.expectedResponseInfo.PhoneInfo {
				assert.Equal(t, testCase.expectedResponseInfo.PhoneInfo[i].Operator, resp.PhoneInfo[i].Operator)
				assert.Equal(t, testCase.expectedResponseInfo.PhoneInfo[i].Region, resp.PhoneInfo[i].Region)
				assert.Equal(t, testCase.expectedResponseInfo.PhoneInfo[i].MNP, resp.PhoneInfo[i].MNP)
			}
			assert.Equal(t, testCase.expectedResponseInfo.Cost, resp.Cost)
		})
	}
}
