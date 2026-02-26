
# sha1( input, outputEncoding )

> **Note:** A module with a better and standard API exists.
> 
> 
> The [crypto module](/docs/k6/v1.5.0/javascript-api/crypto/) partially implements the [WebCrypto API](https://www.w3.org/TR/WebCryptoAPI/), supporting more features than [k6/crypto](/docs/k6/v1.5.0/javascript-api/k6-crypto/).

Use [sha1](https://golang.org/pkg/crypto/sha1/) to hash input data.

| Parameter      | Type                 | Description                                                                                                                |
| -------------- | -------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| input          | string / ArrayBuffer | The input string or `ArrayBuffer` object to hash.                                                                          |
| outputEncoding | string               | Describes the type of encoding to use for the hash value. Can be "base64", "base64url", "base64rawurl", "hex" or "binary". |

### Returns

| Type           | Description                                                                                                                                             |
| -------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------- |
| string / Array | The hash digest as string (for "base64", "base64url", "base64rawurl", "hex" `outputEncoding`) or raw array of integers (for "binary" `outputEncoding`). |

### Example

```javascript
import crypto from 'k6/crypto';

export default function () {
  let hash = crypto.sha1('hello world!', 'hex');
  console.log(hash);
  const binArray = [104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100, 33];
  hash = crypto.sha1(new Uint8Array(binArray).buffer, 'hex');
  console.log(hash);
}
```

The above script should result in the following being printed during execution:

```bash
INFO[0000] 430ce34d020724ed75a196dfc2ad67c77772d169
INFO[0000] 430ce34d020724ed75a196dfc2ad67c77772d169
```
