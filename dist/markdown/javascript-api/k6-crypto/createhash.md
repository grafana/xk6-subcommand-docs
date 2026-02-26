
# createHash( algorithm )

> **Note:** A module with a better and standard API exists.
> 
> 
> The [crypto module](/docs/k6/v1.5.0/javascript-api/crypto/) partially implements the [WebCrypto API](https://www.w3.org/TR/WebCryptoAPI/), supporting more features than [k6/crypto](/docs/k6/v1.5.0/javascript-api/k6-crypto/).

Creates a hashing object that can then be fed with data repeatedly, and from which you can extract a hash digest whenever you want.

| Parameter | Type   | Description                                                                                                                                                       |
| --------- | ------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| algorithm | string | The name of the hashing algorithm you want to use. Can be any one of "md4", "md5", "sha1", "sha256", "sha384", "sha512", "sha512_224", "sha512_256", "ripemd160". |

### Returns

| Type   | Description                                                                                  |
| ------ | -------------------------------------------------------------------------------------------- |
| object | A Hasher object. |

### Example

```javascript
import crypto from 'k6/crypto';

export default function () {
  console.log(crypto.sha256('hello world!', 'hex'));
  const hasher = crypto.createHash('sha256');
  hasher.update('hello ');
  hasher.update('world!');
  console.log(hasher.digest('hex'));
}
```

The above script should result in the following being printed during execution:

```bash
INFO[0000] 7509e5bda0c762d2bac7f90d758b5b2263fa01ccbc542ab5e3df163be08e6ca9
INFO[0000] 7509e5bda0c762d2bac7f90d758b5b2263fa01ccbc542ab5e3df163be08e6ca9
```
