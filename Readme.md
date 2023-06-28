<h1>Запуск</h1>

> docker compose up -d


Обновить баланс (добавить или уменьшить)
> http://localhost:9090/api/auth/student/balance/add

Создать нового пользователя
> http://localhost:9090/api/auth/student/new

Свагер сервиса авторизации
> http://localhost:24000/docs/index.html
Свагер сервиса авторизации
> http://localhost:24001/docs/index.html


Получить решение (и автоматом записывается в баланс)
> http://localhost:9090/api/algosolver/task/solution

Изменить цену задачки
> http://localhost:9090/api/auth/student/new

Получить список задач по пользователю
> http://localhost:9090/api/auth/student/new


Примечания:
1) Префикс /api/auth обязателен.Так как стоит шина (api-gateway)
В зависимости от этого префикса данный сервер будет перенаправлять на нужный
2) Префикс /api/algosolver обязателен. Для решения задач
3) SoftWeather.postman_collection.json это коллекция из postman. Можно ее заимпортить к себе и через нее тестировать
4) Тесты не успел
5) Создал три таблицы:
   1) user -  хранение пользовательских данных
   2) tasks - хранение данных о задачках.Каждой задачке присвоил определенный тип.По которому через switch case решаем ту или иную задачку
   3) history - хранение истории использование сервиса
6) В каждой директории есть папка .github. Можно настроить CICD  через github actions