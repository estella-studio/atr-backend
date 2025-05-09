# LEON Backend

## Deploy

- Copy and modify `.env`

```
cp .env.example .env; vim .env
```

|Variable|Value|
|:---|:---|
|`LIMITER_MAX`|Max number of recent connections during `LIMITER_EXPIRATION_MINUTE` before sending a 429 response|
|`LIMITER_EXPIRATION_MINUTE`|Time before resetting the `LIMITER_MAX` count|
|`APP_PORT`|The backend server will run on this port (make sure to not use well-known port (0 - 1023))|
|`DB_NAME`|Database name|
|`DB_USERNAME`|Database user|
|`DB_PASSWORD`|Database user password|
|`DB_HOST`|Database host|
|`DB_PORT`|Database port|
|`JWT_SECRET_KEY`|JWT secret key|
|`JWT_EXPIRED_DAYS`|JWT expiration (in days)|

### Local

- Build

```
go build app/main.go
```

- Run

```
./main
```

> Make sure MySQL is running

### Docker

- Build Docker Image

```
docker build -t leon-backend:latest .
```

- Run Docker Container

```
docker run -d --name leon-backend --restart=always --network=host leon-backend:latest
```

- Run MySQL Container

> [IMPORTANT]
> Change `<your root password>` with your MySQL root password

```
docker run -d --name mysql -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=<your root password> -p 3306:3306 --restart=always mysql
```

> Check server logs with `docker logs leon-backend`

## API Docs

> Append after `/api/v1/`

|Request|Route Handler|Function|Note|
|:---|:---|:---|:---|
|`GET`|/ping|Test server latency|Any request body will be ignored|
|`GET`|/users/info|Get user info|Requires Bearer Token|
|`GET`|/data/get|Get save data from save id|Requires Bearer Token|
|`GET`|/data/list|List save data|Requires Bearer Token|
|`GET`|/data/listpaged/?offset=`n`&limit=`n`|List save data (paged)|Requires Bearer Token|
|`POST`|/users/register|Register new user|-|
|`POST`|/users/login|Login|-|
|`POST`|/data/add|Upload / save data to database|Requires Bearer Token, `form-data` key must be equal to `data`. Only 1 data can be accepted per request|
|`PATCH`|/users/update|Update user info|Requires Bearer Token|
|`DELETE`|/users/delete|Soft delete user|Requires Bearer Token|

### Sample API Response

#### Get User Info `/users/info`

- Request Body

None

- Response Body

```json
{
    "message": "user authenticated",
    "payload": {
        "id": "c3aed2d1-a009-4e99-b05d-c26ac8fdeff8",
        "email": "user0@gmail.com",
        "username": "user0",
        "name": "ウツミアオバ",
        "created_at": "2025-05-04T14:15:52+07:00",
        "updated_at": "2025-05-04T17:24:06+07:00"
    },
    "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImMzYWVkMmQxLWEwMDktNGU5OS1iMDVkLWMyNmFjOGZkZWZmOCIsImV4cCI6MTc0ODk0NjQwOX0.IdG3w3aJpvBfeBYQWxNnLS27WdDnzm7_YOKpDZQ5_VJE1XqqMFDzfp5zUQwp0WHA6BIR9w3MGxOd0G3cqgXVOg"
}
```

#### Retrieve Save Data `/data/get`

- Request Body

|Key|Type|Min|Max|Required|
|:---|:---|:---|:---|:---|
|id|string|36|36|required|

- Response Body

```json
{
    "message": "retrieved save data",
    "payload": {
        "data": "iVBORw0KGgoAAAANSUhEUgAAAj....."
    }
}

```

#### List Save Data `/data/list`

- Request Body

None

- Response Body

```json
{
    "message": "retrieved save data list",
    "payload": [
        {
            "id": "961f080a-62b0-40c4-a266-4d41edb58b45",
            "created_at": "2025-05-09T07:48:51+07:00"
        },
        {
            "id": "ad91ae46-15f7-484c-bbf3-83e59d1cb9e6",
            "created_at": "2025-05-07T16:51:02+07:00"
        }
    ]
}
```

#### List Save Data Paged `/data/listpaged/?offset=n&limit=n`

- Request Body

None

- Response Body

```json
{
    "message": "retrieved save data list",
    "payload": [
        {
            "id": "0f28309e-6b27-40f8-8943-0ebaca4806f9",
            "created_at": "2025-05-07T11:46:32+07:00"
        }
    ]
}
```

#### Register `/users/register`

- Request Body

```json
{
    "email": "user1@gmail.com",
    "username": "user1",
    "password": "passwordTest"
    "name": "Syafa",
}
```

|Key|Type|Min|Max|Required|
|:---|:---|:---|:---|:---|
|email|string|-|128|optional|
|username|string|3|64|required|
|password|string|8|256|required|
|name|string|3|128|optional|

- Response Body

```json
{
    "message": "user registered",
    "payload": {
        "id": "dca0ba20-a4f1-42c2-87db-2ac087449ef1",
        "email": "user1@gmail.com",
        "username": "user1",
        "name": "Syafa",
        "created_at": "2025-05-04T20:03:26.68+07:00",
        "updated_at": "2025-05-04T20:03:26.68+07:00"
    }
}
```

#### Login `/users/login`

- Request Body

```json
{
    "username": "user0",
    "password": "passwordTest"
}
```

|Key|Type|Min|Max|Required|
|:---|:---|:---|:---|:---|
|username|string|3|64|required|
|password|string|8|256|required|

- Response Body

```json
{
    "message": "user authenticated",
    "payload": {
        "id": "c3aed2d1-a009-4e99-b05d-c26ac8fdeff8",
        "email": "user0@gmail.com",
        "username": "user0",
        "name": "ウツミアオバ",
        "created_at": "2025-05-04T14:15:52+07:00",
        "updated_at": "2025-05-04T17:24:06+07:00"
    },
    "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJJRCI6ImMzYWVkMmQxLWEwMDktNGU5OS1iMDVkLWMyNmFjOGZkZWZmOCIsImV4cCI6MTc0ODk0NjQwOX0.IdG3w3aJpvBfeBYQWxNnLS27WdDnzm7_YOKpDZQ5_VJE1XqqMFDzfp5zUQwp0WHA6BIR9w3MGxOd0G3cqgXVOg"
}
```

#### Upload Data `/data/add`

- Request Body (`form-data`)

|Key|Type|Value|Required|
|:---|:---|:---|:---|
|`data`|File|any file|required|

- Response Body

```json
{
    "message": "data saved",
    "payload": {
        "id": "961f080a-62b0-40c4-a266-4d41edb58b45",
        "user_id": "dca0ba20-a4f1-42c2-87db-2ac087449ef1",
        "created_at": "2025-05-09T07:48:51.159+07:00"
    }
}

```

#### Update User `/users/update`

- Request Body

```json
{
    "email": "test@gmail.com"
    "username": "user64"
    "password": "newPassword"
    "name": "Syafa",
}
```

|Key|Type|Min|Max|Required|
|:---|:---|:---|:---|:---|
|email|string|-|128|optional|
|username|string|3|64|optional|
|password|string|8|256|optional|
|name|string|3|128|optional|

- Response Body

```json
{
    "message": "user updated",
    "payload": {
        "id": "0c2a4992-17d6-47a8-b970-a183a033c125",
        "email": "test@gmail.com",
        "username": "user64",
        "name": "Syafa",
        "created_at": "2025-05-04T14:10:56+07:00",
        "updated_at": "2025-05-04T17:26:30+07:00"
    }
}
```

#### Soft Delete User `/users/delete`

- Request Body

None

- Response Body

None
