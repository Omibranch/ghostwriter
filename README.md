# ghostwriter ✍️

> Automatically types realistic-looking code in your terminal. For live streams, recordings, or just looking busy.

`ghostwriter` writes code snippets character by character, simulating human typing speed — including realistic typos, pauses before keywords, and natural hesitation.

---

## Install

```bash
git clone https://github.com/Omibranch/ghostwriter
cd ghostwriter
./install.sh
```

Requires Go 1.21+.

---

## Usage

```bash
ghostwriter                         # random language, medium speed
ghostwriter --lang go               # Go only
ghostwriter --lang python           # Python only
ghostwriter --lang typescript       # TypeScript only
ghostwriter --lang javascript       # JavaScript only
ghostwriter --lang rust             # Rust only
ghostwriter --lang bash             # Bash only
ghostwriter --speed slow            # ~80ms/char — hunt-and-peck
ghostwriter --speed medium          # ~45ms/char — competent dev (default)
ghostwriter --speed fast            # ~20ms/char — senior mode
ghostwriter --speed turbo           # ~8ms/char — suspiciously fast
ghostwriter --lang go --speed fast  # combine flags
```

Press `Ctrl+C` to stop.

---

## Features

- **6 languages**: Go, Python, TypeScript, JavaScript, Rust, Bash
- **4 speed profiles**: slow, medium, fast, turbo
- **Realistic typos**: ~2.5% chance of a wrong character, immediately corrected with backspace
- **Thinking pauses**: hesitation before comments (`//`, `/*`), function keywords
- **Snippet variety**: 8 real-world code snippets (HTTP servers, caches, event buses, deploy scripts)
- **No dependencies**: pure Go standard library

---

## Example output

```
# server.go

package main

import (
    "context"
    "fmt"
    "log"
    "net/http|"    ← typo corrected
    "net/http"
    "time"
)
...
```

---

## Use cases

- Live coding streams where you want filler content
- Demo recordings that need to look like coding is happening
- Making your screen look busy in an open-plan office
- Testing terminal rendering / font ligatures
- Scaring your rubber duck

---

## License

MIT

---

---

# ghostwriter ✍️

> Автоматически печатает реалистичный код в терминале. Для стримов, записей или просто чтобы выглядеть занятым.

`ghostwriter` вводит сниппеты кода посимвольно, имитируя скорость человека — включая опечатки, паузы перед ключевыми словами и естественные заминки.

---

## Установка

```bash
git clone https://github.com/Omibranch/ghostwriter
cd ghostwriter
./install.sh
```

Требуется Go 1.21+.

---

## Использование

```bash
ghostwriter                         # случайный язык, средняя скорость
ghostwriter --lang go               # только Go
ghostwriter --lang python           # только Python
ghostwriter --lang typescript       # только TypeScript
ghostwriter --lang javascript       # только JavaScript
ghostwriter --lang rust             # только Rust
ghostwriter --lang bash             # только Bash
ghostwriter --speed slow            # ~80мс/символ — двухпальцевый метод
ghostwriter --speed medium          # ~45мс/символ — нормальный разработчик (по умолчанию)
ghostwriter --speed fast            # ~20мс/символ — режим сеньора
ghostwriter --speed turbo           # ~8мс/символ — подозрительно быстро
ghostwriter --lang go --speed fast  # комбинация флагов
```

`Ctrl+C` для остановки.

---

## Возможности

- **6 языков**: Go, Python, TypeScript, JavaScript, Rust, Bash
- **4 скоростных профиля**: slow, medium, fast, turbo
- **Реалистичные опечатки**: ~2.5% шанс неверного символа, сразу исправляется backspace
- **Паузы «подумать»**: заминка перед комментариями (`//`, `/*`), ключевыми словами функций
- **Разнообразие сниппетов**: 8 реальных примеров кода (HTTP-серверы, кеши, event bus, деплой)
- **Без зависимостей**: чистая стандартная библиотека Go

---

## Сценарии использования

- Стримы с живым кодингом, где нужен контент для заполнения
- Запись демо, где нужно имитировать написание кода
- Занятый вид в open-space офисе
- Тестирование рендеринга шрифтов / лигатур в терминале

---

## Лицензия

MIT
