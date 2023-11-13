## Update sch

### Reputation

Now I wanted to add a `reputation` field to the users tables.

```sql
ALTER TABLE users ADD COLUMN reputation INTEGER NOT NULL DEFAULT 0;
```

Now to change the schema, we can create a new migration with the following command.

```bash
./migrate create -ext sql -dir db/migration -seq add_reputation
```


After running the command, two files will be created in the `db/migration` directory. The up file will contain the commands to add the column and the down file will contain the commands to remove the column.

```sql
-- +migrate Up
ALTER TABLE users ADD COLUMN reputation INTEGER NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE users DROP COLUMN reputation;
```

### Likes 

Now I wanted to add a `likes` field to the posts tables.

```sql

ALTER TABLE posts ADD COLUMN likes INTEGER NOT NULL DEFAULT 0;
```

Now to change the schema, we can create a new migration with the following command.

```bash
./migrate create -ext sql -dir db/migration -seq add_likes
```

After running the command, two files will be created in the `db/migration` directory. The up file will contain the commands to add the column and the down file will contain the commands to remove the column.

```sql

-- +migrate Up
ALTER TABLE posts ADD COLUMN likes INTEGER NOT NULL DEFAULT 0;

-- +migrate Down
ALTER TABLE posts DROP COLUMN likes;
```



I think the best way to get into OS is to use a bunch of projects, packages and tools. When you use them, you will find bugs, missing features, and other things that you can fix.``


## Docker network

The container are running on different networks. So they can't communicate with each other.
We can connect them using the following ways
- using ip address of the container
- using docker network
- using docker compose

### Using ip address of the container

```bash

docker inspect <container_name> | grep "IPAddress"
```

For postgress, you can use this in the connection string.

### Using docker network

```bash

doker network create <network_name>
```

```bash 

docker network connect <network_name> <container_name>
```

```bash

docker network inspect <network_name>
```

