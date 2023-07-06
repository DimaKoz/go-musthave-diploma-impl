# go-musthave-diploma-tpl ![Golang](https://img.shields.io/badge/-Golang%20❤️-05122A?style=flat&logo=go&logoColor=white)&nbsp; [![codecov](https://codecov.io/gh/DimaKoz/go-musthave-diploma-impl/branch/master/graph/badge.svg?token=VHZ6CU8FP6)](https://codecov.io/gh/DimaKoz/go-musthave-diploma-impl) [![Go Report Card](https://goreportcard.com/badge/github.com/DimaKoz/go-musthave-diploma-impl)](https://goreportcard.com/report/github.com/DimaKoz/go-musthave-diploma-impl)

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.
