# Онлайн-библиотека музыкальных произведений
<h3 align="left">    
    <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white">
    <img src="https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white">
    <img src="https://img.shields.io/badge/git-%23F05033.svg?style=for-the-badge&logo=git&logoColor=white">
    <img src="https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white">
    <img src="https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white">
    <img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
</h3> 

## 1. Описание проекта

<br>Этот учебный сервис предоставляет API для хранения музыкальных произведений, включая операции добавления информации о музыкальных произведениях в базу данных, обновление информации о произведениях, получении информации и удалении произведений. Сервис написан на языке Go и использует PostgreSQL для хранения данных. Схема БД создаётся через миграции внутри кода.</br>
<br>При запросе от клиента добавления песни, сервис обращается к стороннему серверу, в котором хранится информация о запрошенной песне, и сохраняет полученную от удалённого сервера информацию в свою БД.</br>
<br>Файл конфигурации находится в корне проекта: .env.</br>
<br>По-умолчанию настроен уровень логирования debug, кроме него доступен уровень info.</br>

## 2. Стек технологий

- Golang
- PostgreSQL & pgAdmin
- Docker
- Postman
- Git
- GitHub
- REST API
- Linux
- OpenAPI & Swagger
- RestAPI
- ООП

## 3. Установка

3.1. **Клонируйте репозиторий**:
   ```bash
   git clone https://github.com/jedyEvgeny/online-music-library.git
 ```

3.2. **Перейдите в каталог через терминал**:

   ```bash
   cd online-music-library
 ```
3.3. **Запустите проект в терминале**:
Если PostgreSQL у вас не установлен, его можно поднять, используя файл файл compose.yaml, ПО docker и docker-compose (должны быть установлены и запущены):
```bash
   cd build
 ```
```bash
   docker-compose up -d
 ```
```bash
   cd ..
 ```

3.4. **Запустите сервис в терминале**:

   ```bash
   go run cmd/app/main.go
 ```
 Сервис готов обрабатывать запросы клиентов.

 ## 4. Спецификация OpenAPI
 Для просмотра спецификации из веб-интерфейса и тестовых запросов, после запуска приложения, перейдите в браузере по пути:
 ```bash
 http://localhost:8080/swagger/
 ```

 Для обновления спецификации по стандарту OpenAPI, введите в терминал из домашнего каталога:

 ```bash
 swag init -d ./cmd/app/,./internal/app/endpoint/
 ```
 где флаг --dir value или -d value — Каталоги, которые вы хотите разобрать, разделенные запятыми; файл общей информации должен быть в первом (по-умолчанию: "./")

 ## 5. Тестовое удалённое хранилище
При разработке приложения не было доступа к удалённому серверу с базой данных песен. Реализован мок-сервер, моделирующий поведение удалённого сервера. Доступные для добавления песни:
```bash
{
 "group": "Muse",
 "song": "Supermassive Black Hole"
}
```
```bash
{
 "group": "Muse",
 "song": "test 1"
}
```
```bash
{
 "group": "Muse",
 "song": "test 2"
}
```
```bash
{
 "group": "Nirvana",
 "song": "Smells Like Teen Spirit"
}
```
```bash
{
 "group": "The Beatles",
 "song": "Hey Jude"
}
```
```bash
{
 "group": "The Beatles",
 "song": "test 1"
}
```
```bash
{
 "group": "The Beatles",
 "song": "test 2"
}
```