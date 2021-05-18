/*
Package swagger CTOgram API.

<br/>
В каждом запросе необходимо указать заголовки:
```
ch-platform: ios/android/site/msite
```
<br/>
В каждом авторизованном запросе необходимо указать http-заголовок:
```
Authorization: token
```
token можно получить в ответе авторизации (/profile/auth)
<br/>
### Константы
```
AppUrl     = "https://ctogram.kz"
AppOneLink = "http://onelink.to/ctogram"

ServiceShareUrlPattern = AppUrl + "/#/services/{id}"
OrderShareUrlPattern   = AppUrl + "/#/orders/{id}"

PlatformUndefined = 0
PlatformAndroid   = 1
PlatformIOS       = 2
PlatformSite      = 3
PlatformMSite     = 4
PlatformWhatsapp  = 5
PlatformOkauto    = 6
PlatformAitu      = 7

UsrTypeClient = 1
UsrTypeServer = 2

AdKindService = 1
AdKindOrder   = 2

PartTypeNew       = 1
PartTypeDetailing = 2
PartTypeUsed      = 4

ChatMsgTypeText        = 1
ChatMsgTypeImage       = 2
ChatMsgTypeVideo       = 3
ChatMsgTypeResponse    = 4
ChatMsgTypePriceChange = 5
ChatMsgTypeOrderChange = 6
ChatMsgTypePhoneCall   = 7

UsrAdmChatMsgTypeText         = 1
UsrAdmChatMsgTypeAnnouncement = 2

PushTypeRefreshNumbers    = "refresh-numbers"
PushTypeNewOrder          = "new-order"
PushTypeNewChat           = "new-chat"
PushTypeChatMessage       = "chat-message"
PushTypeChatMessagesRead  = "chat-messages-read"
PushTypeServiceComment    = "service-comment"
PushTypeCommentRequest    = "comment-request"
PushTypeAccountRefill     = "account-refill"
PushTypeAccountBonus      = "account-bonus"
PushTypeSmsBalanceAlarm   = "sms-balance-alarm"
PushTypeCarsChanged       = "usr_cars-changed"
PushTypeExpiredOrders     = "expired_orders"
PushTypeProAdPurchased    = "pro_ad_purchased"
PushTypeUsrAdmChatMessage = "usr_adm_chat_message"
PushTypeAutoPushContent   = "auto_push_content"

PermAll              = "*"
PermConfigE          = "-p-config-e"
PermDicCED           = "-p-dic-ced"
PermRoleCED          = "-p-group-ced"
PermAnnouncementCED  = "-p-announcement-ced"
PermEV               = "-p-ev"
PermCED              = "-p-ced"
PermBankTrCED        = "-p-bank-tr-ced"
PermAdSetPAF         = "-p-ad-set_paf"
PermProAdPriceCED    = "-p-pro_ad_price-ced"
PermSendNotification = "-p-send_notification"
PermSMSBalanceAlarmL = "-p-sms_balance_alarm-l"
PermFeedbackCED      = "-p-feedback-ced"
PermAllReports       = "-p-all_reports"
PermOptionE          = "-p-option-e"

ClosedAdTypeNotFound = 1
ClosedAdTypeFound    = 2
ClosedAdTypeAuto     = 3
```
<br/>
### Работа с файлами
Файлы(фотки) из АПИ приходят в виде __path__. Чтобы отобразить ее на экране надо добавить префикс: https://api.ctogram.kz/fs/__path__. Что бы залить фото на сервер надо использовать АПИ `/fs/`, подробно можно посмотреть ниже по документации.

Для большинства картин, которые возвращает API, действуют следующие параметры(query):
```
.../photo.jpg?w=200&h=200&m=fit/fill
```
Где:
+ w - желаемая ширина картинки
+ h - желаемая высота картинки
+ m - метод изменения размера. Допустимые значения:
	fit - картина полностью поместится в область (сохраняя соотношение)
	fill - картина может обрезаться, но покроет всю область
<br/>
### Websocket
websocket доступен по адресу `wss://api.ctogram.kz/websocket?auth_token=<token>`. По вебсокету приходят такие же данные как и в пушах.


    Schemes: https, http
    Host: api.ctogram.kz
    BasePath: /
    Version: 1.0.0

    Consumes:
    - application/json

    Produces:
    - application/json

swagger:meta
*/
package swagger
