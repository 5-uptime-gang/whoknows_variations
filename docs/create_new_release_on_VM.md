# Create a new folder with the timestap of the release time

STAMP=$(date +%F_%H-%M-%S)
RDIR="/srv/devopsuser/releases/$STAMP"
sudo -u devopsuser mkdir -p "$RDIR"

# GO to the folder containing main.go (This is subject to change, so check where it is)
cd /srv/devopsuser/repo/go_rewrite/src/backend
# Build the project to the newly created release folder (. is referencing the folder you are currently at, and uses the main.go file as standard (If main() is in another file, specify the filename))
sudo -u devopsuser go build -o '$RDIR/app' .

## TODO Challenges with the app not having files in the release folder