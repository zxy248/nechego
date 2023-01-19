package tools

import "fmt"

type Knife struct {
	Durability float64
}

func (k Knife) String() string {
	return fmt.Sprintf("🔪 Нож (%.f%%)", k.Durability*100)
}
