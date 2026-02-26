
# KMSClient

> **Caution:** In some cases, using this library&apos;s operations might impact performance and skew your test results.
> 
> 
> To ensure accurate results, consider executing these operations in the `setup` and `teardown` [lifecycle functions](/docs/k6/v1.5.0/using-k6/test-lifecycle). These functions run before and after the test run and have no impact on the test results.

`KMSClient` interacts with the AWS Key Management Service.

With it, the user can list all the Key Management Service keys in the caller's AWS account and region. They can also generate symmetric data keys to use outside of AWS Key Management Service.

Both the dedicated `kms.js` jslib bundle and the all-encompassing `aws.js` bundle include the `KMSClient`.

### Methods

| Function                                                                                                       | Description                                                                          |
| :------------------------------------------------------------------------------------------------------------- | :----------------------------------------------------------------------------------- |
| listKeys               | List the all the Key Management Service keys in the caller's AWS account and region. |
| generateDataKey | Generate a symmetric data key for use outside of the AWS Key Management Service.     |

### Throws

`KMSClient` methods throw errors in case of failure.

| Error                   | Condition                                                 |
| :---------------------- | :-------------------------------------------------------- |
| `InvalidSignatureError` | when using invalid credentials                            |
| `KMSServiceError`       | when AWS replied to the requested operation with an error |

### Example

```javascript
import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.2/index.js';

import {
  AWSConfig,
  KMSClient,
} from 'https://jslib.k6.io/aws//kms.js';

const awsConfig = new AWSConfig({
  region: __ENV.AWS_REGION,
  accessKeyId: __ENV.AWS_ACCESS_KEY_ID,
  secretAccessKey: __ENV.AWS_SECRET_ACCESS_KEY,
});

const kms = new KMSClient(awsConfig);
const keyAlias = 'alias/k6-key';

export async function setup() {
  // Create a symmetric data key
  return {
    dataKey: await kms.generateDataKey(keyAlias, 32),
  };
}

export default async function (data) {
  // Use the data key to encrypt data
}

export function handleSummary(data) {
  return {
    'stdout': textSummary(data, { indent: ' ', enableColors: true }),
    './test-run.key': data.dataKey,
  };
}
```

