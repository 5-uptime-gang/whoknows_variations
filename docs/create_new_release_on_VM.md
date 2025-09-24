# 1. Create a new folder with a timestamp
STAMP=$(date +%F_%H-%M-%S)
RDIR="/srv/adminuser/releases/$STAMP"
mkdir -p "$RDIR"

# 2. Go to the folder containing main.go
# In your project, that's in src
cd /srv/adminuser/repo/go_rewrite/src

# 3. Build the project into the release folder
# ⚠️ use double quotes so $RDIR expands
go build -o "$RDIR/app" .

# 4. Copy static files (your public/ dir) into the release folder
rsync -a ./public "$RDIR/"

# 6. Run the new release (test)
# For foreground run:
"/srv/adminuser/$RDIR/app"

# For background run:
nohup "$RDIR/app" > /srv/adminuser/app.log 2>&1 &


# Then you can see the process running:
ps aux | grep "app"

# To Check the logs for the app:
cat /srv/adminuser/app.log