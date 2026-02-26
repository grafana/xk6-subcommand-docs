
# dblclick(x, y[, options])

Mouse double clicks on the `x` and `y` coordinates. Similar to `mouse.click()`, but it simulates a double click.

| Parameter      | Type   | Description                                                                                                     |
| -------------- | ------ | --------------------------------------------------------------------------------------------------------------- |
| x              | number | The x-coordinate to double click on.                                                                            |
| y              | number | The y-coordinate to double click on.                                                                            |
| options        | object | Optional.                                                                                                       |
| options.button | string | The mouse button to double click. Possible values are `'left'`, `'right'`, and `'middle'`. Default is `'left'`. |
| options.delay  | number | The delay in milliseconds between `mousedown` and `mouseup` events. Default is `0`.                             |

### Returns

| Type            | Description                                                             |
| --------------- | ----------------------------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the mouse double click action is complete. |