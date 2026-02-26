
# AesKeyGenParams

The `AesKeyGenParams` object represents the object that should be passed as the algorithm parameter into the generateKey operation when generating an AES key.

## Properties

| Property | Type     | Description                                                                         |
| :------- | :------- | :---------------------------------------------------------------------------------- |
| name     | `string` | The name of the algorithm. Possible values are `AES-CBC`, `AES-CTR`, and `AES-GCM`. |
| length   | `number` | The length of the key in bits. Possible values are 128, 192 or 256.                 |

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
}
```

