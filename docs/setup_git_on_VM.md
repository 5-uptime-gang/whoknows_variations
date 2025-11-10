sudo apt update
sudo apt install -y git

## Navigate to user folder and make repo folder
cd /srv/devopsuser
sudo mkdir -p /srv/devopsuser/repo
sudo chown -R devopsuser:devopsuser /srv/devopsuser

# Dont add a password (Just press enter)
sudo -u devopsuser -H ssh-keygen -t ed25519 -C "deploy@vm" -f /srv/devopsuser/.ssh/gh_deploy

# Copy the gets from gh_deploy to Github.com -> settings -> deploy keys
sudo -u devopsuser -H cat /srv/devopsuser/.ssh/gh_deploy.pub

# Create a config folder and add the text Host github.com....
sudo -u devopsuser -H bash -lc 'cat >> ~/.ssh/config <<EOF
Host github.com
  HostName github.com
  User git
  IdentityFile ~/.ssh/gh_deploy
  IdentitiesOnly yes
EOF'

## Clone the project using the devopsuser
sudo -u devopsuser -H git clone git@github.com:5-uptime-gang/whoknows_variations.git /srv/devopsuser/repo