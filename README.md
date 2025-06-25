# HTTP server для обработки задач
Есть задачи, на обработку которых уходит 3-5 минут
Создан HTTP API, с помощью которго можно создавать\удалять залачи, получать результат их работы, статус задачи. Результат запроса статуса возвращает дату создания задачи, время работы и текущий статус выполнения

## API Reference

#### Create task

```http
  GET /create
```

#### Get task status

```http
  GET /status/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### Get task result

```http
  GET /data/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### Delete task

```http
  GET /delete/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |


## Tech Stack
* Go
* HTTP

