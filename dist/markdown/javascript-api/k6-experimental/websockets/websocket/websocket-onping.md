
# WebSocket.onping

A handler for a WebSocket connection `ping` event.
For multiple, simultaneous event handlers, use `WebSocket.addEventListener()`.

### Example

_A k6 script that initiates a WebSocket connection and sets up a handler for the `ping` event._

```javascript
import { WebSocket } from 'k6/experimental/websockets';

export default function () {
  const ws = new WebSocket('wss://quickpizza.grafana.com/ws');

  ws.onping = () => {
    console.log('A ping happened!');
    ws.close();
  };

  ws.onclose = () => {
    console.log('WebSocket connection closed!');
  };

  ws.onopen = () => {
    ws.send(JSON.stringify({ event: 'SET_NAME', new_name: `Croc ${__VU}` }));
  };
  ws.onerror = (err) => {
    console.log(err);
  };
}
```

