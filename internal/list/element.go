package list

// Element is an element of a linked list.
type Element[T any] struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element[T]

	// The list to which this element belongs.
	list *List[T]

	// The value stored with this element.
	Value T
}

// Next returns the next list element or nil.
func (e *Element[T]) Next() *Element[T] {
	p := e.next
	if e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Prev returns the previous list element or nil.
func (e *Element[T]) Prev() *Element[T] {
	p := e.prev
	if e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}
