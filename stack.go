package main

type Item[T any] struct {
	value T
	next *Item[T]
}

type Stack[T any] struct {
	head *Item[T]
	count int
}

func (this *Stack[T]) push(value T) {
	newItem := Item[T] {
		value: value,
		next: this.head,
	}

	this.head = &newItem
	this.count++
}

func (this *Stack[T]) pop() T {
	head := this.head

	this.head = head.next
	this.count--

	return head.value
}

func (this *Stack[T]) peek() T {
	return this.head.value
}

func (this *Stack[T]) len() int {
	return this.count
}

func (this *Stack[T]) foreach(iterator func(item T) (stop bool)) {
	for head := this.head; head != nil; head = head.next {
		stop := iterator(head.value)
		if stop {
			return
		}
	}
}
