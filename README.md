# xk6-kubernetes-jobs

This is a [k6](https://github.com/loadimpact/k6) extension using the [xk6](https://github.com/k6io/xk6) system.

| :exclamation: This is a proof of concept, isn't supported by the k6 team, and may break in the future. USE AT YOUR OWN RISK! |
|------|

## Build

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Clone `xk6`:
  ```shell
  git clone https://github.com/k6io/xk6.git
  cd xk6
  ```

2. Build the binary:
  ```shell
  CGO_ENABLED=1 go run ./cmd/xk6/main.go build master \
    --with github.com/dgzlopes/xk6-kubernetes-jobs
  ```

## Example

```javascript
import { sleep } from 'k6';
import kubernetes from 'k6/x/kubernetes-jobs';

const client = new kubernetes.Client();

export default function () {
  client.create("pi-small","perl","perl -Mbignum=bpi -wle print bpi(20)")
  client.create("pi-big","perl","perl -Mbignum=bpi -wle print bpi(2000)")
  console.log(`Jobs: ${client.list()}`);
  client.deleteAll();
  sleep(2);
  console.log(`Jobs: ${client.list()}`);
}
```

Result output:

```
$ ./k6 run example.js

          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

  execution: local
     script: ../example.js
     output: -

  scenarios: (100.00%) 1 scenario, 1 max VUs, 10m30s max duration (incl. graceful stop):
           * default: 1 iterations for each of 1 VUs (maxDuration: 10m0s, gracefulStop: 30s)

INFO[0001] Jobs: pi-big,pi-small                         source=console
INFO[0003] Jobs:                                         source=console

running (00m03.0s), 0/1 VUs, 1 complete and 0 interrupted iterations
default ✓ [======================================] 1 VUs  00m03.0s/10m0s  1/1 iters, 1 per VU

     data_received........: 0 B 0 B/s
     data_sent............: 0 B 0 B/s
     iteration_duration...: avg=2.98s min=2.98s med=2.98s max=2.98s p(90)=2.98s p(95)=2.98s
     iterations...........: 1   0.333754/s
     vus..................: 1   min=1 max=1
     vus_max..............: 1   min=1 max=1
```
