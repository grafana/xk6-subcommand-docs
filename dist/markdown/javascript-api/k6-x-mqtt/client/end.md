
# Client.end()

Disconnect from the broker synchronously. When the disconnection is complete, the `end` event is triggered.

## Signature

```javascript
client.end(options)
```

## Parameters

| Parameter | Type | Description |
| :-------- | :--- | :---------- |
| options | object | Optional disconnect configuration |
| options.tags | object | Custom tags for metrics (key-value pairs) |
