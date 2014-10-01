btPodcast
=========

Self-hosted server daemon written in golang for managing and serving/streaming
media content to all your devices through regular podcast clients.


### Development (Mac OS X)

```bash
brew install go
echo "export GOROOT=/usr/local/opt/go/libexec" >> .bash_profile
echo "export GOPATH=~/.go" >> .bash_profile
echo "export PATH=${GOPATH}/bin:${PATH}" >> .bash_profile
go get github.com/pilu/fresh
cd /path/to/btpodcast
fresh
```
