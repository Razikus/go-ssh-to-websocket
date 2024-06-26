# GO SSH TO WEBSOCKET

This project does exactly what it says

It's an proxy written in GO that allows to connect to SSH server over WebSocket - just from your browser

I wrote it for other project, but maybe will be useful for somebody


![GO SSH TO WEBSOCKET](https://raw.githubusercontent.com/Razikus/go-ssh-to-websocket/main/image.png)


## Features

- SSH connection via WebSockets
- Web-based terminal emulation using [xterm.js](https://xtermjs.org/)
- Environment variable configuration for SSH credentials (password based) and settings
- Dockerized

## Running in docker

```

docker run -p 8280:8280 -e SSH_USER=USER -e SSH_PASS=PASS -e SSH_HOST=HOST -e SSH_PORT=PORT --rm docker.io/razikus/sshtows:1.0.1

```

Go to http://localhost:8280 and you will see terminal

## Configuration

Available env variables

```
SSH_USER="your_username"
SSH_PASS="your_password"
SSH_HOST="ssh.example.com"
SSH_PORT="22"
MOUNT_HTML="true"
```

## Docker compose with basic auth and simple SSH container

- 13.05.2024 - just added simple SSH container to make ability to log in inside the system out of the box.

In docker-compose there is example how to setup basic proxy with basic auth 

In order to change credentials from default (admin : admin) user need to follow traefik tutorial
```
https://doc.traefik.io/traefik/middlewares/http/basicauth/
```

Here you can create htpasswd online (remember to escape $ with $$)

```
https://hostingcanada.org/htpasswd-generator/
```
