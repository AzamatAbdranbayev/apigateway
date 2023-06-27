<h1>Запуск</h1>

> docker compose up -d


Обновить баланс (добавить или уменьшить)
> http://localhost:9090/api/auth/student/balance/add

Создать нового пользователя
> http://localhost:9090/api/auth/student/new

Свагер сервиса авторизации
> http://localhost:24000/docs/index.html#/Users/post_student_new


Примечания:
1) Префикс /api/auth обязателен.Так как стоит шина (api-gateway)
В зависимости от этого префикса данный сервер будет перенаправлять на нужный