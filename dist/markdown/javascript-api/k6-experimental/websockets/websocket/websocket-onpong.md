
# WebSocket.onpong

A handler for a WebSocket connection `pong` event.
For multiple, simultaneous event handlers, use `WebSocket.addEventListener()`.

### Example

_A k6 script that initiates a WebSocket connection and setups a handler for the `pong` event._

```javascript
import { WebSocket } from 'k6/experimental/websockets';

export default function () {
  const ws = new WebSocket('ws://localhost:10000');

  ws.onpong = () => {
    console.log('A pong happened!');
    ws.close();
  };

  ws.onopen = () => {
    ws.ping();
  };
}
```

The preceding example uses a WebSocket echo server, which you can run with the following command:

```bash
docker run --detach --rm --name ws-echo-server -p 10000:8080 jmalloc/echo-server
```

