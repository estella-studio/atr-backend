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
