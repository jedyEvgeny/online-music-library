package main

import (
	"jedyEvgeny/online-music-library/internal/pkg/app"
	"log"
)

// @title						Онлайн-библиотека музыки
// @version					1.0
// @description				Проект Effective Mobile
// @contact.name				Евгений
// @contact.url				https://github.com/jedyEvgeny
// @contact.email				KEF1991@yandex.ru
// @license.name				MIT
// @license.url				http://opensource.org/licenses/MIT
// @host						localhost:8080
// @BasePath					/
// @accept						json
// @produce					json text/plain
// @schemes					http
// @externalDocs.description	"Резерв для дополнительного описания API"
// @externalDocs.url			https://t.me/+ZGac_D1V4wFjYzRi
// @x-name						{"environment": "production", "version": "1.0.0", "team": "backend"}
// @tag.name					items
// @tag.description			Операции с товарами
// @tag.docs.url				https://t.me/EvKly
// @tag.docs.description		Связь с автором в Телеграм
func main() {
	a := app.New()
	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
