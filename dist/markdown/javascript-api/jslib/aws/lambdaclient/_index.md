
# LambdaClient

> **Caution:** In some cases, using this library&apos;s operations might impact performance and skew your test results.
> 
> 
> To ensure accurate results, consider executing these operations in the `setup` and `teardown` [lifecycle functions](/docs/k6/v1.5.0/using-k6/test-lifecycle). These functions run before and after the test run and have no impact on the test results.

`LambdaClient` interacts with [AWS Lambda](https://aws.amazon.com/lambda/). With it, you can invoke a Lambda function.

Both the dedicated `lambda.js` jslib bundle and the all-encompassing `aws.js` bundle include the `LambdaClient`.

### Methods

| Function                                                                                                                  | Description                     |
| :------------------------------------------------------------------------------------------------------------------------ | :------------------------------ |
| invoke(name, payload, [options]) | Invokes an AWS Lambda function. |

### Throws

`LambdaClient` methods will throw errors in case of failure.

| Error                 | Condition                                                   |
| :-------------------- | :---------------------------------------------------------- |
| InvalidSignatureError | When invalid credentials were provided.                     |
| LambdaInvocationError | When AWS replied to the requested invocation with an error. |

### Examples

```javascript
import {
  AWSConfig,
  LambdaClient,
} from 'https://jslib.k6.io/aws//lambda.js';
import { check } from 'k6';

const awsConfig = new AWSConfig({
  region: __ENV.AWS_REGION,
  accessKeyId: __ENV.AWS_ACCESS_KEY_ID,
  secretAccessKey: __ENV.AWS_SECRET_ACCESS_KEY,
  sessionToken: __ENV.AWS_SESSION_TOKEN,
});

const lambdaClient = new LambdaClient(awsConfig);

export default async function () {
  const response = await lambdaClient.invoke('add-numbers', JSON.stringify({ x: 1, y: 2 }));

  check(response, {
    'status is 200': (r) => r.statusCode === 200,
    'payload is 3': (r) => r.payload === 3,
  });
}
```

