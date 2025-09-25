Deployment setup with a single user (adminuser)

Use your existing VM user
Instead of creating a new service account, we’ll just run and deploy everything as adminuser.

1. Create deployment folders

sudo mkdir -p /srv/adminuser/releases /var/lib/adminuser /var/log/adminuser /etc/adminuser

2. Give ownership to adminuser

sudo chown -R adminuser:adminuser /srv/adminuser /var/lib/adminuser /var/log/adminuser /etc/adminuser

3. What each folder is for

/srv/adminuser/releases → built binaries + static frontend files (timestamped per deploy).

/var/lib/adminuser → persistent data, e.g. whoknows.db.

/var/log/adminuser → app logs.

/etc/adminuser → config/env files.

# Installing go on the server (This is to get the newest version)
4. sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go