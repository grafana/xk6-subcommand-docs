
# Selection

Represents a set of nodes in a DOM tree.

Selections have a jQuery-compatible API, but with two caveats:

- CSS and screen layout are not processed, thus calls like css() and offset() are unavailable.
- DOM trees are read-only, you can't set attributes or otherwise modify nodes.

(Note that the read-only nature of the DOM trees is purely to avoid a maintenance burden on code with seemingly no practical use - if a compelling use case is presented, modification can easily be implemented.)

| Method                                                                                                                                           | Description                                                                                                                                            |
| ------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Selection.attr(name)                                 | Get the value of an attribute for the first element in the Selection.                                                                                  |
| Selection.children([selector])                   | Get the children of each element in the set of matched elements, optionally filtered by a selector.                                                    |
| Selection.closest(selector)                       | Get the first element that matches the selector by testing the element itself and traversing up through its ancestors                                  |
| Selection.contents()                             | Get the children of each element in the set of matched elements, including text and comment nodes.                                                     |
| Selection.data([key])                                | Return the value at the named data store for the first element in the set of matched elements.                                                         |
| Selection.each(fn)                                   | Iterate and execute a function for each matched element.                                                                                               |
| Selection.eq(index)                                    | Reduce the set of matched elements to the one at the specified index.                                                                                  |
| Selection.filter(selector)                         | Reduce the set of matched elements to those that match the selector or pass the function's test.                                                       |
| Selection.find(selector)                             | Find the selection descendants, filtered by a selector.                                                                                                |
| Selection.first()                                   | Reduce the set of matched elements to the first in the set.                                                                                            |
| Selection.get(index)                                  | Retrieve the Element (k6/html) matched by the selector                      |
| Selection.has(selector)                               | Reduce the set of matched elements to those that have a descendant that matches the selector                                                           |
| Selection.html()                                     | Get the HTML contents of the first element in the set of matched elements                                                                              |
| Selection.is(selector)                                 | Check the current matched set of elements against a selector or element and return true if at least one of these elements matches the given arguments. |
| Selection.last()                                     | Reduce the set of matched elements to the final one in the set.                                                                                        |
| Selection.map(fn)                                     | Pass each selection in the current matched set through a function, producing a new Array containing the return values.                                 |
| Selection.nextAll([selector])                     | Get all following siblings of each element in the set of matched elements, optionally filtered by a selector.                                          |
| Selection.next([selector])                           | Get the immediately following sibling of each element in the set of matched element                                                                    |
| Selection.nextUntil([selector], [filter])       | Get all following siblings of each element up to but not including the element matched by the selector.                                                |
| Selection.not(selector)                               | Remove elements from the set of matched elements                                                                                                       |
| Selection.parent([selector])                       | Get the parent of each element in the current set of matched elements, optionally filtered by a selector.                                              |
| Selection.parents([selector])                     | Get the ancestors of each element in the current set of matched elements, optionally filtered by a selector.                                           |
| Selection.parentsUntil([selector], [filter]) | Get the ancestors of each element in the current set of matched elements, up to but not including the element matched by the selector.                 |
| Selection.prevAll([selector])                     | Get all preceding siblings of each element in the set of matched elements, optionally filtered by a selector.                                          |
| Selection.prev([selector])                           | Get the immediately preceding sibling of each element in the set of matched elements.                                                                  |
| Selection.prevUntil([selector], [filter])       | Get all preceding siblings of each element up to but not including the element matched by the selector.                                                |
| Selection.serialize()                           | Encode a set of form elements as a string in standard URL-encoded notation for submission.                                                             |
| Selection.serializeArray()                 | Encode a set of form elements as an array of names and values.                                                                                         |
| Selection.serializeObject()               | Encode a set of form elements as an object.                                                                                                            |
| Selection.size()                                     | Return the number of elements in the Selection.                                                                                                        |
| Selection.slice(start [, end])                      | Reduce the set of matched elements to a subset specified by a range of indices.                                                                        |
| Selection.text()                                     | Get the text content of the selection.                                                                                                                 |
| Selection.toArray()                               | Retrieve all the elements contained in the Selection, as an array.                                                                                     |
| Selection.val()                                       | Get the current value of the first element in the set of matched elements.                                                                             |

### Example

```javascript
import { parseHTML } from 'k6/html';
import http from 'k6/http';

export default function () {
  const res = http.get('https://k6.io');
  const doc = parseHTML(res.body);
  const pageTitle = doc.find('head title').text();
  const langAttr = doc.find('html').attr('lang');
}
```

