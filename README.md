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

1. Inside your `$GOPATH` (see http://golang.org for further info) :  
```
$ go get github.com/vil1/gerra/comment-added github.com/vil1/gerra/change-merged
```

2. Copy or move `$GOPATH/bin/comment-added` and `$GOPATH/bin/change-merged` inside the `$site_path/hooks` folder of your gerrit installation 

3. Add a `hooks.conf` configuration file (see the example provided)
