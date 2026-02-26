
# createHMAC( algorithm, secret )

> **Note:** A module with a better and standard API exists.
> 
> 
> The [crypto module](/docs/k6/v1.5.0/javascript-api/crypto/) partially implements the [WebCrypto API](https://www.w3.org/TR/WebCryptoAPI/), supporting more features than [k6/crypto](/docs/k6/v1.5.0/javascript-api/k6-crypto/).

Creates a HMAC hashing object that can then be fed with data repeatedly, and from which you can extract a signed hash digest whenever you want.

| Parameter |         Type         | Description                                                                                                                         |
| --------- | :------------------: | :---------------------------------------------------------------------------------------------------------------------------------- |
| algorithm |        string        | The hashing algorithm to use. One of `md4`, `md5`, `sha1`, `sha256`, `sha384`, `sha512`, `sha512_224`, `sha512_256` or `ripemd160`. |
| secret    | string / ArrayBuffer | A shared secret used to sign the data.                                                                                              |

### Returns

| Type   | Description                                                                                  |
| ------ | :------------------------------------------------------------------------------------------- |
| object | A Hasher object. |

### Example

```javascript
import crypto from 'k6/crypto';

export default function () {
  console.log(crypto.hmac('sha256', 'a secret', 'my data', 'hex'));
  const hasher = crypto.createHMAC('sha256', 'a secret');
  hasher.update('my ');
  hasher.update('data');
  console.log(hasher.digest('hex'));
}
```

The above script should result in the following being printed during execution:

```bash
INFO[0000] 82f669c8fde13aef6d6977257588dc4953dfac505428f8fd6b52e19cd96d7ea5
INFO[0000] 82f669c8fde13aef6d6977257588dc4953dfac505428f8fd6b52e19cd96d7ea5
```
