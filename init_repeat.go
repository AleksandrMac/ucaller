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
	Phone uint64 `json:"phone"`

	// код, который будет последними цифрами в номере телефона
	Code uint16 `json:"code"`

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
	s.RLock()
	_, ok := s.Calls[uid]
	s.RUnlock()
	if !ok {
		return nil, fmt.Errorf("ucaller_id не найден или истек")
	}
	if len(s.Calls[uid]) > int(s.freeRepeatNumber) {
		clearCalls(uid, s)
		return nil, fmt.Errorf("превышено количество бесплатных повторов")
	}

	body, err := s.Get(fmt.Sprintf("/initRepeat?service_id=%d&key=%s&ucaller_id=%d", s.id, s.secretKey, uid))
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
	s.Lock()
	// добавляем r.ID повтора, в цепочку звонков основаных входящим uid
	s.Calls[uid] = append(s.Calls[uid], r.ID)
	s.Unlock()
	return r, nil
}
