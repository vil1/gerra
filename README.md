# Gerra 

Gerra is a set of hooks that connect Gerrit code review events to JIRA issues workflow transitions. 

As of version [2.2.1](http://gerrit.googlecode.com/svn/documentation/2.2.1/config-hooks.html), Gerrit provides hooks for the following events : 

* patchset-created
* comment-added
* change-abandoned
* change-restored
* change-merged
* ref-updated

Currently, only `comment-added` and `change-merged` are supported.


## Installation

### Prerequisites

You need to :

* have Go 1.1 installed and configured (see http://golang.org for further info) 
* make sure that the user running your gerrit process is allowed to connect to gerrit's sshd port
* have credentials to connect to your JIRA's REST API
 
### Installing hooks

1. Get and build the hooks :  
```
$ cd $GOPATH
$ go get github.com/vil1/gerra/comment-added github.com/vil1/gerra/change-merged
```

2. Copy or move `$GOPATH/bin/comment-added` and `$GOPATH/bin/change-merged` inside the `$site_path/hooks` folder of your gerrit installation 

3. Add a `hooks.conf` configuration file (see the example provided)

