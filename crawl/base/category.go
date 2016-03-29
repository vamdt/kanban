package base

type CategoryItemInfo struct {
	Id     int
	Pid    int
	Factor int
	Leaf   bool
	Name   string
	Tag    string
}

type CategoryItem struct {
	Id   int
	Name string
	Info []CategoryItemInfo
	Sub  Category
}

type Category map[string]CategoryItem

func NewCategory() *Category {
	c := make(Category)
	return &c
}

func NewCategoryItem(name string) *CategoryItem {
	return &CategoryItem{Name: name}
}

func (p *CategoryItem) Assembly(data []CategoryItemInfo) {
	for i := len(data) - 1; i > -1; i-- {
		if p.Id != data[i].Pid {
			continue
		}

		name := data[i].Name

		if data[i].Leaf {
			p.Info = append(p.Info, data[i])
		} else {
			if p.Sub == nil {
				p.Sub = *NewCategory()
			}

			if _, ok := p.Sub[name]; !ok {
				item := NewCategoryItem(name)
				item.Id = data[i].Id
				p.Sub[name] = *item
			}
			item := p.Sub[name]
			item.Assembly(data)
			p.Sub[name] = item
		}
	}
}

func (p *CategoryItem) AddStock(id, tag string) {
	info := CategoryItemInfo{Name: id, Tag: tag}
	p.Info = append(p.Info, info)
}

func (p *CategoryItem) initSub() {
	if p.Sub == nil {
		p.Sub = *NewCategory()
	}
}
