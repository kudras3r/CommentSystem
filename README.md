# CommentSystem

Test tast in Ozon.

## Installation

```bash
git clone https://github.com/kudras3r/CommentSystem.git
cd CommentSystem/
```

In root create .env file:

```env
DB_HOST=comment-system-db # or edit docker-compose
DB_USER=ozon_keker
DB_PASS=1234
DB_NAME=comm_sys_db
DB_PORT=5432

LOG_LEVEL=DEBUG # INFO

SERVER_HOST=0.0.0.0
SERVER_POST=8080
```

### Docker

Run:
```bash
sudo docker-compose up --build
```

### Manually

In previuos step DB_HOST set in localhost / 0.0.0.0

```bash
cd cmd/
go run main.go --storage=inmemory
```
or
```bash
go run main.go --storage=db
```


## Roadmap

- [27.03.25] make basic structure | first server run
- [28.03.25] restructurize project | add pg db | add first resolvers | add migrations | add .env load func
- [29.03.25] scheme regenerate | bug fixes
- [30.03.25] add service level | add logger | inmemory rework
- [31.03.25] add logging | add docker | bug fixes



## Authors

- [@kudras3r](https://www.github.com/kudras3r)

