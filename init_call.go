package ucaller

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type InitCall struct {
	// REQUIRED. Без + с кодом страны. Номер телефона пользователя, которому будет совершен звонок с авторизацией.
	// Цифровое значение
	Phone *uint64 `json:"phone"`

	// 4 цифры. На эти 4 цифры будет заканчиваться исходящий от нас звонок к вашему пользователю.
	// Если `code` не передан, он будет автоматически сгенерирован нами.
	// Цифровое значение
	Code *uint16 `json:"code"`

	// До 64 символов. Можете передать информацию о вашем пользователе, например его никнейм, E-mail и так далее.
	// В статистике сервиса сможете найти ту или иную авторизацию по этому параметру.
	// Набор символов
	Client *string `json:"client"`

	// До 64 символов. Для активации идемпотентности запроса.
	// Рекомендуем использовать V4 UUID.
	// Набор символов
	Unique *string `json:"unique"`
}

type ResponseInitCall struct {
	// true в случае успеха, false в случае неудачи
	Status bool `json:"status"`

	// уникальный ID в системе uCaller, который позволит проверять статус и инициализировать метод initRepeat
	ID ID `json:"ucaller_id"`

	// номер телефона, куда мы совершили звонок
	Phone string `json:"phone"`

	// код, который будет последними цифрами в номере телефона
	Code interface{} `json:"code"`

	// идентификатор пользователя переданный клиентом
	Client string `json:"client"`

	// появляется только если вами был передан параметр `unique`
	UniqueRequestID uuid.UUID `json:"unique_request_id"`

	// появляется при переданном параметре `unique`, если такой запрос уже был инициализирован ранее
	Exists bool `json:"exists"`

	Error string `json:"error"`
}

func (s *Service) InitCall(ic *InitCall) (r *ResponseInitCall, err error) {
	url, err := GetInitCallURL(s, ic)
	if err != nil {
		return nil, err
	}

	body, err := s.Get(url)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return r, errors.New(r.Error)
	}

	s.calls[r.ID] = append(s.calls[r.ID], r.ID)
	// запускаем горутину которая удалит информаци о цепочке звонков с истекшей
	// возможностью для инициализации повтора звонка
	go func() {
		time.Sleep(s.freeRepeatTime)
		clearCalls(r.ID, s)
	}()

	//     Авторизация успешно инициализирована
	//     Мы совершаем звонок на номер +79991234567
	//     Последние четыре цифры номера, с которого мы совершаем звонок, заканчиваются на 7777
	return r, nil
}

func GetInitCallURL(s *Service, ic *InitCall) (string, error) {
	if ic.Phone == nil {
		return "", errors.New("phone is required")
	}
	url := fmt.Sprintf("/initCall?service_id=%d&key=%s&phone=%d", s.id, s.secretKey, *ic.Phone)

	if ic.Code != nil {
		if *ic.Code > MaxCode {
			return "", errors.New("'code' не должен превышать 9999")
		}
		url += fmt.Sprintf("&code=%d", *ic.Code)
	}
	if ic.Client != nil {
		url += fmt.Sprintf("&client=%s", *ic.Client)
	}
	if ic.Unique == nil {
		id := uuid.NewString()
		ic.Unique = &id
	}
	url += fmt.Sprintf("&unique=%s", *ic.Unique)
	return url, nil
}

// clearCalls очищает цепочку звонков начатых с firstUCallerID
func clearCalls(uid ID, s *Service) {
	s.mux.RLock()
	_, ok := s.calls[uid]
	s.mux.RUnlock()
	if ok {
		s.mux.Lock()
		delete(s.calls, uid)
		s.mux.Unlock()
	}
}
