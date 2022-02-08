# banner-rotation
Проектная работа "Ротация баннеров"

Техническое задание:

https://github.com/OtusGolang/final_project/blob/master/02-banners-rotation.md

### Команды

- `make build` - сборка бинарника
- `make run` - запуск бинарника (DB in memory и без отправки в kafka)
- `make build-img` - сборка проекта в докере  (DB in memory и без отправки в kafka)
- `make run-img` - сборка проекта в докере и запуск  (DB in memory и без отправки в kafka)
- `make run-detached-img` - сборка проекта в докере и запуск в detached режиме (DB in memory и без отправки в kafka)
- `make stop-detached-img` - остановка и удаление контейнера с проектом
- `make test` - запуск unit-тестов для проекта
- `make lint` - запуск golangci-lint для проекта
- `make compose-up` - запуск docker-compose с проектом, postgres, kafka, zookeeper
- `make test-integration` - запуск интеграционных тестов
- `make compose-down` - остановка docker-compose с проектом, postgres, kafka, zookeeper

### Интеграционное тестирование

Последовательное выполнение шагов

1. Проект с DB in memory и без отправки в kafka
- `make run-detached-img`
- `make test-integration`
- `make stop-detached-img`

2. Проект с postgres, kafka, zookeeper
- `make compose-up`
- `make test-integration`
- `make compose-down`