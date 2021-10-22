# UCaller
адаптер для работы с сервисом https://ucaller.ru/

Оффициальная документация https://ucaller.ru/doc

## Example
```
package main

import "github.com/AleksandrMac/ucaller"

func main() {
    service, _ := ucaller.New(
        &ucaller.InputData{
            SecretKey:        "zxcvbasdfgqwertZXCVBASDFGQWERT23",
            ID:               103000,
            FreeRepeatTime:   5 * time.Minute,
            FreeRepeatNumber: 2,
        }, nil)

    responseInitCall, err = service.InitCall(
        &ucaller.InitCall{
            Phone:  &phone,
            Code:   &code,
            Client: &client,
            Unique: &unique,
        })
}
```
