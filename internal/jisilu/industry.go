package jisilu

// 行业
type Node struct {
	Name   string
	Number string
	Level  int
}
type Industry struct {
	list  []*Node
	index map[string]int
}

func (i *Industry) Add(level int, number string, name string) {

	i.list = append(i.list, &Node{
		Number: number,
		Level:  level,
		Name:   name,
	})
	if i.index == nil {
		i.index = make(map[string]int)
	}
	i.index[number] = len(i.list) - 1
}

// 往上找
func (i *Industry) TopLevel(nm string) *Node {
	if n, ok := i.index[nm]; ok {
		for i.list[n].Level != 1 && n > 0 {
			n--
		}
		return i.list[n]
	}
	return nil
}

// 往下找
func (i *Industry) SubLevel(nm string) (res []*Node) {
	if n, ok := i.index[nm]; ok {
		res = append(res, i.list[n])
		n++
		for n < len(i.list) {
			if i.list[n].Level == 1 {
				break
			}
			res = append(res, i.list[n])
			n++
		}
	}
	return
}
