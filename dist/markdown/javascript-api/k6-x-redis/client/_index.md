
# Client

`Client` is a [Redis](https://redis.io) client to interact with a Redis server, sentinel, or cluster. It exposes a promise-based API, which users can interact with in an asynchronous manner.

Though the API intends to be thorough and extensive, it does not expose the whole Redis API. Instead, the intent is to expose Redis for use cases most appropriate to k6.

## Usage

### Single-node server

You can create a new `Client` instance that connects to a single Redis server by passing a URL string.
It must be in the format:

```
redis[s]://[[username][:password]@][host][:port][/db-number]
```

Here's an example of a URL string that connects to a Redis server running on localhost, on the default port (6379), and using the default database (0):

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client('redis://localhost:6379');
```

A client can also be instantiated using an options object to support more complex use cases, and for more flexibility:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client({
  socket: {
    host: 'localhost',
    port: 6379,
  },
  username: 'someusername',
  password: 'somepassword',
});
```

### TLS

You can configure a TLS connection in a couple of ways.

If the server has a certificate signed by a public Certificate Authority, you can use the `rediss` URL scheme:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client('rediss://example.com');
```

Otherwise, you can supply your own self-signed certificate in PEM format using the socket.tls object:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client({
  socket: {
    host: 'localhost',
    port: 6379,
    tls: {
      ca: [open('ca.crt')],
    },
  },
});
```

Note that for self-signed certificates, k6's insecureSkipTLSVerify option must be enabled (set to `true`).

#### TLS client authentication (mTLS)

You can also enable mTLS by setting two additional properties in the socket.tls object:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client({
  socket: {
    host: 'localhost',
    port: 6379,
    tls: {
      ca: [open('ca.crt')],
      cert: open('client.crt'), // client certificate
      key: open('client.key'), // client private key
    },
  },
});
```

### Cluster client

You can connect to a cluster of Redis servers by using the cluster configuration property, and passing 2 or more node URLs:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client({
  cluster: {
    // Cluster options
    maxRedirects: 3,
    readOnly: true,
    routeByLatency: true,
    routeRandomly: true,
    nodes: ['redis://host1:6379', 'redis://host2:6379'],
  },
});
```

Or the same as above, but passing socket objects to the nodes array instead of URLs:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client({
  cluster: {
    nodes: [
      {
        socket: {
          host: 'host1',
          port: 6379,
        },
      },
      {
        socket: {
          host: 'host2',
          port: 6379,
        },
      },
    ],
  },
});
```

### Sentinel client

A [Redis Sentinel](https://redis.io/docs/management/sentinel/) provides high availability features, as an alternative to a Redis cluster.

You can connect to a sentinel instance by setting additional options in the object passed to the `Client` constructor:

```javascript
import redis from 'k6/x/redis';

const client = new redis.Client({
  username: 'someusername',
  password: 'somepassword',
  socket: {
    host: 'localhost',
    port: 6379,
  },
  // Sentinel options
  masterName: 'masterhost',
  sentinelUsername: 'sentineluser',
  sentinelPassword: 'sentinelpass',
});
```

## Real world example

```javascript
import { check } from 'k6';
import http from 'k6/http';
import redis from 'k6/x/redis';
import exec from 'k6/execution';
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

export const options = {
  scenarios: {
    redisPerformance: {
      executor: 'shared-iterations',
      vus: 10,
      iterations: 200,
      exec: 'measureRedisPerformance',
    },
    usingRedisData: {
      executor: 'shared-iterations',
      vus: 10,
      iterations: 200,
      exec: 'measureUsingRedisData',
    },
  },
};

// Instantiate a new redis client
const redisClient = new redis.Client(`redis://localhost:6379`);

// Prepare an array of rating ids for later use
// in the context of the measureUsingRedisData function.
const ratingIDs = new Array(0, 1, 2, 3, 4, 5, 6, 7, 8, 9);

export async function measureRedisPerformance() {
  // VUs are executed in a parallel fashion,
  // thus, to ensure that parallel VUs are not
  // modifying the same key at the same time,
  // we use keys indexed by the VU id.
  const key = `foo-${exec.vu.idInTest}`;

  await redisClient.set(key, 1);
  await redisClient.incrBy(key, 10);
  const value = await redisClient.get(key);
  if (value !== '11') {
    throw new Error('foo should have been incremented to 11');
  }

  await redisClient.del(key);
  if ((await redisClient.exists(key)) !== 0) {
    throw new Error('foo should have been deleted');
  }
}

export async function setup() {
  await redisClient.sadd('rating_ids', ...ratingIDs);
}

export async function measureUsingRedisData() {
  // Pick a random rating id from the dedicated redis set,
  // we have filled in setup().
  const randomID = await redisClient.srandmember('rating_ids');
  const url = `https://quickpizza.grafana.com/api/ratings/${randomID}`;
  const res = await http.asyncRequest('GET', url, {
    headers: { Authorization: 'token abcdef0123456789' },
  });

  check(res, { 'status is 200': (r) => r.status === 200 });

  await redisClient.hincrby('k6_rating_fetched', url, 1);
}

export async function teardown() {
  await redisClient.del('rating_ids');
}

export function handleSummary(data) {
  redisClient
    .hgetall('k6_rating_fetched')
    .then((fetched) => Object.assign(data, { k6_rating_fetched: fetched }))
    .then((data) => redisClient.set(`k6_report_${Date.now()}`, JSON.stringify(data)))
    .then(() => redisClient.del('k6_rating_fetched'));

  return {
    stdout: textSummary(data, { indent: '  ', enableColors: true }),
  };
}
```

## API

### Key value methods

| Method                                                                                                                                  | Redis command                                        | Description                                                                 |
| :-------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------------- | :-------------------------------------------------------------------------- |
| `Client.set(key, value, expiration)` | **[SET](https://redis.io/commands/set)**             | Set `key` to hold `value`, with a time to live equal to `expiration`.       |
| `Client.get(key)`                    | **[GET](https://redis.io/commands/get)**             | Get the value of `key`.                                                     |
| `Client.getSet(key, value)`       | **[GETSET](https://redis.io/commands/getset)**       | Atomically sets `key` to `value` and returns the old value stored at `key`. |
| `Client.del(keys)`                   | **[DEL](https://redis.io/commands/del)**             | Removes the specified keys.                                                 |
| `Client.getDel(key)`              | **[GETDEL](https://redis.io/commands/getdel)**       | Get the value of `key` and delete the key.                                  |
| `Client.exists(keys)`             | **[EXISTS](https://redis.io/commands/exists)**       | Returns the number of `key` arguments that exist.                           |
| `Client.incr(key)`                  | **[INCR](https://redis.io/commands/incr)**           | Increments the number stored at `key` by one.                               |
| `Client.incrBy(key, increment)`   | **[INCRBY](https://redis.io/commands/incrby)**       | Increments the number stored at `key` by `increment`.                       |
| `Client.decr(key)`                  | **[DECR](https://redis.io/commands/decr)**           | Decrements the number stored at `key` by one.                               |
| `Client.decrBy(key, decrement)`   | **[DECRBY](https://redis.io/commands/decrby)**       | Decrements the number stored at `key` by `decrement`.                       |
| `Client.randomKey()`           | **[RANDOMKEY](https://redis.io/commands/randomkey)** | Returns a random key's value.                                               |
| `Client.mget(keys)`                 | **[MGET](https://redis.io/commands/mget)**           | Returns the values of all specified keys.                                   |
| `Client.expire(key, seconds)`     | **[EXPIRE](https://redis.io/commands/expire)**       | Sets a timeout on key, after which the key will automatically be deleted.   |
| `Client.ttl(key)`                    | **[TTL](https://redis.io/commands/ttl)**             | Returns the remaining time to live of a key that has a timeout.             |
| `Client.persist(key)`            | **[PERSIST](https://redis.io/commands/persist)**     | Removes the existing timeout on key.                                        |

### List methods

| Method                                                                                                                                  | Redis command                                  | Description                                                                     |
| :-------------------------------------------------------------------------------------------------------------------------------------- | :--------------------------------------------- | :------------------------------------------------------------------------------ |
| `Client.lpush(key, values)`        | **[LPSUH](https://redis.io/commands/lpush)**   | Inserts all the specified values at the head of the list stored at `key`.       |
| `Client.rpush(key, values)`        | **[RPUSH](https://redis.io/commands/rpush)**   | Inserts all the specified values at the tail of the list stored at `key`.       |
| `Client.lpop(key)`                  | **[LPOP](https://redis.io/commands/lpop)**     | Removes and returns the first element of the list stored at `key`.              |
| `Client.rpop(key)`                  | **[RPOP](https://redis.io/commands/rpop)**     | Removes and returns the last element of the list stored at `key`.               |
| `Client.lrange(key, start, stop)` | **[LRANGE](https://redis.io/commands/lrange)** | Returns the specified elements of the list stored at `key`.                     |
| `Client.lindex(key, start, stop)` | **[LINDEX](https://redis.io/commands/lindex)** | Returns the specified element of the list stored at `key`.                      |
| `Client.lset(key, index, element)`  | **[LSET](https://redis.io/commands/lset)**     | Sets the list element at `index` to `element`.                                  |
| `Client.lrem(key, count, value)`    | **[LREM](https://redis.io/commands/lrem)**     | Removes the first `count` occurrences of `value` from the list stored at `key`. |
| `Client.llen(key)`                  | **[LLEN](https://redis.io/commands/llen)**     | Returns the length of the list stored at `key`.                                 |

### Hash methods

| Method                                                                                                                                         | Redis command                                    | Description                                                                                          |
| :--------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------- | :--------------------------------------------------------------------------------------------------- |
| `Client.hset(key, field, value)`           | **[HSET](https://redis.io/commands/hset)**       | Sets the specified field in the hash stored at `key` to `value`.                                     |
| `Client.hsetnx(key, field, value)`       | **[HSETNX](https://redis.io/commands/hsetnx)**   | Sets the specified field in the hash stored at `key` to `value`, only if `field` does not yet exist. |
| `Client.hget(key, field)`                  | **[HGET](https://redis.io/commands/hget)**       | Returns the value associated with `field` in the hash stored at `key`.                               |
| `Client.hdel(key, fields)`                 | **[HDEL](https://redis.io/commands/hdel)**       | Deletes the specified fields from the hash stored at `key`.                                          |
| `Client.hgetall(key)`                   | **[HGETALL](https://redis.io/commands/hgetall)** | Returns all fields and values of the hash stored at `key`.                                           |
| `Client.hkeys(key)`                       | **[HKEYS](https://redis.io/commands/hkeys)**     | Returns all fields of the hash stored at `key`.                                                      |
| `Client.hvals(key)`                       | **[HVALS](https://redis.io/commands/hvals)**     | Returns all values of the hash stored at `key`.                                                      |
| `Client.hlen(key)`                         | **[HLEN](https://redis.io/commands/hlen)**       | Returns the number of fields in the hash stored at `key`.                                            |
| `Client.hincrby(key, field, increment)` | **[HINCRBY](https://redis.io/commands/hincrby)** | Increments the integer value of `field` in the hash stored at `key` by `increment`.                  |

### Set methods

| Method                                                                                                                                   | Redis command                                            | Description                                                              |
| :--------------------------------------------------------------------------------------------------------------------------------------- | :------------------------------------------------------- | :----------------------------------------------------------------------- |
| `Client.sadd(key, members)`          | **[SADD](https://redis.io/commands/sadd)**               | Adds the specified members to the set stored at `key`.                   |
| `Client.srem(key, members)`          | **[SREM](https://redis.io/commands/srem)**               | Removes the specified members from the set stored at `key`.              |
| `Client.sismember(key, member)` | **[SISMEMBER](https://redis.io/commands/sismember)**     | Returns if member is a member of the set stored at `key`.                |
| `Client.smembers(key)`           | **[SMEMBERS](https://redis.io/commands/smembers)**       | Returns all the members of the set values stored at `keys`.              |
| `Client.srandmember(key)`     | **[SRANDMEMBER](https://redis.io/commands/srandmember)** | Returns a random element from the set value stored at `key`.             |
| `Client.spop(key)`                   | **[SPOP](https://redis.io/commands/spop)**               | Removes and returns a random element from the set value stored at `key`. |

### Miscellaneous

| Method                                                                                                                                         | Description                         |
| :--------------------------------------------------------------------------------------------------------------------------------------------- | :---------------------------------- |
| `Client.sendCommand(command, args)` | Send a command to the Redis server. |

