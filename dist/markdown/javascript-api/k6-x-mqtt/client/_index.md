
# Client

The `Client` class provides a high-level, event-driven interface for interacting with MQTT brokers. It supports both synchronous and asynchronous operations for connecting, subscribing, publishing, and unsubscribing.

## Constructor

```javascript
new Client(options)
```

Creates a new MQTT client instance.

### Parameters

| Parameter | Type | Description |
| :-------- | :--- | :---------- |
| options | object | Optional client configuration |
| options.client_id | string | Client identifier (must be unique per broker connection). If not provided, the broker assigns one. |
| options.username | string | Username for broker authentication |
| options.password | string | Password for broker authentication |
| options.credentials_provider | function | Function returning `{username, password}` for dynamic authentication |
| options.will | object | Last Will and Testament message configuration |
| options.will.topic | string | Topic for the will message |
| options.will.payload | string | Payload for the will message |
| options.will.qos | number | QoS level for the will message (0, 1, or 2) |
| options.will.retain | boolean | Whether the will message should be retained |
| options.tags | object | Custom tags for metrics (key-value pairs) |

## Properties

| Property | Type | Description |
| :------- | :--- | :---------- |
| connected | boolean | Read-only. Indicates if the client is currently connected to the broker. |

## QoS

Quality of Service enumeration for message delivery guarantees:

| Value | Name | Description |
| :---- | :--- | :---------- |
| 0 | QoS.AtMostOnce | Fire and forget. Message delivered at most once, no acknowledgment. |
| 1 | QoS.AtLeastOnce | Message delivered at least once, with acknowledgment. |
| 2 | QoS.ExactlyOnce | Message delivered exactly once, guaranteed and duplicate-free. |

## Methods

| Method | Description |
| :----- | :---------- |
| connect() | Connect to an MQTT broker |
| reconnect() | Reconnect to the broker using previous parameters |
| subscribe() | Subscribe to one or more topics |
| subscribeAsync() | Subscribe to topics asynchronously |
| unsubscribe() | Unsubscribe from topics |
| unsubscribeAsync() | Unsubscribe from topics asynchronously |
| publish() | Publish a message to a topic |
| publishAsync() | Publish a message asynchronously |
| on() | Register event handlers |
| end() | Disconnect from the broker |
| endAsync() | Disconnect from the broker asynchronously |

## Example

### Basic Usage

```javascript
import { Client } from "k6/x/mqtt";

export default function () {
  const client = new Client()

  client.on("connect", () => {
    console.log("Connected to MQTT broker")
    client.subscribe("greeting")

    client.publish("greeting", "Hello MQTT!")
  })

  client.on("message", (topic, message) => {
    const str = String.fromCharCode.apply(null, new Uint8Array(message))
    console.info("topic:", topic, "message:", str)
    client.end()
  })

  client.on("end", () => {
    console.log("Disconnected from MQTT broker")
  })

  client.connect(__ENV["MQTT_BROKER_ADDRESS"] || "mqtt://broker.emqx.io:1883")
}
```
