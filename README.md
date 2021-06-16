
# Echo CMS
RESTful API service for cms application using [Echo framework](https://echo.labstack.com/). Admin page source code can be found [here](https://github.com/muhammadardie/admin-cms) and frontpage at [here](https://github.com/muhammadardie/react-cms) 

## Features

- API Documentation `Swagger (auto generate)`
- Authentication `Json Web Token`
- CRUD operations `MongoDB`
- Caching `Redis`
- Environment variables config
- Middlewares `CORS, Rate Limit, Logger, Recover, Custom, etc`

## System Requirements

- Golang
- MongoDB
- Redis

## Environment Variables

| **Key**          | **Description** 					 |
| :--------------- | :---------------------------------- |
| MONGODB_URL      | URL to connect to MongoDB instance  |
| MONGODB_NAME     | MongoDB database name				 |
| ACCESS_SECRET    | JWT key for access token			 |
| REFRESH_SECRET   | JWT key for refresh token			 |
| REDIS_ADDRESS    | URL to connect to Redis instance	 |
| REDIS_PASSWORD   | Redis Password                 	 |

## Demo
APP: https://echo-cms-app.herokuapp.com
API Documentation: https://echo-cms-app.herokuapp.com/swagger/index.html
