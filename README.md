![Go](https://github.com/mecitsemerci/go-todo-app/workflows/Go/badge.svg?branch=master)

# Go + Angular Todo APP Project Template

This repository is a todo sample go and angular web project built according to Clean Architecture.  

## Technologies
* Go Web Framework (gin-gonic)
* Containerize (Docker)
* Swagger (swaggo)
* Database
    * Mongodb (default)
    * Redis
* Dependency Injection (wire by google)
* Unit/Integration Tests (testify)
* Tracing (opentracing)
* WebUI (Angular 11)

### Web UI Preview
![GitHub Logo](https://github.com/mecitsemerci/blog/blob/master/src/images/angular_ui.gif?raw=true)

### Open API Doc Preview
![GitHub Logo](https://github.com/mecitsemerci/blog/blob/master/src/images/swagger_ui.jpg?raw=true)


## Layers and Dependencies

### `cmd` (application run)
Main application executive folder. Don't put a lot of code in the application directory.
The directory name for each application should match the name of the executable you want to have (e.g., /cmd/myapp).
It's common to have a small main function that imports and invokes the code from the /internal and /pkg directories and nothing else.

### `internal` (application codes)
Private application and library code. This is the code you don't want others importing in their applications or libraries.
* **core** includes application core files (domain objects, interfaces). It has no dependencies on another layer. 
* **pkg** includes external dependencies files and implementation of core interfaces.

### `test` (integration tests)
Application integration test folder.

### `web` (web ui)
Web application specific components: static web assets, server side templates and SPAs.

### `docs` (openapi docs)
open api (swagger) docs files. Swaggo generates automatically. 

    swag init -g ./cmd/api/main.go -o ./docs


## Usage

Open your terminal and clone the repository

    git clone https://github.com/mecitsemerci/go-todo-app.git

The application uses mongodb for default database so run makefile command

    make docker-mongo-start

This command builds all docker services so if it's ok check that application urls.  

Application | URL | Purpose
------------ | -------------| -------------
Angular UI | http://localhost:5000 | Todo APP Project
Swagger UI | http://localhost:8080/swagger/index.html | Todo API OpenAPI Docs
Jaeger UI | http://localhost:16686 | Opentracing Dashboard


By the way the application supports redis, if you use redis run that command

    make docker-redis-start

This command builds docker services so if it's ok check same application urls.

## Local Development
  ### Configuration
  The application uses environment variables. Environment variable names and values as follows by default. 
  ```
    # MONGO
    MONGO_URL=mongodb://127.0.0.1:27017
    MONGO_TODO_DB=TodoDb
    MONGO_CONNECTION_TIMEOUT=20
    MONGO_MAX_POOL_SIZE=10
    
    # REDIS
    REDIS_URL=127.0.0.1:6379
    REDIS_TODO_DB=0
    REDIS_CONNECTION_TIMEOUT=20
    REDIS_MAX_POOL_SIZE=10
    
    # JAEGER
    JAEGER_AGENT_HOST=localhost
    JAEGER_AGENT_PORT=6831
    JAEGER_SAMPLER_PARAM=1
    JAEGER_SAMPLER_TYPE=probabilistic
    JAEGER_SERVICE_NAME=go-todo-app
    JAEGER_DISABLED=false
  ```  

  ### Dependency Injection

  The project uses google wire for dependency injection. It comes from **wire_gen.go** for MongoDB. 
  Docker compose files generates automatically **wire_gen.go**. You can check under the `/internal/wired/wire_gen.go` but if you want to use redis you must regenerate for redis.
  
    make wire-redis
  
  This command generates **wire_gen.go** with redis provider. When you check the **wire_gen.go** file, you will see that it generates again for redis

  ```go
  // Injectors from redis.go:
  
  func InitializeTodoController() (handler.TodoHandler, error) {
      client, err := redisdb.ProvideRedisClient()
      if err != nil {
          return handler.TodoHandler{}, err
      }
      todoRepository := redisdb.ProvideTodoRepository(client)
      idGenerator := redisdb.ProvideIDGenerator()
      todoService := services.ProvideTodoService(todoRepository, idGenerator)
      todoHandler := handler.ProvideTodoHandler(todoService)
      return todoHandler, nil
  }
  ```

  If you want to use MongoDB again, you can run the command below.

    make wire-mongo
  
  You can observe the change in **wire_gen.go** after each change.
  
  ### Swagger
  
  You can use the code below to extract the open api documents to `/docs` folder.

    make swag

  ### Tests
  Existing tests are for demonstration purposes only

  Unit Test run command  

    make unit-test
  
  Integration Test run command
  
    make integration-test