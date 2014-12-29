podcastd
=========

Podcast server daemon written in golang for turning your media directories into
realtime podcast feeds. Allows you to download or stream your media content to all your
devices through regular podcast clients and also through any web browser.


### Features

- Configuration file allows including media files from unlimited directories
and their sub-directories.

- Web interface and podcast feeds are secured by basic HTTP auth (configurable
username and password).

- Media directories are monitored and podcast feeds are updated in realtime as
the contents of the directories change.

- Movies and TV show information is scraped from the web for correct title,
release date, summary, runtime, genres, thumbnails and even trailers.

- A separate movie trailer feed is available to quickly read descriptions and
watch trailers of all your movies to make deciding on what to watch as
frictionless as possible.

- A web interface that displays details information about all your media, movie
trailer links, and links to all the custom podcast feeds to copy and paste into
your podcast client of choice. Media can also be played directly through the
web browser through the web interface.

![Podcastd Web Interface](https://github.com/ryanss/podcastd/raw/master/images/web.png)

![Podcastd Feed on iPhone](https://github.com/ryanss/podcastd/raw/master/images/iphone.png)


### Development Notes (Mac OS X)

```bash
brew install go
echo "export GOROOT=/usr/local/opt/go/libexec" >> .bash_profile
echo "export GOPATH=~/.go" >> .bash_profile
echo "export PATH=${GOPATH}/bin:${PATH}" >> .bash_profile
go get -u github.com/pilu/fresh
cd /path/to/podcastd
fresh -c fresh.conf
```
