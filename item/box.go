package item

import "fmt"

type Box struct {
	From    string
	Content *Item
}

func (b Box) String() string {
	return fmt.Sprintf("📦 Посылка (%s)", b.From)
}
