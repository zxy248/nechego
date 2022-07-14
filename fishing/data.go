package fishing

var speciesData = map[Species]struct {
	name          string
	normalWeight  float64
	maximumWeight float64
	constitution  Constitution
	pricePerKg    float64
	predator      bool
}{
	Pike:        {"Щука", 2.00, 35.00, Long, 400, true},
	Perch:       {"Окунь", 0.20, 5.00, Belly, 350, true},
	Zander:      {"Судак", 2.00, 18.00, Long, 400, true},
	Ruffe:       {"Ерш", 0.10, 0.20, Regular, 50, false},
	VolgaZander: {"Берш", 1.30, 2.00, Long, 430, true},
	Asp:         {"Жерех", 2.50, 4.50, Regular, 400, true},
	Chub:        {"Голавль", 0.75, 6.00, Regular, 300, false},
	Snakehead:   {"Змееголов", 3.00, 10.00, Long, 150, true},
	Burbot:      {"Налим", 4.50, 24.00, Long, 450, true},
	Eel:         {"Угорь", 2.00, 8.50, Long, 1500, true},
	Catfish:     {"Сом", 20.00, 150.00, Long, 500, true},
	Salmon:      {"Лосось", 4.00, 8.00, Regular, 1000, true},
	Grayling:    {"Хариус", 0.70, 1.40, Regular, 800, false},
	Trout:       {"Форель", 2.00, 10.00, Regular, 1000, true},
	Char:        {"Голец", 0.01, 0.025, Long, 50, false},
	Sturgeon:    {"Осетр", 18.00, 80.00, Long, 5000, true},
	Sterlet:     {"Стерлядь", 1.50, 8.00, Long, 1300, true},
	Carp:        {"Карп", 1.50, 24.00, Belly, 360, false},
	Goldfish:    {"Карась", 0.50, 5.00, Belly, 70, false},
	Tench:       {"Линь", 1.50, 7.50, Belly, 400, false},
	Bream:       {"Лещ", 1.00, 7.50, Belly, 100, false},
	Ide:         {"Язь", 1.00, 7.50, Regular, 300, false},
	Roach:       {"Плотва", 0.20, 2.00, Regular, 280, false},
	BigheadCarp: {"Толстолобик", 1.20, 16.00, Regular, 200, false},
	WhiteBream:  {"Белоглазка", 0.10, 0.80, Belly, 50, false},
	Rudd:        {"Красноперка", 0.30, 2.00, Belly, 100, false},
	Bleak:       {"Уклейка", 0.02, 0.06, Regular, 400, false},
	Nase:        {"Подуст", 0.40, 1.60, Regular, 180, false},
	Taimen:      {"Таймень", 4.00, 70.00, Long, 900, true},
}

var outcomeDescriptions = map[Outcome]string{
	Lost:     "Вы не смогли выудить рыбу.",
	Off:      "Рыба сорвалась с крючка.",
	Tear:     "Рыба сорвала леску.",
	Seagrass: "Рыба скрылась в водорослях.",
	Slip:     "Рыба выскользнула из рук.",
	Release:  "Вы отпустили рыбу обратно в воду.",
	Collect:  "Вы оставили рыбу себе.",
}
