
# KinesisClient

> **Caution:** In some cases, using this library&apos;s operations might impact performance and skew your test results.
> 
> 
> To ensure accurate results, consider executing these operations in the `setup` and `teardown` [lifecycle functions](/docs/k6/v1.5.0/using-k6/test-lifecycle). These functions run before and after the test run and have no impact on the test results.

`KinesisClient` interacts with the AWS Kinesis service.

With it, you can perform operations such as creating streams, putting records, listing streams, and reading records from streams. For a full list of supported operations, see [Methods](#methods).

Both the dedicated `kinesis.js` jslib bundle and the all-encompassing `aws.js` bundle include the `KinesisClient`.

### Methods

| Function                                                                                                                                                                | Description                                            |
| :---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------- |
| createStream(streamName, [options])                                     | Create a new Kinesis stream                            |
| deleteStream(streamName)                                                | Delete a Kinesis stream                                |
| listStreams([options])                                                   | List available Kinesis streams                         |
| putRecords(streamName, records)                                           | Put multiple records into a Kinesis stream             |
| getRecords(shardIterator, [options])                                      | Get records from a Kinesis stream shard                |
| listShards(streamName, [options])                                         | List shards in a Kinesis stream                        |
| getShardIterator(streamName, shardId, shardIteratorType, [options]) | Get a shard iterator for reading records from a stream |

### Throws

KinesisClient methods will throw errors in case of failure.

| Error                 | Condition                                                  |
| :-------------------- | :--------------------------------------------------------- |
| InvalidSignatureError | When invalid credentials are provided.                    |
| KinesisServiceError   | When AWS replies to the requested operation with an error. |

### Examples

```javascript
import { check } from 'k6';
import exec from 'k6/execution';

import {
  AWSConfig,
  KinesisClient,
} from 'https://jslib.k6.io/aws//kinesis.js';

const awsConfig = new AWSConfig({
  region: __ENV.AWS_REGION,
  accessKeyId: __ENV.AWS_ACCESS_KEY_ID,
  secretAccessKey: __ENV.AWS_SECRET_ACCESS_KEY,
});

const kinesis = new KinesisClient(awsConfig);
const testStreamName = 'test-stream';

export default async function () {
  // List available streams
  const streams = await kinesis.listStreams();
  console.log('Available streams:', streams.streamNames);

  // Check if our test stream exists
  if (!streams.streamNames.includes(testStreamName)) {
    // Create the stream if it doesn't exist
    await kinesis.createStream(testStreamName, { shardCount: 1 });
    console.log(`Created stream: ${testStreamName}`);
  }

  // Put some records into the stream
  const records = [
    {
      data: JSON.stringify({ message: 'Hello from k6!', timestamp: Date.now() }),
      partitionKey: 'test-partition-1',
    },
    {
      data: JSON.stringify({ message: 'Another message', timestamp: Date.now() }),
      partitionKey: 'test-partition-2',
    },
  ];

  const putResult = await kinesis.putRecords(testStreamName, records);
  console.log('Put records result:', putResult);

  // List shards in the stream
  const shards = await kinesis.listShards(testStreamName);
  console.log('Stream shards:', shards.shards);

  // Get a shard iterator for reading records
  if (shards.shards.length > 0) {
    const shardId = shards.shards[0].shardId;
    const shardIterator = await kinesis.getShardIterator(testStreamName, shardId, 'TRIM_HORIZON');

    // Get records from the shard
    const getResult = await kinesis.getRecords(shardIterator.shardIterator);
    console.log('Retrieved records:', getResult.records);
  }
}
```

#### Stream management

```javascript
import {
  AWSConfig,
  KinesisClient,
} from 'https://jslib.k6.io/aws//kinesis.js';

const awsConfig = new AWSConfig({
  region: __ENV.AWS_REGION,
  accessKeyId: __ENV.AWS_ACCESS_KEY_ID,
  secretAccessKey: __ENV.AWS_SECRET_ACCESS_KEY,
});

const kinesis = new KinesisClient(awsConfig);

export default async function () {
  const streamName = 'my-test-stream';

  // Create a stream with on-demand billing
  await kinesis.createStream(streamName, {
    streamModeDetails: {
      streamMode: 'ON_DEMAND',
    },
  });

  // List all streams
  const streams = await kinesis.listStreams();
  console.log('All streams:', streams.streamNames);

  // Clean up - delete the stream
  await kinesis.deleteStream(streamName);
  console.log(`Deleted stream: ${streamName}`);
}
```

