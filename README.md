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

## To do List:

[x] Publicar libro
[ ] Visualizar libro
[ ] Puntuar libro
[ ] Borrar puntaje
