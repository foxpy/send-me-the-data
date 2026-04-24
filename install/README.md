#### Files required to install SmtD as a systemd service

- `smtd.conf` - environment variables, includes a secret in a form of postgresql url
- `smtd.service` - systemd unit file

The idea is that SmtD will run under its own `smtd` user which will store data in `/var/lib/smtd` owned by `smtd:smtd` with umask `0027`.

Create user: `useradd -r -s /usr/bin/nologin smtd`.
TODO: add `tmpfiles.d entry` for user `smtd` with `/var/lib/smtd` for its data.
Edit `pg_hba.conf`: `local smtd smtd peer`. Also disable all other authentication options which allow anyone else to access database `smtd`.

Now the database connection should be possible:
- from UNIX user `smtd`
- with UNIX socket located in `/run/postgres`
- without any password
