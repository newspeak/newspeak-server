go api server source
====================

this directory contains the sources of the `go api server`.

at the same time, this directory is an intellij ide project. see below for intellij ide install instructions.

command line compilation
------------------------

instructions to get started quickly, without intellij ide.

first, install [go](http://golang.org):

 * on osx using [homebrew](http://brew.sh): `brew install go`
 * on ubuntu (noninteractive): `sudo DEBIAN_FRONTEND='noninteractive' apt-get -o Dpkg::Options::='--force-confnew' -y install golang`

then:

    cd newspeak-server/shared/newspeak
    export GOPATH=$(pwd)
    go get github.com/fitstar/falcore
    go get github.com/Mistobaan/go-apns
    go get github.com/orfjackal/gospec/src/gospec
    go get github.com/peterbourgon/g2s
    go get github.com/bradfitz/gomemcache/memcache
    go build main
    ./bin/main

intellij
--------

intellij is used as ide because it provides good go language support. install instructions:

 * install intellij from [jetbrains.com](http://jetbrains.com)
 * install the ["golang.org support plugin"](http://plugins.jetbrains.com/plugin?pluginId=5047)

intellij run configuration: select menu -> run -> edit configuration

 * add "go application"
 * set "script" to "out/production/newspeak/go-bins/main/main" binary
 * select menu -> build -> rebuild project

go sdk configuration: select menu -> file -> project structure...

 * press "+" button -> and select "go sdk"
