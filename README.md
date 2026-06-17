Сервис для подачи заявок на майнкрафт сервер ElysiumSMP

---

Сборка в Docker
```bash
docker build . -t svc-applications:latest
docker run --rm svc-applications:latest
```

Для запуска требуется postgresql база данных, тестовую БД можно запустить так:
```bash
docker compose up -d
```

---

## Эндпоинты
### GET /applications/
Получает список всех заявок которые только были и информацию по ним

Пример данных в поле data ответа:
```jsonc
[
  {
    "id": "79b3a269-3f8a-4e3e-be86-27765680282b",
    "user_id": "ec43aada-449f-4dea-996b-ae5178725f82",
    "user_ip": "8.8.8.8", // Не отображается в /applications/me
    "age": 16,
    "about": "Я крутой человек вообще по жизни",
    "inviter": "Меня приглосил друг aktos",
    "ai_categories": { 
      "friends": ["AlexCraft", "Foks_f", "Mister"]
      "servers": ["BeeLand", "FW"]
      "categories": { // Произвольного формата json object
        "playstyle": "Roleplay",
        "experience": "Experienced (playing since 2011)"
      }
    },
    "ai_decision": "Accept", // Enum: 'Pending', 'Accept', 'Reject'
    "ai_answer": "Вы приняты! Ура поздравляем капец", // Пользователь это увидит
    "ai_comment": "Подробной по шагам объяснение почему такое решение было" // Не отображается в /applications/me
    "admin_decision": "Pending",
    "admin_id": null
  },
  // ... другие заявки
]
```

### GET /applications/:uuid
Возвращает тоже самое что и /applications/ но только по конкретному пользователю

### GET /applications/me
Возвращает тоже самое /applicatons/:uuid, для текущего пользователя, но с чуть урезанными данными

### POST /applicatons/:id/decision
Принимает:
```jsonc
{
  "admin_decision": "Accept" // Enum: 'Accept', 'Reject'
}
```

Возвращает пустой data

### POST /applications/create
Принимает:
```jsonc
{
  "age": 16,
  "about": "Билдер, но могу быть ленивым. Не умею строить редстоун фермы. Очень вспыльчив когда обижают. Меня из сервера Amigo знают как добрый и честный человек. Не умею пвпехатся из-за слабого пк, который тянет 30 FPS.",
  "inviter": "larik",
}
```

Возвращает вот такой поток данных, построчно (Не по API Guidelines):
```
Классификация заявки...
Выносим финальный вердикт...
{"action": "Accept", "answer": "Вы приняты! Ура поздравляем капец"}

```

---

Для конфигурации используются вот такие ENV переменные:
- DB_ADDRESS - адрес базы данных (По умолчанию "localhost:5432")
- DB_USER - юзер для подключения к базе данных (По умолчанию: "root")
- DB_PASS - пароль для доступа к базе данных (По умолчанию: "")
- DB_NAME - имя базы данных (По умолчанию "svc-applications")
- DB_ARGS - аргументы для подключения к базе данных (По умолчанию: "sslmode=disable")
- TRUSTED_PROXIES - разрешённые прокси в виде айпи или подсетей через запятую (По умолчанию: "127.0.0.1")
- GOOGLE_API_KEY - API ключ для Gemini (ОБЯЗАТЕЛЬНО)
