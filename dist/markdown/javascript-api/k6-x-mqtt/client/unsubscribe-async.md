
# Client.unsubscribeAsync()

Unsubscribe from topics asynchronously.

## Signature

```javascript
await client.unsubscribeAsync(topics, options)
```

## Parameters

| Parameter | Type | Description |
| :-------- | :--- | :---------- |
| topics | string \| string[] | Topic or array of topics to unsubscribe from |
| options | object | Optional unsubscribe configuration |
| options.tags | object | Custom tags for metrics (key-value pairs) |

## Returns

A promise that resolves when the unsubscription is complete.
