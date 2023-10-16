package format

import "fmt"

const CannotSetName = "⚠️ Не удалось установить имя."

func LongName(maxLen int) string {
	return fmt.Sprintf("⚠️ Максимальная длина имени %d символов.", maxLen)
}

func NameSet(n string) string {
	return fmt.Sprintf("Имя %s установлено ✅", Name(n))
}

func YourName(n string) string {
	return fmt.Sprintf("Ваше имя: %s 🔖", Name(n))
}
