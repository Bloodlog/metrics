# go-musthave-metrics-tpl

Шаблон репозитория для трека «Сервер сбора метрик и алертинга».

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без префикса `https://`) для создания модуля.

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m main template https://github.com/Yandex-Practicum/go-musthave-metrics-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/main .github
```

Затем добавьте полученные изменения в свой репозиторий.

## Запуск автотестов

Для успешного запуска автотестов называйте ветки `iter<number>`, где `<number>` — порядковый номер инкремента. Например, в ветке с названием `iter4` запустятся автотесты для инкрементов с первого по четвёртый.

При мёрже ветки с инкрементом в основную ветку `main` будут запускаться все автотесты.

Подробнее про локальный и автоматический запуск читайте в [README автотестов](https://github.com/Yandex-Practicum/go-autotests).


![example workflow](https://github.com/Bloodlog/metrics/actions/workflows/.github/workflows/golangci-lint.yml/badge.svg?event=push)

![example workflow](https://github.com/Bloodlog/metrics/actions/workflows/.github/workflows/mertricstest.yml/badge.svg?event=push)

![example workflow](https://github.com/Bloodlog/metrics/actions/workflows/.github/workflows/statictest.yml/badge.svg?event=push)

* [pprof:](http://127.0.0.1:6060/debug/pprof/)
* [swagger:](http://127.0.0.1:8080/swagger/index.html#/)
* [Go doc](http://localhost:8081/)
* [Go doc dto](http://127.0.0.1:8081/pkg/metrics/internal/server/dto/)
* [Go doc api handlers](http://127.0.0.1:8081/pkg/metrics/internal/server/handlers/api/)
* [Go doc web handlers](http://localhost:8081/pkg/metrics/internal/server/handlers/web/)