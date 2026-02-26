
# WebSocket

Creates a WebSocket instance for connection to a remote host.

The following events can close the connection:

- remote host close event.
- WebSocket.close().
- k6 VU interruption based on test configuration or CLI commands.

| Parameter | Type   | Description                                                                                                                                          |
| --------- | ------ | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| url       | string | The URL to which to connect (e.g. "ws://localhost:10000").                                                                                           |
| protocols | array  | Not yet implemented, reserved for the future use.                                                                                                    |
| params    | object | Params object containing additional request parameters. |

### Returns

| Type      | Description                      |
| --------- | -------------------------------- |
| WebSocket | An instance of WebSocket object. |

### Example

_A k6 script that initiates a WebSocket connection._

```javascript
import { WebSocket } from 'k6/experimental/websockets';

export default function () {
  const ws = new WebSocket('ws://localhost:10000');

  ws.onopen = () => {
    console.log('WebSocket connection established!');
    ws.close();
  };
}
```

The preceding example uses a WebSocket echo server, which you can run with the following command:

```bash
docker run --detach --rm --name ws-echo-server -p 10000:8080 jmalloc/echo-server
```

