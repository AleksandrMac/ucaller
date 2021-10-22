package ucaller

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResponseService struct {
	// true в случае успеха, false в случае неудачи
	Status bool `json:"status"`

	// ID сервиса
	ServiceStatus uint64 `json:"service_status"`

	// Название сервиса
	Name string `json:"name"`

	// Время создания сервиса в unix формате
	CreationTime uint64 `json:"creation_time"`

	// Время последнего не кэшированного обращения к API сервиса в unix формате
	LastRequest uint64 `json:"last_request"`

	// E-mail адрес владельца сервиса
	Owner string `json:"owner"`

	// Информация о том, где будет использоваться сервис
	UseDirection string `json:"use_direction"`

	// Состояние тестового режима на текущий момент
	NowTest bool `json:"now_test"`

	// Информация о тестовом состоянии
	TestInfo TestInfo `json:"test_info"`

	Error string `json:"error"`
}

type TestInfo struct {
	// Оставшееся количество бесплатных тестовых обращений
	TestRequests uint64 `json:"test_requests"`

	// Верифицированный номер телефона для тестовых обращений
	VerifiedPhone uint64 `json:"verified_phone"`
}

func (s *Service) GetService() (r *ResponseService, err error) {
	body, err := s.Get(fmt.Sprintf("/getService?service_id=%d&key=%s", s.id, s.secretKey))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return nil, errors.New(r.Error)
	}
	return
}
