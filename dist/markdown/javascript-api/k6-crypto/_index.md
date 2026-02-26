
# k6/crypto

> **Note:** A module with a better and standard API exists.
> 
> 
> The [crypto module](/docs/k6/v1.5.0/javascript-api/crypto/) partially implements the [WebCrypto API](https://www.w3.org/TR/WebCryptoAPI/), supporting more features than [k6/crypto](/docs/k6/v1.5.0/javascript-api/k6-crypto/).

The `k6/crypto` module provides common hashing functionality available in the GoLang [crypto](https://golang.org/pkg/crypto/) package.

| Function                                                                                                                | Description                                                                                                                  |
| ----------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| createHash(algorithm)                   | Create a Hasher object, allowing the user to add data to hash multiple times, and extract hash digests along the way.        |
| createHMAC(algorithm, secret)           | Create an HMAC hashing object, allowing the user to add data to hash multiple times, and extract hash digests along the way. |
| hmac(algorithm, secret, data, outputEncoding) | Use HMAC to sign an input string.                                                                                            |
| md4(input, outputEncoding)                     | Use MD4 to hash an input string.                                                                                             |
| md5(input, outputEncoding)                     | Use MD5 to hash an input string.                                                                                             |
| randomBytes(int)                       | Return an array with a number of cryptographically random bytes.                                                             |
| ripemd160(input, outputEncoding)         | Use RIPEMD-160 to hash an input string.                                                                                      |
| sha1(input, outputEncoding)                   | Use SHA-1 to hash an input string.                                                                                           |
| sha256(input, outputEncoding)               | Use SHA-256 to hash an input string.                                                                                         |
| sha384(input, outputEncoding)               | Use SHA-384 to hash an input string.                                                                                         |
| sha512(input, outputEncoding)               | Use SHA-512 to hash an input string.                                                                                         |
| sha512_224(input, outputEncoding)       | Use SHA-512/224 to hash an input string.                                                                                     |
| sha512_256(input, outputEncoding)       | Use SHA-512/256 to hash an input string.                                                                                     |

| Class                                                                              | Description                                                                                                                                                                                           |
| ---------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Hasher | Object returned by crypto.createHash(). It allows adding more data to be hashed and to extract digests along the way. |

