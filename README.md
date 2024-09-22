# Bank Management System

This project is the backend component of bank management system, implemented using Docker, Go, Postgresql and kafka. It provides a reliable and efficient backend system for library management system, more manageable ones.
High decopling also undirectly achived by utilizing asyncronous way.


## MVP Features

- User Management
- Bank Management

## Prerequisites

Make sure you have the following prerequisites installed:

- Docker (version >= 27.1.2)
- Docker Compose (version >= v1.29.2)

## Configuration

Modified Create ```app.yml``` file to in config  directory. ```app.yml``` supposed to be like this:

```
server:
  port: 8083
  read_timeout: 3
  write_timeout: 3

# for worker, its okay to not have server config

db:
  host: 172.17.0.1
  port: 5432
  user: postgres
  password: haris123
  db_name: bank_db
  ssl_mode: disable
  min_conn: 5
  max_conn: 500

kafka:
  broker: 172.17.0.1:9092
```

## How To Run



#### 1. Fun Way:


Thanks to docker-compose.

1. Clone repo

```
git clone git@github.com:Harisatul/bank-be.git
```

2. Suppose you are in root folder. running docker compose up command:

```
docker-compose up -d
```

3. Execute sql migration file ```sql_dum.sql``` on migration folder:

```
cd /migration

// excute sql_dump
```

#### 1. Also Fun Way:

build binary

1. Clone repo

```
git clone git@github.com:Harisatul/bank-be.git
```
2. Go to service folder

```
cd bank-backend
```

3. Build Binary

```
go build .
```

4. Execute Binary

```
./bank-backend serve-http
```
note: bank-backend has argument ```serve-http``` and bank-worker has argument ```serve```

5. Execute sql migration file ```sql_dump.sql``` on migration folder:

```
cd migration

// excute sql_dump.sql
```

### Data Migration:

```sql_dump.sql``` contain DDL sql to build db schema. since there are no pre-required data, you can try the system as soon as possible start with registering account.


## API USAGE

The backend component provides the API endpoints for Bank Management System. To interact with the backend, you
can use an API testing tool such as Postman.

Import postman collection [file]() to your Postman. 

## ERD

![img.png](img.png)
