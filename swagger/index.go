/*
Package swagger GmsTemp API.

<br/>
В каждом авторизованном запросе необходимо указать http-заголовок:
```
Authorization: token
```
token можно получить:
 + в ответе авторизации (/profile/auth)
 + в ответе регистрации (/profile/reg)
<br/>
### Константы
```
AppUrl     = "https://gms_temp.kz"

UsrTypeUndefined = 0
UsrTypeAdmin = 1
```
<br/>
### Работа с файлами
Файлы(фотки) из АПИ приходят в виде __path__. Чтобы отобразить ее на экране надо добавить префикс: https://gms_temp.kz/api/fs/__path__. Что бы залить фото на сервер надо использовать АПИ `/fs/`, подробно можно посмотреть ниже по документации.

Для большинства картин, которые возвращает API, действуют следующие параметры(query):
```
.../photo.jpg?w=200&h=200&m=fit|fill
```
Где:
+ w - желаемая ширина картинки
+ h - желаемая высота картинки
+ m - метод изменения размера. Допустимые значения:
	fit - картина полностью поместится в область (сохраняя соотношение)
	fill - картина может обрезаться, но покроет всю область
<br/>
### Websocket
websocket доступен по адресу `wss://gms_temp.kz/api/websocket?auth_token=<token>`. По вебсокету приходят такие же данные как и в пушах.


    Schemes: https, http
    Host: gms_temp.kz
    BasePath: /api
    Version: 1.0.0

    Consumes:
    - application/json

    Produces:
    - application/json

swagger:meta
*/
package swagger
