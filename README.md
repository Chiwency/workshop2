## 1. Start Up

### Pre-request

+ Make
+ Ubuntu
+ Docker

We use docker compose to start up the lab environment.

Start up four container: `master1`, `master2`, `worker`, `zookeeper`.
```bash
make start-up
```
Check the status of the containers:
```bash
docker ps
```

The output should be like this:
```
CONTAINER ID   IMAGE                 COMMAND                  CREATED          STATUS          PORTS                                                   NAMES
93e4b1d54a7e   myapp/master:latest   "/app/bin/master"        39 seconds ago   Up 35 seconds                                                           master1       
ac11bbef0daf   myapp/master:latest   "/app/bin/master"        39 seconds ago   Up 35 seconds                                                           master2       
b274b83c3c53   zookeeper             "/docker-entrypoint.…"   39 seconds ago   Up 36 seconds   2181/tcp, 2888/tcp, 3888/tcp, 8080/tcp                  zookeeper     
e75c2c49fda1   myapp/worker:latest   "/app/bin/worker"        39 seconds ago   Up 36 seconds                                                           worker    
```
If there is any container is not running, you can restart it by:
```bash
make stop-and-remove-all
make start-up
```


## 2. REMOTE PROCEDURE CALLS
The unit test defined in `./master/rpc_test.go`  simulates a rpc call from master to worker.

We can test the  RPC by running: 

```bash
make rpc-call
```
The output should be like this:
```
docker-compose exec worker /bin/bash -c "go test -v ./master"
=== RUN   TestPerformTask
2024/05/01 08:22:42 We got: Hello, worker! gatech
--- PASS: TestPerformTask (0.00s)
PASS
ok      github.com/chiwency/workshop2/master    0.004s
```
We can see the test passed, and the master send "Hello, worker!" and received the response from worker.
```bash
2024/05/01 08:22:42 We got: Hello, worker! gatech
```


## 3. LEADER ELECTION
The master take election every 5 seconds to choose a leader. Once a leader is elected, it will do a rpc call ts address as the input, and print the log to show it's leadership.
We can see the log by running:
```bash
make log
```
or you can see the any of the container log by:
```bash
make log-master1
make log-master2
make log-worker
```
The output should be like this:
```
master1  | 2024/05/01 08:32:32 I am the follower
master1  | 2024/05/01 08:32:37 I am the follower
master1  | 2024/05/01 08:32:42 I am the follower
master1  | 2024/05/01 08:32:47 I am the follower

master2  | 2024/05/01 08:32:45 I am the leader
master2  | 2024/05/01 08:32:45 I am the leader
master2  | 2024/05/01 08:32:45 I am the leader
master2  | 2024/05/01 08:32:45 I am the leader
master2  | 2024/05/01 08:32:45 I am the leader

worker  | 2024/05/01 08:19:30 server listening at [::]:7070
worker  | 2024/05/01 08:19:31 Performing task: leader_election
worker  | 2024/05/01 08:22:42 Performing task: 1
```

## 4.Failover
If you find the master2 is the leader, you can stop it by:
```bash
make stop-master2
```
After a while, you can then see the log of master1:
```bash
make log-master1
```
The output should be like this:
```
master1  | 2024/05/01 08:32:32 I am the follower
master1  | 2024/05/01 08:32:37 I am the follower
master1  | 2024/05/01 08:32:42 I am the follower
master1  | 2024/05/01 08:32:47 I am the follower
master1  | 2024/05/01 08:32:52 I am the leader
master1  | 2024/05/01 08:32:52 I am the leader
master1  | 2024/05/01 08:32:52 I am the leader
master1  | 2024/05/01 08:32:52 I am the leader
master1  | 2024/05/01 08:32:52 I am the leader
master1  | 2024/05/01 08:32:52 I am the leader
```
You can see the master1 become the leader after master2 is stopped.

Then you can restart master2 by:
```bash
make start-master2
```
After a while, you can then see the log of master1:
```bash
make log-master1
```
Or if you find the master1 is the leader, vice versa.


## 5.Clean UP
To clean up the environment, you can run:
```bash
make clean
```