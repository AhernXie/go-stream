package stream

import (
	"fmt"
	"testing"
)

func Test_stream(t *testing.T) {
	var arr = []int{2, 1, 3, 6, 5, 4}
	collect := Just(arr...).Filter(func(item int) bool {
		return item%2 == 0
	}).Sort(func(a int, b int) bool {
		return a < b
	}).Skip(1).Head(1).Collect()
	fmt.Println(collect)
}

type Person struct {
	Name string
	Age  int
}

func Test_collect_map(t *testing.T) {
	list := []*Person{{Name: "a", Age: 1}, {Name: "b", Age: 2}, {Name: "c", Age: 3}}
	print(list)
}

func Page[T any, E any](pageNo int64, pageSize int64, t []T, filterBefore FilterFunc[T], fn MapFunc[T, E], filterAfter FilterFunc[E], less LessFunc[E]) ([]E, int, error) {
	var length = len(t)
	if len(t) == 0 {
		return make([]E, 0), 0, nil
	}
	var ts = t
	if filterBefore != nil {
		ts = Just(t...).Filter(func(item T) bool {
			return filterBefore(item)
		}).Collect()
		length = len(ts)
	}
	just := MapAndJust(ts, func(item T) E {
		return fn(item)
	})
	if filterAfter != nil {
		es := just.Filter(func(item E) bool {
			return filterAfter(item)
		}).Collect()
		just = Just[E](es...)
		length = len(es)
	}
	if less != nil {
		just = just.Sort(func(before E, after E) bool {
			return less(before, after)
		})
	}
	collect := just.Skip(pageSize * (pageNo - 1)).Head(pageSize).Collect()
	return collect, length, nil
}
