
# fs

> **Caution:** This is an experimental module.
> 
> 
> While we intend to keep experimental modules as stable as possible, we may need to introduce breaking changes. This could happen at future k6 releases until the module becomes fully stable and graduates as a k6 core module. For more information, refer to the [versioning and stability guarantees document](/docs/k6/v1.5.0/reference/versioning-and-stability-guarantees/#experimental).
> 
> 
> Experimental modules maintain a high level of stability and follow regular maintenance and security measures. Feel free to [open an issue](https://github.com/grafana/k6/issues) if you have any feedback or suggestions.

The k6 filesystem experimental module provides a memory-efficient way to handle file interactions within your test scripts. It currently offers support for opening files, reading their content, seeking through their content, and retrieving metadata about them.

### Memory efficiency

One of the key advantages of the filesystem module is its memory efficiency. Unlike the traditional open function, which loads a file multiple times into memory, the filesystem module reduces memory usage by loading the file as little as possible and sharing the same memory space between all VUs. This approach reduces the risk of encountering out-of-memory errors, especially in load tests involving large files.

### Notes on usage

An important consideration when using the filesystem module is its handling of external file modifications. Once k6 loads a file, it behaves like a "view" over its contents. If you modify the underlying file during a test, k6 will not reflect those changes in the loaded File instance.

## API Overview

The module exports functions and objects to interact with the file system:

| Function/Object                                                                                 | Description                                                                                          |
| ----------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| open         | Opens a file and returns a promise resolving to a `File` instance.                                   |
| File         | Represents a file with methods for reading, seeking, and obtaining file stats.                       |
| SeekMode | Enum for specifying the reference point for seek operations. Includes `Start`, `Current`, and `End`. |

## Example

```javascript
import { open, SeekMode } from 'k6/experimental/fs';

const file = await open('bonjour.txt');

export default async function () {
  // Seek to the beginning of the file
  await file.seek(0, SeekMode.Start);

  // About information about the file
  const fileinfo = await file.stat();
  if (fileinfo.name != 'bonjour.txt') {
    throw new Error('Unexpected file name');
  }

  const buffer = new Uint8Array(4);

  let totalBytesRead = 0;
  while (true) {
    // Read into the buffer
    const bytesRead = await file.read(buffer);
    if (bytesRead == null) {
      // EOF
      break;
    }

    // Do something useful with the content of the buffer
    totalBytesRead += bytesRead;

    // If bytesRead is less than the buffer size, we've read the whole file
    if (bytesRead < buffer.byteLength) {
      break;
    }
  }

  // Check that we read the expected number of bytes
  if (totalBytesRead != fileinfo.size) {
    throw new Error('Unexpected number of bytes read');
  }
}
```

