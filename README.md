# todo-list

> Это REST API pet-project, написанный на Go (вдохновлен zhashkevych и GolangLessons).

## Оглавление
- [Стек технологий](#стек-технологий)
- [Установка](#установка)
- [Реализованное](#реализованное)

## Стек технологий

- **Язык программирования**: Go
- **База данных**: PostgreSQL
- **Аутентификация**: JWT
- **Контейнеризация**: Docker
- **Документация**: в процессе разработки

## Установка

Инструкция по установке проекта на **Linux** (Ubuntu):

1. Клонируйте репозиторий:
```bash
git clone https://github.com/magneless/todo-app.git
```
2. Перейдите в директорию проекта:
```bash
cd todo-app
```
3. Установите зависимости:
```bash
go mod download
```
4. Развертывание БД:
   1. Установите Docker, если он еще не установлен:
   ```bash
   sudo apt update
   sudo apt install docker.io
   sudo systemctl start docker
   sudo systemctl enable docker
   ```
   2. Скачайте образ PostgreSQL:
   ```bash
   docker pull postgres
   ```
   3. Запустите контейнер PostgreSQL:
   ```bash
   docker run --name postgres -e POSTGRES_PASSWORD=qwerty -d -p 5436:5432 postgres
   ```
   4. Запустите утилиту migrate:
   ```bash
   migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up
   ```
5. Настройте конфигурационный файл **cmd/todo-app/config/local.yaml**, например:
```yaml
env: local
storage:
  host: localhost
  port: 5436
  username: postgres
  dbname: postgres
  sslmode: disable
http_server:
  address: localhost:8082
  timeout: 4s
  idle_timeout: 60s
```
6. Создайте в директории **todo-app** файл **.env**:
```conf
DB_PASSWORD=qwerty # пароль от базы данных
CONFIG_PATH=cmd/todo-app/config/local.yaml # путь к конфигу, если не перемещали, оставьте так
```
7. Запустите проект, если создали **.env**:
```bash
go run cmd/todo-app/main.go
```
   Или так, если не создавали **.env**:
```bash
DB_PASSWORD=qwerty CONFIG_PATH=cmd/todo-app/config/local.yaml go run cmd/todo-app/main.go
```

## Реализованное

- auth/sign-up
- auth/sign-in

