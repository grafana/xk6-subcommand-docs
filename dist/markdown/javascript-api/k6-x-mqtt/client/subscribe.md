
# Client.subscribe()

Subscribe to topics synchronously.

## Signature

```javascript
client.subscribe(topic, options)
```

## Parameters

| Parameter | Type | Description |
| :-------- | :--- | :---------- |
| topic | string \| string[] | Topic or array of topics to subscribe to |
| options | object | Optional subscription configuration |
| options.qos | number | Quality of Service level (0, 1, or 2). Default: 0 |
| options.tags | object | Custom tags for metrics (key-value pairs) |
