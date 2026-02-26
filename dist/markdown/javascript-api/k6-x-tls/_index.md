
# k6/x/tls

> **Note:** This module is implemented as an official extension and is available natively in k6, requiring no additional installation or build steps. Refer to the [extensions documentation](/docs/k6/v1.5.0/extensions/explore) for available extensions and details.

The `k6/x/tls` module provides functionality for TLS certificate validation and inspection, allowing you to fetch certificate details and validate their properties directly from your k6 tests.

## Use cases

- Fetch TLS certificate information from any host
- Validate certificate expiration and properties
- Access certificate details (for example subject, issuer, fingerprint)

## API

### getCertificate( hostname )

Fetches TLS certificate information from a specified host.

The `tls.getCertificate` function retrieves TLS certificate details from a given hostname and returns a promise that resolves to a certificate object containing various properties such as expiration date, subject, issuer, and fingerprint information.

#### Parameters

| Parameter | Type | Description |
| --- | :--- | --- |
| hostname | string | The hostname to fetch the TLS certificate from. For example, "example.com" or "k6.io". |

#### Returns

A promise resolving to a certificate object containing the following properties:

| Property | Type | Description |
| --- | --- | --- |
| subject | object | The certificate subject information (pkixName object) |
| issuer | object | The certificate issuer information (pkixName object) |
| issued | number | The certificate issued timestamp in milliseconds since Unix epoch |
| expires | number | The certificate expiration timestamp in milliseconds since Unix epoch |
| fingerprint | string | The certificate fingerprint |

## Examples

### Check if a certificate is expired

```javascript
import tls from "k6/x/tls";
import { check } from "k6";

export default async function () {
  const cert = await tls.getCertificate("example.com");

  check(cert, {
    "certificate is not expired": (c) => c.expires > Date.now(),
  });

  console.log(`Certificate expires: ${new Date(cert.expires)}`);
}
```
