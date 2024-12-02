# API of orders service

1. `PUT` `/orders/assign?order_id=...&executor_id=...&zone_id=...`
   - `400`, json `{"message": "Missing parameters, please provide order_id, executor_id and zone_id"}`, если хотя бы один из параметров `order_id`, `executor_id` или `zone_id` отсутствует в запросе.
   - `400`, json `{"message": "AssignedOrder with provided orderId already exists"}`, если происходит попытка зарегистрировать назначенный заказ с `order_id`, уже известным сервису. Заметим, что если подобное обращение происходит вследствие дублирования одного и того же запроса, то результат эквивалентен единичному применению запроса, поэтому данная операция идемпотентна.
   - `200`, json `{"message": "Successfully created"}` в случае успешной регистрации назначенного заказа.
   - `500`, json `{"message": "Unknown error"}` в случае любой внутренней ошибки сервиса. В интересах безопасности содержание ошибки не раскрывается пользователям данной ручки.
2. `GET` `/order/acquire?executor_id=...`
   - `400`, json `{"message": "Missing parameters, please provide executor_id"}`, если в запросе отсутствует `executor_id`.
   - `200`, json `{"message": "Successfully acquired", "order_profile": {...}}` в случае успешного получения назначенного заказа. В `order_profile` будет лежать заказ в схеме, идентичной схеме в БД (см. ниже).
   - `200`, json `{"message": "There are no orders assigned to you"}` в случае, если каждый заказ, когда либо назначенный на исполнителя, был выполнен или отменён.
   - `500`, json `{"message": "Unknown error"}` в случае любой внутренней ошибки сервиса. В интересах безопасности содержание ошибки не раскрывается пользователям данной ручки.
3. `POST` `/order/cancel?order_id=...`
   - `400`, json `{"message": "Missing parameters, please provide order_id"}`, если в запросе отсутствует `order_id`.
   - `200`, json `{"message": "Successfully canceled"` в случае успешной отмены заказа.
   - `200`, json `{"message": "AssignedOrder with OrderId $$order_id does not exist or has already been canceled"}` в случае, если данный заказ уже отменён, или выполнен, или никогда не существовал.
   - `500`, json `{"message": "Unknown error"}` в случае любой внутренней ошибки сервиса. В интересах безопасности содержание ошибки не раскрывается пользователям данной ручки.
