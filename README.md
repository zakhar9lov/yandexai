# Пакет для обращения к Yandex Foundation Models

## Пример использования

```go
package main

import (
	"fmt"
	"github.com/zakhar9lov/yandexai/yandexgpt"
)

func main() {
	key := "NFHEb2hf48HFhde4d7fFDSHFBSFNJejdmwk"
	url := "model.url"
	folderid := "kjgndfjklangkj"

	gpt := yandexgpt.NewConnection(key, url, folderid)

	resp, err := gpt.Promt("Переведи текст", "Hello, Alisa!")
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

```