## podunpacker

A tool to examine and extract files from .pod archives used in certain PlayStation 2 (PS2) games.

### Installing

[Install the go compiler for your operating system](http://golang.org/doc/install) and [configure your workspace environment](http://golang.org/doc/install#gopath).  Download podunpacker using the following command.

```ShellSession
[dogshoes@oxcart ~]# go get github.com/dogshoes/podunpacker
```

Finally, ensure that your $GOPATH/bin is in $PATH for ease of use.

```ShellSession
[dogshoes@oxcart ~]# export PATH=$PATH:$GOPATH/bin
```

### Use

```ShellSession
[dogshoes@oxcart scratch]# podunpacker -e /mnt/cdrom/LANGUAGE.POD 
Input is a POD file, version 5.
Found 6 files.
* world\de\global.txt (17 bytes / 17 bytes), u1: 00000000, u2: 50987748, u3: 268f53a9
* world\en\global.txt (142 bytes / 142 bytes), u1: 00000000, u2: 4f987748, u3: 9d4b08e4
* world\en\test_sound.txt (128 bytes / 128 bytes), u1: 00000000, u2: 4f987748, u3: 9a17181d
* world\es\global.txt (17 bytes / 17 bytes), u1: 00000000, u2: 4f987748, u3: a6134f07
* world\fr\global.txt (18 bytes / 18 bytes), u1: 00000000, u2: 50987748, u3: 090da5a2
* world\it\global.txt (18 bytes / 18 bytes), u1: 00000000, u2: 50987748, u3: 98869414 
```
