sudo sh -c 'echo "deb https://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt-get update
sudo apt-get -y install postgresql

# export GOROOT = /usr/local/go-1.20
# export GOPATH = $HOME/go/
# export GOROOT = $GOPATH/bin:$GOROOT/bin:$PATH
# export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
# export PATH=$PATH:/path/to/jadx-temp/bin