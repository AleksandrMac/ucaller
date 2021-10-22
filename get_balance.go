package ucaller

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ResponseBalance struct {
	// true в случае успеха, false в случае неудачи
	Status bool `json:"status"`

	// Остаточный баланс на рублевом счете аккаунта
	Rub float32 `json:"rub_balance"`

	// Остаточный бонусный баланс
	Bonus float32 `json:"bonus_balance"`

	// Кодовое значение вашего тарифного плана
	Tariff string `json:"tariff"`

	// Название тарифного плана
	TariffName string `json:"tariff_name"`

	Error string `json:"error"`
}

func (s *Service) GetBalance() (r *ResponseBalance, err error) {
	body, err := s.Get(fmt.Sprintf("/getBalance?service_id=%d&key=%s", s.id, s.secretKey))
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
