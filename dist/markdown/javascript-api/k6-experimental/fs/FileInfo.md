
# FileInfo

The `FileInfo` class represents information about a file.

## Properties

| Property | Type   | Description                    |
| :------- | :----- | :----------------------------- |
| name     | string | The name of the file.          |
| size     | number | The size of the file in bytes. |

## Example

```javascript
import { open, SeekMode } from 'k6/experimental/fs';

const file = await open('bonjour.txt');

export default async function () {
  // Retrieve information about the file
  const fileinfo = await file.stat();
  if (fileinfo.name != 'bonjour.txt') {
    throw new Error('Unexpected file name');
  }
}
```

