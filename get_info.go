package ucaller

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type PhoneInfo struct {
	Operator string `json:"operator"` // Оператор связи
	Region   string `json:"region"`   // регион субъеккта Российской федерации
	MNP      string `json:"mnp"`      // Если у номера был сменен оператор - MNP покажет нового оператора
}

type ResponseInfo struct {
	// true в случае успеха, false в случае неудачи
	Status bool `json:"status"`

	// запрошенный uCaller ID
	ID ID `json:"ucaller_id"`

	// время, когда была инициализирована авторизация
	InitTime uint64 `json:"init_time"`

	// Статус звонка,
	// -1 = информация проверяется (от 1 сек до 1 минуты),
	//  0 = дозвониться не удалось,
	//  1 = звонок осуществлен
	CallStatus int8 `json:"call_status"`

	// является ли этот uCaller ID повтором (initRepeat), если да, будет добавлен first_ucaller_id с первым uCaller ID этой цепочки
	IsRepeated bool `json:"is_repeated"`

	// возможно ли инициализировать бесплатные повторы (initRepeat)
	Repeatable bool `json:"repeatable"`

	// Появляется в случае repeatable: true, говорит о количестве возможных повторов
	RepeatTimes uint8 `json:"repeat_times"`

	// цепочка  uCaller ID инициализированных повторов (initRepeat)
	RepeatedUIDs []ID `json:"repeated_ucaller_ids"`

	// ключ идемпотентности (если был передан)
	UniqueRequestID uuid.UUID `json:"unique"`

	// идентификатор пользователя переданный клиентом (если был передан)
	Client string `json:"client"`

	// номер телефона пользователя, куда мы совершали звонок
	Phone uint64 `json:"phone"`

	// код авторизации
	Code uint64 `json:"code"`

	// ISO код страны пользователя
	CountryCode string `json:"country_code"`

	// static.ucaller.ru/flag/ru.svg, изображение флага страны пользователя
	CountryImage string `json:"country_image"`

	// информация по телефону
	PhoneInfo []PhoneInfo `json:"phone_info"`

	// сколько стоила эта авторизация клиенту
	Cost float32 `json:"cost"`

	Error string `json:"error"`
}

func (s *Service) GetInfo(uid ID) (r *ResponseInfo, err error) {
	s.RLock()
	if _, ok := s.Calls[uid]; !ok {
		return nil, fmt.Errorf("ucaller_id не найден или истек")
	}
	s.RUnlock()

	body, err := s.Get(fmt.Sprintf("/getInfo?service_id=%d&key=%s&ucaller_id=%d", s.id, s.secretKey, uid))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, r)
	if err != nil {
		return nil, err
	}

	if !r.Status {
		return nil, errors.New(r.Error)
	}
	return
}
