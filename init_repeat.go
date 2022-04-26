package ucaller

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ResponseInitRepeat struct {
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

	// показывает, что осуществлена повторная авторизация
	FreeRepeated bool `json:"free_repeated"`

	Error string `json:"error"`
}

func (s *Service) InitRepeat(uid ID) (r *ResponseInitRepeat, err error) {
	s.mux.RLock()
	_, ok := s.calls[uid]
	s.mux.RUnlock()
	if !ok {
		return nil, fmt.Errorf("ucaller_id не найден или истек")
	}
	if len(s.calls[uid]) > int(s.freeRepeatNumber) {
		clearCalls(uid, s)
		return nil, fmt.Errorf("превышено количество бесплатных повторов")
	}

	body, err := s.Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&uid=%d", s.id, s.secretKey, uid))
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
	//     Повторная бесплатная авторизация успешно инициализирована
	s.mux.Lock()
	// добавляем r.ID повтора, в цепочку звонков основаных входящим uid
	s.calls[uid] = append(s.calls[uid], r.ID)
	s.mux.Unlock()
	return r, nil
}
