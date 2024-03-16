<p align="center">
    <img src="logo.png" width="450" height="500">
</p>

# kk - Golang boilerplate CLI tool
> [!TIP]
> Create, setup and extend a Golang project like a puzzle in seconds

## How to use it?
### Install

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
Paths may be different depending on the path given in the command arguments above.

## Example project
kk generates projects like this https://github.com/waler4ik/kk-example (chi router, REST endpoints).

## Testing locally
After running `kk init`, build and start your service with the commands below. Then check your endpoints using a tool of your choice.
```zsh
docker compose build && docker compose up -d
```
## Upcoming features (commands)
- [ ] Add websocket command
- [ ] Switch router command (e.g chi -> gorilla)
- [ ] Add kubernetes scripts command
- [ ] Add github scripts command
- [ ] Generate swagger/openapi specification command

## Similar approaches 
- https://github.com/hay-kot/scaffold It's a more general approach. It gives you the possibility to write and use your own templates.
