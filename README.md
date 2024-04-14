## avito-backend-trainee-assignment
### Установка

1. Убедитесь, что у вас установлен Docker и Docker Compose.
2. Склонируйте этот репозиторий на свой компьютер:

    ```bash
    git clone https://github.com/riddledwithknife/avito-backend-trainee-assignment.git
    ```

### Запуск приложения

1. Перейдите в корневую директорию репозитория:

    ```bash
    cd avito-backend-trainee-assignment
    ```
   
2. Запустите команду

    ```bash
   make
    ```
   
3. Приложение поднимется на порте 8080 в docker compose кластере

### Использование

После успешного запуска приложения вы можете использовать следующие эндпоинты для взаимодействия с баннерами:

- `GET /user_banner`: Получение баннера для пользователя.
- `GET /banner`: Получение всех баннеров с фильтрацией по фиче и/или тегу.
- `POST /banner`: Создание нового баннера.
- `PATCH /banner/{id}`: Обновление содержимого баннера.
- `DELETE /banner/{id}`: Удаление баннера.
- `GET /CreateInitUsers`: Создание user и admin пользователей.

### Postman
[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/34267819-79a51a99-a3cf-46d1-bc5a-0969cb23cdbe?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D34267819-79a51a99-a3cf-46d1-bc5a-0969cb23cdbe%26entityType%3Dcollection%26workspaceId%3Db0d7c48c-7145-4011-aae8-780fc948ff32)