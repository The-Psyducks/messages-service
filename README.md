# go-microservice-template



## Como usar este template

Para iniciar el proyecto de go, correr el siguiente comando:

```bash
cd server
go mod init ${PROJECTNAME}
go mod tidy
```
para reflejar los ultimos cambios:
```bash
docker compose build
```

para correr el servicio:

```bash

docker compose up service
```

para correr los tests:

```bash
docker compose up tests
```

