goutil
======

## os
Some( _most_ ) related function wrapper from linux environment

GetPsAuxCount = ```ps aux | grep xx | wc -l```

IsFreeMemoryLessThanMB = ```cat /proc/meminfo | grep MemFree```

## eff_bytesize
[ref](http://golang.org/doc/progs/eff_bytesize.go)

```fmt.Println(ByteSize(1), ByteSize(10), ByteSize(1000), ByteSize(1024))```

Output: _fmt.Println(ByteSize(1), ByteSize(10), ByteSize(1000), ByteSize(1024))_

```fmt.Println(2.5*MB, ByteSize(1e13), 1*MB==1024*KB)```

OUtput: _fmt.Println(2.5*MB, ByteSize(1e13), 1*MB==1024*KB)_

