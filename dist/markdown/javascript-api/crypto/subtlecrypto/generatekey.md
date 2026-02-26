
# generateKey

The `generateKey()` generates a new cryptographic key and returns it as a CryptoKey object or a CryptoKeyPair object that can be used with the Web Crypto API.

## Usage

```
generateKey(algorithm, extractable, keyUsages)
```

## Parameters

| Name          | Type                                                       | Description                                                                                                                                                      |
| :------------ | :--------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `algorithm`   | a `string` or algorithm object with a single `name` string | The type of key to generate. It can be either a string with any of the currently supported algorithms as a value or any of the generation key parameter objects. |
| `extractable` | `boolean`                                                  | Whether the key can be exported using exportKey.                        |
| `keyUsages`   | `Array<string>`                                            | An array of strings describing what operations can be performed with the key. Key usages could vary depending on the algorithm.                                  |

### Supported algorithms

| AES-CBC                                                                                        | AES-CTR                                                                                        | AES-GCM                                                                                        | AES-KW | ECDH                                                                                                          | ECDSA                                                                                         | HMAC                                                                                                    | RSA-OAEP                                                                                                          | RSASSA-PKCS1-v1_5                                                                                                 | RSA-PSS                                                                                                           |
| :--------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------- | :----- | :------------------------------------------------------------------------------------------------------------ | :-------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------------ | :---------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------- | :---------------------------------------------------------------------------------------------------------------- |
| ✅ AesCbcParams | ✅ AesCtrParams | ✅ AesGcmParams | ❌     | ✅ EcdhKeyDeriveParams | ✅ EcdsaParams | ✅ HmacKeyGenParams | ✅ RsaHashedImportParams | ✅ RsaHashedImportParams | ✅ RsaHashedImportParams |

## Return Value

A `Promise` that resolves with the generated key as a CryptoKey object or a CryptoKeyPair object.

### Algorithm specific input

|                        | HMAC                                                                                                  | AES                                                                                                 | ECDH                                                                                              | ECDSA                                                                                             | RSA-OAEP                                                                                                        | RSASSA-PKCS1-v1_5                                                                                               | RSA-PSS                                                                                                         |
| :--------------------- | :---------------------------------------------------------------------------------------------------- | :-------------------------------------------------------------------------------------------------- | :------------------------------------------------------------------------------------------------ | :------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| Parameters type to use | `HmacKeyGenParams` | `AesKeyGenParams` | `EcKeyGenParams` | `EcKeyGenParams` | `RSAHashedKeyGenParams` | `RSAHashedKeyGenParams` | `RSAHashedKeyGenParams` |
| Possible key usages    | `sign`, `verify`                                                                                      | `encrypt`, `decrypt`                                                                                | `deriveKey`, `deriveBits`                                                                         | `sign`, `verify`                                                                                  | `encrypt`, `decrypt`                                                                                            | `sign`, `verify`                                                                                                | `sign`, `verify`                                                                                                |

## Throws

| Type          | Description                                                                                   |
| :------------ | :-------------------------------------------------------------------------------------------- |
| `SyntaxError` | Raised when the `keyUsages` parameter is empty, but the key is of type `secret` or `private`. |

## Example

```javascript
export default async function () {
  const key = await crypto.subtle.generateKey(
    {
      name: 'AES-CBC',
      length: 256,
    },
    true,
    ['encrypt', 'decrypt']
  );

  console.log(JSON.stringify(key));
}
```

