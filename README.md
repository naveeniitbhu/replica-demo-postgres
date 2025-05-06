# INTRO:

This repos contains details on how to setup postgres replica and replica and test it

## Steps

1. Install postgres
2. Check if role postgres exist. If not create it
3. psql -U $(whoami) -d postgres -c "\du"
4. CREATE ROLE postgres LOGIN SUPERUSER;
5. Open priamry postgresql.conf
   listen_addresses = '\*' # Listen on all interfaces (for local testing, 'localhost' is also fine)
   wal_level = replica # Enable logical decoding as well (required for streaming replication)
   max_wal_senders = 5 # Maximum concurrent connections for sending WAL
   wal_keep_size = 1GB # Minimum size of WAL files to keep for standby servers
   max_wal_timeout = 600s
6. Open pg_hba.conf
   host replication all 127.0.0.1/32 trust
7. mkdir /path/to/standby-data
8. ls -ld /Users/kaala/github/system-design/standby-data
9. chmod 0700 /Users/kaala/github/system-design/standby-data
10. pg_basebackup -h localhost -U postgres -D /path/to/standby-data -P -R
11. pg_ctl -D /Users/kaala/github/system-design/standby-data -l logfile -o "-p 5433 -c config_file=/Users/kaala/github/system-design/standby-data/postgresql.conf" start
    '-- On the primary server:
    CREATE DATABASE test_replication;
    \c test_replication;
    CREATE TABLE my_table (id SERIAL PRIMARY KEY, data TEXT);
    INSERT INTO my_table (data) VALUES ('Hello from primary!');-- On the standby server:\c test_replication;SELECT \* FROM my_table; '
12. Now you can connect via replica server and check if it is reflecting or not.
13. Node script is also there to test it.
14. In node script, execute first the primary connect function and then replica function.

## Important paths and commands:

1. /opt/homebrew/var/postgresql@17
2. /opt/homebrew/var/log
3. psql -h localhost -p 5432 -U postgres
4. psql -h localhost -p 5433 -U postgres
