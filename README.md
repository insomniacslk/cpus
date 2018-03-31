# Turn CPUs on and off on Linux

```
$ ./cpus -h
Usage of ./cpus:
    status : print CPU online status. If no argument is passed, this is the default
    on     : turn CPUs on. Optionally pass a list of CPU # to turn on selectively
    off    : turn CPUs off. Optionally pass a list of CPU # to turn off selectively
```

```
$ sudo ./cpus off
Changing status for cpus: [0 1 2 3 4 5 6 7] to offline
CPU0 is online
CPU1 is offline
CPU2 is offline
CPU3 is offline
CPU4 is offline
CPU5 is offline
CPU6 is offline
CPU7 is offline
```

```
$ sudo ./cpus on 4 5 6 7
Changing status for cpus: [4 5 6 7] to online
CPU0 is online
CPU1 is offline
CPU2 is offline
CPU3 is offline
CPU4 is online
CPU5 is online
CPU6 is online
CPU7 is online
```

You'd better use Linux capabilities rather than running as root though.


## TODO

Implement a daemon mode that periodically changes the CPUs that are on and off
to avoid loading always the same ones.
