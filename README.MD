# Wallet Service
<h3 align="left">    
    <img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white">
    <img src="https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white">
    <img src="https://img.shields.io/badge/git-%23F05033.svg?style=for-the-badge&logo=git&logoColor=white">
    <img src="https://img.shields.io/badge/github-%23121011.svg?style=for-the-badge&logo=github&logoColor=white">
    <img src="https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white">
    <img src="https://img.shields.io/badge/Linux-FCC624?style=for-the-badge&logo=linux&logoColor=black">
</h3> 

## 1. Описание проекта

Этот учебный сервис предоставляет API для хранения музыкальных произведений, включая операции добавления произведений в базу данных, обновление информации о произведениях, получении информации и удалении произведений. Сервис написан на языке Go и использует PostgreSQL для хранения данных. Схема БД создаётся через миграции внутри кода.

## 2. Стек технологий

- Golang
- PostgreSQL & pgAdmin
- Docker
- Postman
- Git
- GitHub
- REST API
- Linux

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
Если PostgreSQL у вас не установлен, его можно поднять, используя docker и docker-compose из проекта.
    ```bash
   cd cd build
 ```
    ```bash
   docker-compose up -d
 ```
    ```bash
   cd ..
 ```

3.3. **Запустите проект в терминале**:

   ```bash
   go run cmd/app/main.go
 ```