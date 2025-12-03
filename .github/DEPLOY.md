# Инструкция по настройке деплоя

## Необходимые Secrets в GitHub

Для работы workflow необходимо настроить следующие secrets в настройках репозитория GitHub:
`Settings` → `Secrets and variables` → `Actions` → `New repository secret`

### Обязательные secrets:

1. **SERVER_HOST** - IP адрес или домен сервера для деплоя
   - Пример: `192.168.1.100` или `example.com`

2. **SSH_PRIVATE_KEY** - Приватный SSH ключ для подключения к серверу
   - Сгенерируйте ключ: `ssh-keygen -t ed25519 -C "github-actions"`
   - Скопируйте содержимое приватного ключа (обычно `~/.ssh/id_ed25519`)
   - Добавьте публичный ключ на сервер: `ssh-copy-id root@SERVER_HOST`

3. **HTTP_SERVER_USER** - Имя пользователя для Basic Auth
   - Пример: `admin`

4. **HTTP_SERVER_PASSWORD** - Пароль для Basic Auth
   - Пример: `secure_password_123`

### Опциональные secrets:

5. **SSH_PORT** - Порт SSH (по умолчанию 22)
   - Если не указан, используется порт 22

## Настройка сервера

1. Убедитесь, что на сервере установлен Docker:
   ```bash
   curl -fsSL https://get.docker.com -o get-docker.sh
   sh get-docker.sh
   ```

2. Убедитесь, что Docker запущен:
   ```bash
   systemctl enable docker
   systemctl start docker
   ```

3. Настройте SSH доступ для root пользователя:
   - Добавьте публичный ключ в `~/.ssh/authorized_keys`
   - Убедитесь, что SSH сервис запущен

4. Откройте порт 8080 в firewall (если используется):
   ```bash
   # Для ufw
   ufw allow 8080/tcp
   
   # Для firewalld
   firewall-cmd --permanent --add-port=8080/tcp
   firewall-cmd --reload
   ```

## Запуск деплоя

Workflow автоматически запускается при:
- Push в ветки `main` или `master`
- Ручном запуске через `Actions` → `Deploy to Server` → `Run workflow`

## Проверка деплоя

После успешного деплоя проверьте:
```bash
# Проверка контейнера
docker ps | grep url-shortener

# Проверка логов
docker logs url-shortener

# Проверка доступности
curl http://SERVER_HOST:8080
```

