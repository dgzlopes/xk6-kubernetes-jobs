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