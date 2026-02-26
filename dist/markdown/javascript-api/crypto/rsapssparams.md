
# RsaPssParams

The `RsaPssParams` represents the object that should be passed as the algorithm parameter into `sign` or `verify` when using the RSA-PSS algorithm.

## Properties

| Property   | Type     | Description                                                                 |
| :--------- | :------- | :-------------------------------------------------------------------------- |
| name       | `string` | An algorithm name. Should be `RSA-PSS`.                                     |
| saltLength | `number` | A long integer representing the length of the random salt to use, in bytes. |

> **Caution:** Since under the hood we use Golang's SDK the salt length 0 is not supported. In that case the maximum possible salt length will be used.

