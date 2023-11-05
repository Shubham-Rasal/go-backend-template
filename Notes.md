## Update schema

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