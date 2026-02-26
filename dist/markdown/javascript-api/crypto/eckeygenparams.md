
# EcKeyGenParams

The `EcKeyGenParams` object represents the object that should be passed as the algorithm parameter into the generateKey operation when generating key pairs for ECDH or ECDSA algorithms.

## Properties

| Property   | Type     | Description                                                                                              |
| :--------- | :------- | :------------------------------------------------------------------------------------------------------- |
| name       | `string` | An algorithm name. Possible values are `ECDH` or `ECDSA`.                                                |
| namedCurve | `string` | A elliptic curve's name to use for key pair generation. Possible values are `P-256`, `P-384` or `P-521`. |

## Example

```javascript
export default async function () {
  const keyPair = await crypto.subtle.generateKey(
    {
      name: 'ECDSA',
      namedCurve: 'P-256',
    },
    true,
    ['sign', 'verify']
  );

  console.log(JSON.stringify(keyPair));
}
```

