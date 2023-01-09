package game

type Product struct {
	Price int
	Item  *Item
}

type Market struct {
	P    []*Product
	keys map[int]*Product
}

func (m *Market) Add(p *Product) {
	m.P = append(m.P, p)
}

func (m *Market) Products() []*Product {
	m.keys = map[int]*Product{}
	for i, p := range m.P {
		m.keys[i] = p
	}
	return m.P
}

func (u *User) Buy(m *Market, key int) (p *Product, ok bool) {
	p, ok = m.keys[key]
	if !ok {
		return nil, false
	}
	if ok := u.SpendMoney(p.Price); !ok {
		return nil, false
	}
	delete(m.keys, key)
	for i, v := range m.P {
		if v == p {
			m.P[i] = m.P[len(m.P)-1]
			m.P = m.P[:len(m.P)-1]
		}
	}
	u.Inventory = append(u.Inventory, p.Item)
	return p, true
}
