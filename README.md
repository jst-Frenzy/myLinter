# myLinter

Кастомный линтер для Go, проверяющий сообщения логов на соответствие правилам:
- Начинаются со строчной буквы
- Только английские символы
- Без специальных символов и эмодзи
- Без чувствительных данных

Интегрируется с golangci-lint (версия 2 и выше).

## Установка и использование

### Этап 1. Установка myLinter
В вашем проекте выполните:
```bash
go get github.com/jst-Frenzy/myLinter@latest
go mod tidy
```

### Этап 2. Настройка конфигурационных файлов
Создайте в корне проекта два файла:
.golangci.yml
```yaml
version: "2"

linters:
  enable:
    - MyLinter

  settings:
    custom:
      MyLinter:
        type: module
        path: github.com/jst-Frenzy/myLinter
        original-url: github.com/jst-Frenzy/myLinter
        settings:
          config: ./mylinter-config.json
```
и mylinter-config.json
```json
{
  "sensitive_words": ["jwt: ", "pwd: ", "refresh= "],
  "checks": {
    "special_chars": true,
    "capital_letter": false,
    "english_only": true,
    "sensitive_data": false
  }
}
```

### Этап 3. Запуск
Выполните golangci-lint run