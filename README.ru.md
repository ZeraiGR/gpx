# gpx

[English version](README.md)

`gpx` — CLI-утилита для управления **пресетами переменных окружения**
для Go-разработки (`GOPROXY`, `GOPRIVATE`, `GONOSUMDB`, `GOTOOLCHAIN` и любые другие).

Поддерживает:
- переключение **профилей** (несколько переменных за раз)
- временное применение через `eval "$(gpx use <profile>)"`
- постоянное применение через `gpx apply` (управляемый блок в rc-файле)
- редактирование конфига через CLI (`gpx profile ...`)
- подсветку **активного профиля** в `gpx list` (`*` — последний use/apply)

---

## Важное ограничение

CLI-программа **не может навсегда изменить окружение родительской оболочки**.

Поэтому:
- безопасный режим — печать `export ...` и применение через `eval`
- постоянный режим — явная команда `gpx apply`

---

## Установка

### Сборка из исходников

```bash
go test ./...
go build -o gpx ./cmd/gpx
install -m 0755 gpx /usr/local/bin/gpx
```

---

## Быстрый старт

### Создать конфиг

```bash
gpx init
```

### Список профилей (`*` — активный)

```bash
gpx list
```

### Временно применить профиль (рекомендуется)

```bash
eval "$(gpx use public)"
```

### Применить профиль постоянно

```bash
gpx apply public
source ~/.zshrc
```

---

## Команды

### gpx list

Показывает профили.  
Активный профиль (последний use/apply) отмечен `*`.

### gpx status

Показывает текущие значения переменных,
используемых в профилях.

### gpx use <profile>

Печатает команды `export ...`
(с безопасным shell-экранированием).

### gpx set KEY=VALUE [KEY=VALUE ...]

Разовая установка переменных без профиля.

### gpx unset KEY [KEY ...]

Печатает `unset KEY`.

### gpx diff <profile>

Показывает, что изменится относительно текущего окружения.

### gpx apply [flags] <profile>

Записывает управляемый блок в rc-файл:

```
# GPX_BEGIN
export GOPROXY='...'
export GOPRIVATE='...'
# GPX_END
```

Флаги:
- `--rc PATH` — явный путь к rc-файлу (имеет приоритет)
- `--shell zsh|bash` — выбор дефолтного rc
- `--dry-run` — показать результат без записи
- `--backup` — создать резервную копию (по умолчанию выключен)

**Контракт CLI:** флаги должны идти перед позиционными аргументами.

---

## Управление конфигом

### Управление профилями

```bash
gpx profile add corp
gpx profile rm corp
gpx profile rename corp avito
```

### Показать содержимое профиля

```bash
gpx profile show public
```

### Установить / удалить переменные в профиле

```bash
gpx profile set corp GOPROXY=https://proxy.corp.local,direct GOPRIVATE=github.com/mycorp/*
gpx profile unset corp GOPRIVATE GONOSUMDB
```

---

## Конфигурация

Путь по умолчанию:

```
~/.config/gpx/config.json
```

---

## Версия

```bash
gpx version
```

---

## Структура проекта

```
cmd/gpx            # вход CLI
internal/app       # сценарии и use-cases
internal/config    # load/save/validate
internal/envx      # env parsing, quoting, export/unset
internal/shell     # apply в rc-файлы (atomic replace)
internal/state     # активный профиль
```

---

## Разработка

```bash
go test ./...
```

---

## Лицензия

MIT
