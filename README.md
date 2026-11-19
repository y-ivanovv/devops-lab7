# DevOps ЛР7

[![CI](https://github.com/y-ivanovv/devops-lab7/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/y-ivanovv/devops-lab7/actions/workflows/ci.yml)
[![Docker](https://github.com/y-ivanovv/devops-lab7/actions/workflows/docker.yml/badge.svg?branch=main)](https://github.com/y-ivanovv/devops-lab7/actions/workflows/docker.yml)

Вариант 2 «Docker + Docker Hub». В проект перенесены сервисы ЛР6. CI собирает их на Go 1.21 и 1.22, запускает тесты Todo и golangci-lint. Отчёт покрытия доступен как артефакт `coverage` в запуске CI.

Docker workflow запускается после успешного CI для push в main и публикует Todo как `yurgenius/todo-service` с тегами `main` и `sha-<7 символов SHA>`. Требуются секреты `DOCKERHUB_USERNAME` и `DOCKERHUB_TOKEN`.

```powershell
go build ./...
go test -coverprofile=coverage.out ./cmd/todo-service
go tool cover -html=coverage.out -o coverage.html
docker pull yurgenius/todo-service:main
docker run --rm -p 127.0.0.1:8080:8080 -p 127.0.0.1:9090:9090 yurgenius/todo-service:main
```

В другом терминале PowerShell 7 можно проверить HTTP:

```powershell
$task = '{"title":"Check container"}' | curl.exe -fsS http://127.0.0.1:8080/tasks -H 'Content-Type: application/json' --data-binary '@-' | ConvertFrom-Json
curl.exe -fsS -X PATCH "http://127.0.0.1:8080/tasks/$($task.id)/done"
curl.exe -fsS "http://127.0.0.1:8080/tasks/$($task.id)"
```
