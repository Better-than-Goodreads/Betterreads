# Betterreads

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

## To do List:

[x] Publicar libro
[ ] Visualizar libro
[ ] Puntuar libro
[ ] Borrar puntaje
