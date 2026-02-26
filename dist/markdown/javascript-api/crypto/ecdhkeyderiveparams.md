
# EcdhKeyDeriveParams

The `EcdhKeyDeriveParams` represents the object that should be passed as the algorithm parameter into `deriveBits`, when using the ECDH algorithm.

ECDH is a secure communication method. Parties exchange public keys and use them with their private keys to generate a unique shared secret key.

## Properties

| Property  | Type                                                                                    | Description                          |
| :-------- | :-------------------------------------------------------------------------------------- | :----------------------------------- |
| name      | `string`                                                                                | An algorithm name. Should be `ECDH`. |
| publicKey | `CryptoKey` | Another party's public key.          |
