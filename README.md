# fullbottle

> Online file service based on seaweedfs

## Based on

- Docker
- Etcd (using for go-micro service discovery)
- Go-micro
- gin (API framework)
- Mysql
- Redis
- Seaweedfs

## Features
- JWT authorization
- Chunk upload
- Resumable upload
- Flash upload (Chinese: 秒传)
- Range download

## Structure

```bash
.
├── api
│   ├── handler
│   ├── middleware
│   ├── route
│   └── util
├── bottle  # provides folder/file curd
│   ├── dao
│   ├── handler
│   ├── proto
│   ├── service
│   └── util
├── common  # common function, like db, redis
│   ├── db
│   ├── kv
│   └── log
├── config
├── fs      # docker seaweedfs volume
│   └── data
├── fullbottle-fe   # frontend submodule
│   ├── dist
│   ├── node_modules
│   ├── public
│   └── src
├── mysql   # docker mysql volume
│   └── data
├── nginx   # nginx conf
├── redis   # docker redis volume
│   └── data
├── share   # provides share service
│   ├── dao
│   ├── handler
│   ├── proto
│   └── service
├── upload  # provides upload service
│   ├── handler
│   └── proto
├── user    # provides user service
│   ├── dao
│   ├── handler
│   └── proto
├── util    # simple utils, like hash
└── weed    # weed client, supporting chunk upload
```

## Deploy

**Remember to modify docker-compose.yml to close dev port mapping**

Touch a `.env` in project root, like:

```shell script
# mysql
MYSQL_ROOT_PASSWORD=
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_DATABASE=
MYSQL_URL=mysql:3306

# redis
REDIS_PASSWORD=
REDIS_URL=redis:6379

# app
APP_SECRET=
APP_UPLOAD_SECRET=

# weed
WEED_MASTER=http://weed:9333
```

then run `./build.sh -a` to build go binaries, `./deploy.sh -a` to deploy it by docker-compose
