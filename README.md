goutil
======

## os
Some( _most_ ) related function wrapper from linux environment

GetPsAuxCount = ```ps aux | grep xx | wc -l```

IsFreeMemoryLessThanMB = ```cat /proc/meminfo | grep MemFree```

## eff_bytesize
[ref](http://golang.org/doc/progs/eff_bytesize.go)

```fmt.Println(ByteSize(1), ByteSize(10), ByteSize(1000), ByteSize(1024))```

Output: _1.00B 10.00B 1000.00B 1.00KB_

```fmt.Println(2.5*MB, ByteSize(1e13), 1*MB==1024*KB)```

Output: _2.50MB 9.09TB true_

