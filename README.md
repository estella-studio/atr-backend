# leon-backend

## Docker

Deploy with Docker

- Copy and modify `.env`
`cp .env.example .env; vim .env`

|Variable|Value|
|:---|:---|
|`APP_PORT`|The backend server will run on this port (make sure to not use well-known port (0 - 1023))|
|`DB_NAME`|Database name|
|`DB_USERNAME`|Database user|
|`DB_PASSWORD`|Database user password|
|`DB_HOST`|Database host|
|`DB_PORT`|Database port|
|`JWT_SECRET_KEY`|JWT secret key|
|`JWT_EXPIRED_DAYS`|JWT expiration (in days)|

- Build Docker Image
`docker build -t leon-backend:latest .`

- Run Docker Container
`docker run -d --name leon-backend --restart=always --network=host leon-backend:latest`

- Run MySQL Container
`docker run -d --name mysql -v mysql-volume:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=<your root password> -p 3306:3306 --restart=always mysql`

> Check server logs with `docker logs leon-backend`

## API Docs

|Request|Route Handler|Function|Note|
|:---|:---|:---|:---|
|`GET`|/users/info|Get user info|Requires Bearer Token|
|`POST`|/users/register|Register new user|-|
|`POST`|/users/login|Login|-|
|`PATCH`|/users/update|Update user info|Requires Bearer Token|

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
