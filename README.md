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
├── auth    # provides jwt service
│   ├── handler
│   └── proto
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

touch a `.env` in project root, like:

```shell script
# mysql
MYSQL_ROOT_PASSWORD=
MYSQL_USER=
MYSQL_PASSWORD=
MYSQL_DATABASE=
MYSQL_URL=

# redis
REDIS_PASSWORD=
REDIS_URL=

# app
APP_SECRET=
APP_UPLOAD_SECRET=

# weed
WEED_MASTER=
```

then run `./build.sh -a` to build go binaries, `./deploy.sh -a` to deploy it by docker-compose
