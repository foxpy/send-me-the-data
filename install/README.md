#### Files required to install SmtD as a systemd service

- `smtd.conf` - environment variables, includes a secret in a form of postgresql url
- `smtd.service` - systemd unit file

The idea is that SmtD will run under its own `smtd` user which will store data in `/var/lib/smtd` owned by `smtd:smtd` with umask `0027`.

Create user: `useradd -r -s /usr/bin/nologin smtd`
