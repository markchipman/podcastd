podcastd
=========

Podcast server daemon written in golang for turning your media directories into
podcast feeds. Allows you to download or stream your media content to all your
devices through regular podcast clients.


### Development (Mac OS X)

```bash
brew install go
echo "export GOROOT=/usr/local/opt/go/libexec" >> .bash_profile
echo "export GOPATH=~/.go" >> .bash_profile
echo "export PATH=${GOPATH}/bin:${PATH}" >> .bash_profile
go get -u github.com/pilu/fresh
cd /path/to/podcastd
fresh
```
