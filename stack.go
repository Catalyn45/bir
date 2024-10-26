package main

type Item[T any] struct {
	value T
	next *Item[T]
}

type Stack[T any] struct {
	head *Item[T]
}

func (this *Stack[T]) push(value T) {
	newItem := Item[T] {
		value: value,
		next: this.head,
	}

	this.head = &newItem
}

func (this *Stack[T]) pop() T {
	head := this.head

	this.head = head.next

	return head.value
}

func (this *Stack[T]) peek() T {
	return this.head.value
}
