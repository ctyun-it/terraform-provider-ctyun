package utils

// BidirectionalMap 双向的映射map
type BidirectionalMap[K comparable, V comparable] struct {
	m1 map[K]V // key => value 映射关系
	m2 map[V]K // value => key映射关系
}

func NewBidirectionalMap[K comparable, V comparable]() *BidirectionalMap[K, V] {
	return &BidirectionalMap[K, V]{
		m1: make(map[K]V),
		m2: make(map[V]K),
	}
}

// Put 增加映射关系
func (b *BidirectionalMap[K, V]) Put(key K, value V) {
	b.m1[key] = value
	b.m2[value] = key
}

// GetValue 通过key获取value
func (b BidirectionalMap[K, V]) GetValue(key K) (V, bool) {
	v, ok := b.m1[key]
	return v, ok
}

// GetKey 通过value获取key
func (b BidirectionalMap[K, V]) GetKey(value V) (K, bool) {
	v, ok := b.m2[value]
	return v, ok
}
