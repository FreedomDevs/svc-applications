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

Для конфигурации используются вот такие ENV переменные:
- DB_ADDRESS - адрес базы данных (По умолчанию "localhost:5432")
- DB_USER - юзер для подключения к базе данных (По умолчанию: "root")
- DB_PASS - пароль для доступа к базе данных (По умолчанию: "")
- DB_NAME - имя базы данных (По умолчанию "svc-applications")
- DB_ARGS - аргументы для подключения к базе данных (По умолчанию: "sslmode=disable")
- TRUSTED_PROXIES - разрешённые прокси в виде айпи или подсетей через запятую (По умолчанию: "127.0.0.1")
- GOOGLE_API_KEY - API ключ для Gemini (ОБЯЗАТЕЛЬНО)
