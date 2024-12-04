# Betterreads API

Bettereads API es el backend de prueba tomando como inspiración a goodreads. La idea es crear una plataforma donde los usuarios puedan compartir sus opiniones sobre libros y puntuarlos. 

## Dependencias
El proyecto esta pensado para ser ejecutado en docker con una base de datos local. Por ende la unica dependencia necesaria es tener docker instalado en el sistema.

### Variables de entorno

Aun así se necesita tener el .env dentro de la carpeta src con las variables de entorno necesarias.

```shell
ENVIRONMENT=development
PORT=port
HOST=host
DATABASE_HOST=host
DATABASE_PORT=port
DATABASE_NAME=name
DATABASE_USER=user
DATABASE_PASSWORD=password
JWT_SECRET=any
JWT_DURATION_HOURS=1
```

Y otra .env dentro de /database: 

```shell
POSTGRES_USER=user
POSTGRES_PASSWORD=user123
POSTGRES_DB=db
```

> Notar que la informacion de las variables de entorno puede variar segun la configuración de la base de datos. Si notar que es importante que el user, password y database_name/postgres_db tienen que ser iguales en ambos archivos.

### Levantar el proyecto y la base de datos

Simplemente entonces levantamos con docker compose: 
```shell
docker compose up --build
```


## Documentación

La documentacion esta automatizada con swagger + swag en go.
Para generar la documentación se utiliza necesita instalar la CLI de Swag con :

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

Luego para generar la documentación se utiliza el comando dentro de src/:

```shell
swag init -g cmd/main.go
```

Para visualizar la documentación, teniendo el servidor corriendo simplemente acceda a :
[Link con puerto en 8080](http://localhost:8080/swagger/index.html#/)
link generico : http://localhost:PORT/swagger/index.html#/


