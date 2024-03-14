# kk - Golang boilerplate command tool
> [!TIP]
> Create, setup and extend a Golang project in seconds

## Install kk

```zsh
go install github.com/waler4ik/kk@latest
```
After installation use `kk` command in your shell.

### Init command
Creates a Golang REST server project without endpoints.

```zsh
kk init 
```

### Add command
Creates and wires a REST resource endpoint

```zsh
kk add resource machines/data
```
After adding a resource implement the endpoint business logic in `internal/endpoints/machines/data/provider.go` and `internal/endpoints/machines/data/model.go`. 
Warning: Paths may be different depending on the path given in the command arguments above.

### Testing your service locally
Build and start your service with the commands below. Then check your endpoints using a tool of your choice.
```zsh
docker compose build && docker compose up -d
```

## Upcoming features (commands)
- [ ] Add websocket command
- [ ] Switch router command (e.g chi -> gorilla)
- [ ] Add kubernetes scripts command
- [ ] Add github scripts command
- [ ] Generate swagger/openapi specification command

