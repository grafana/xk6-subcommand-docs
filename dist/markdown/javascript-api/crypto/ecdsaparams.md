
# EcdsaParams

The `EcdsaParams` represents the object that should be passed as the algorithm parameter into `sign` or `verify` when using the ECDSA algorithm.

## Properties

| Property | Type     | Description                                                                                           |
| :------- | :------- | :---------------------------------------------------------------------------------------------------- |
| name     | `string` | An algorithm name. Should be `ECDSA`.                                                                 |
| hash     | `string` | An identifier for the digest algorithm to use. Possible values are `SHA-256`, `SHA-384` or `SHA-512`. |
