#### Postgresql instructions

Edit `pg_hba.conf`: `local smtd smtd peer`. Also disable all other authentication options which allow anyone else to access database `smtd`.

Now the database connection should be possible:
- from UNIX user `smtd`
- with UNIX socket located in `/run/postgres`
- without any password
