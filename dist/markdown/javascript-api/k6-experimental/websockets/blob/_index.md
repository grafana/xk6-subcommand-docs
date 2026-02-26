
# Blob

`Blob` is an interface that represents a blob, which is a file-like object of immutable, raw data; they can be read as text or binary data, or converted into a ReadableStream.

It's the type of the data received on WebSocket.onmessage when `WebSocket.binaryType` is set to `"blob"`. 

A `Blob` instance has the following methods/properties:

| Class/Property     | Description                                                                                                                                                                                 |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Blob.size          | The number of bytes of data contained within the `Blob`.                                                                                                                                     |
| Blob.type          | A string containing the MIME type, or an empty string if the type could not be determined.                                                                                                  |
| Blob.arrayBuffer() | Returns a `Promise` that resolves with the contents of the blob as binary data contained in an `ArrayBuffer`.                                                                               |
| Blob.bytes()       | Returns a `Promise` that resolves with a `Uint8Array` containing the contents of the blob as an array of bytes.                                                                             |
| Blob.slice()       | Returns a new `Blob` object which contains data from a subset of the blob on which it's called.                                                                                             |
| Blob.stream()      | Returns a ReadableStream which upon reading returns the data contained within the `Blob`. |
| Blob.text()        | Returns a `Promise` that resolves with a string containing the contents of the blob, interpreted as `UTF-8`.                                                                                |

