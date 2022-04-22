package ucaller

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	MaxCode            uint16 = 9999
	MaxLenForSecretKey int    = 32
	APIURL             string = "https://api.ucaller.ru/v1.0"
)

type ID uint64

type UCaller interface {
	// InitCall Данный метод позволяет инициализировать авторизацию для пользователя вашего приложения.
	InitCall(ic *InitCall) (*ResponseInitCall, error)

	// InitRepeat В случае, если ваш пользователь не получает звонок инициализированный методом initCall,
	// вы можете два раза и совершенно бесплатно инициализировать повторную авторизацию по uCaller ID,
	// который вы получаете в ответе метода initCall.
	// Повторную авторизацию можно запросить только в течение пяти минут с момента выполнения основной авторизации методом initCall.
	// Все данные, например `code` или `phone`, совпадают с теми же, которые были переданы в первом запросе initCall.
	InitRepeat(uid ID) (*ResponseInitRepeat, error)

	// GetInfo возвращает развернутую информацию по уже осуществленному uCaller ID
	GetInfo(uid ID) (*ResponseInfo, error)

	//  Этот метод возвращает информацию по остаточному балансу.
	GetBalance() (*ResponseBalance, error)

	// Этот метод возвращает информацию по сервису.
	GetService() (*ResponseService, error)
}

//go:generate mockgen -source=ucaller.go -destination=mocks/mock.go
type Requester interface {
	Get(url string) (body []byte, err error)
}

type InputData struct {
	// Секретный ключ вашего сервиса, Строка, 32 символа
	SecretKey string

	// ID вами созданного сервиса, Цифровое значение
	ID uint32

	// FreeRepeatTime время отведенное для совершения бесплатного повторного звонка
	FreeRepeatTime time.Duration

	// FreeRepeatNumber количество бесплатных повторов
	FreeRepeatNumber uint8
}

type Service struct {
	mux              sync.RWMutex
	secretKey        string
	id               uint32
	freeRepeatTime   time.Duration
	freeRepeatNumber uint8
	// cписок ucaller_id звонков. на первом месте ucaller_id из метода initCall,
	// каждый последующий потомок предидущего
	calls map[ID][]ID
	Requester
}

func (s *Service) Calls(id ID) []ID {
	s.mux.RLock()
	ids, ok := s.calls[id]
	s.mux.RUnlock()
	if !ok {
		return nil
	}
	return ids
}

func (s *Service) AddCalls(uID ID) {
	s.mux.RLock()
	s.calls[uID] = append(s.calls[uID], uID)
	s.mux.RUnlock()
}

func New(i *InputData, req Requester) (*Service, error) {
	if len([]rune(i.SecretKey)) == MaxLenForSecretKey {
		if req == nil {
			req = &Client{}
		}
		return &Service{
			secretKey:        i.SecretKey,
			id:               i.ID,
			freeRepeatTime:   i.FreeRepeatTime,
			freeRepeatNumber: i.FreeRepeatNumber,
			calls:            map[ID][]ID{},
			Requester:        req,
		}, nil
	}
	return nil, errors.New("неверная длина ключа, ожидается 32 символа")
}

type Client struct{}

func (c *Client) Get(url string) (body []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", APIURL, url), http.NoBody)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	body, err = io.ReadAll(req.Body)
	return
}
