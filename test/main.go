package main

func cha() {
	ch := make(chan interface{})
	ch1 := make(chan interface{})

	select {
	case <-ch:
		return

	case <-ch1:
		a := 1
		_ = a
		return
	}
}

func baz() {
	switch 0 {
	case 0:
		a := 1
		_ = a
		fallthrough
	case 1:
		a := 1
		_ = a
		break
	case 2:
		break
	}
}

func foo() int {
	v := []int{}
	for range v {
		return 0
	}

	for range v {
		for range v {
			return 0
		}
		return 0
	}

	o := []int{
		0, 1,
	}

	return o[0]
}

func bar() int {
	o := 1
	if o == 1 {
		if o == 0 {
			return 1
		}
		return 0
	}

	return o
}

func main() {
	return
}
