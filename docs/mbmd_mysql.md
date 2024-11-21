## mbmd and mysql

Let mbmd write sensor data to mysql

### Table

Create this table in your database:

```
create table readings (
    device varchar(100),
    measurement varchar(60),
    value float,
    tstamp int,
    description varchar(100),
    unit varchar(10),
    primary key (device, measurement, tstamp),
    index idx_tstamp (tstamp)
)
```

### MySQL parameters

Run mbmd with:

```
mbmd run -a ... -d ... --mysql-host={your-db-host:3306} --mysql-user={your-mysql-user} --mysql-password={mysql-user-password} --mysql-database={mysql-database-name}
```

for example:

```
mbmd run -a192.168.1.123:1502 -dSOLAREDGE:1.0 -dSOLAREDGE:1.1 -r2s --mysql-database=solaredge --mysql-host=127.0.0.1:3306 --mysql-user=mbmd --mysql-password=secret
```
