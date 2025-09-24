First we wanted to create a new user that only has the rights to see the code and folders related to the deployment. Not the whole server.

1. sudo useradd --system --home /srv/devopsuser --shell /usr/sbin/nologin devopsuser

Then setup deployment folders
2. sudo mkdir -p /srv/devopsuser/releases /var/lib/devopsuser /var/log/devopsuser /etc/devopsuser

Give READ access to the newly created user.
3. sudo chown -R devopsuser:devopsuser /srv/devopsuser /var/lib/devopsuser /var/log/devopsuser /etc/devopsuser

# Installing go on the server (This is to get the newest version)
4. sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go