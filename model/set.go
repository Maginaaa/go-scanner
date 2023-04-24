package model

/**
以空结构体作为map的value来实现，空结构体不占内存
*/

type Empty struct{}

var empty Empty

// Set Set类型
type Set struct {
	M map[string]Empty
}

// NewSet 初始化Set
func NewSet() Set {
	return Set{M: make(map[string]Empty)}
}

// Add 添加元素
func (s *Set) Add(val string) {
	s.M[val] = empty
}

// Remove 删除元素
func (s *Set) Remove(val string) {
	delete(s.M, val)
}

// Size 获取长度
func (s *Set) Size() int {
	return len(s.M)
}

// Clear 清空set
func (s *Set) Clear() {
	s.M = make(map[string]Empty)
}

// Exist 查看某个元素是否存在
func (s *Set) Exist(val string) (ok bool) {
	_, ok = s.M[val]
	return
}

// KeySet 获取key列表
func (s *Set) KeySet() []string {
	list := make([]string, 0)
	for item := range s.M {
		list = append(list, item)
	}
	return list
}
