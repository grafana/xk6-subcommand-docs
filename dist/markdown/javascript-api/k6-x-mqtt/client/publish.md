
# Client.publish()

Publish a message synchronously.

## Signature

```javascript
client.publish(topic, payload, options)
```

## Parameters

| Parameter | Type | Description |
| :-------- | :--- | :---------- |
| topic | string | Topic to publish to |
| payload | string \| ArrayBuffer | Message payload (string or binary data) |
| options | object | Optional publish configuration |
| options.qos | number | Quality of Service level (0, 1, or 2). Default: 0 |
| options.retain | boolean | Whether the message should be retained by the broker. Default: false |
| options.tags | object | Custom tags for metrics (key-value pairs) |
