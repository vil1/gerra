# Gerra 

Gerra is a set of hooks that connect Gerrit code review events to JIRA issues workflow transitions. 

As of version [2.5.2](http://gerrit-documentation.googlecode.com/svn/Documentation/2.5.2/config-hooks.html), Gerrit provides hooks for the following events : 

* patchset-created
* draft-published
* comment-added
* change-abandoned
* change-restored
* change-merged
* ref-updated
* cla-signed

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

2. Copy or move the`$GOPATH/bin/comment-added` and `$GOPATH/bin/change-merged` binaries inside the `$site_path/hooks` folder of your gerrit installation and make sure they're executable.

3. Add a `hooks.conf` configuration file 

### Configuration

The `hooks.conf` file _must be present_ for the hooks to work. It basically contains the configuration tidbits needed to connect to the JIRA REST API and Gerrit's command line : 

``` config
[jira]

 baseUrl = http://jira.example.com/rest/api/2
 user = vil1
 password = secret

[gerrit]

 host = localhost
 port = 29418

```

Since it contains non-encrypted sensitive data, you may make sure that this file is not publicly accessible, only the user that runs the gerrit process should be able to read this file.


## Gerrit & JIRA Versions

Gerra is currently tested against :

* Gerrit version **2.5.2**
* JIRA version **5.2.11**
