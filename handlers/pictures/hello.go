package pictures

import (
	"math/rand/v2"

	"github.com/zxy248/nechego/data"
	"github.com/zxy248/nechego/handlers"

	tele "gopkg.in/zxy248/telebot.v3"
)

type Hello struct {
	Queries *data.Queries
}

var helloRe = handlers.NewRegexp("^!(п[рл]ив[а-я]*|хай|зд[ао]ров[а-я]*|ку|здрав[а-я]*)")

func (h *Hello) Match(c tele.Context) bool {
	return helloRe.MatchString(c.Text())
}

func (h *Hello) Handle(c tele.Context) error {
	ids := [...]string{
		"CAACAgEAAx0CY12kCQACBftioLLbAu3OUYYkwzh41Jw_X0zoUQAC_CkAAnj8xgWerLvLS8hteiQE",
		"CAACAgIAAx0CY12kCQACBfRioLLbf4duWXmZ7QS06LxZV-8mugACV4kAAp7OCwACGBkSA_W8tyQE",
		"CAACAgIAAx0CY12kCQACBfVioLLbf043SsXZNqTxwuJsMsOGDQAC3WUAAp7OCwABnFr-z1vx1C8kBA",
		"CAACAgIAAx0CY12kCQACBfZioLLb2oH1vv_0syXaQ3pyE9PU1gACYA8AAtqGuUrThzxY_vc1WCQE",
		"CAACAgIAAx0CY12kCQACBfdioLLb4Xgs-msdqR4AAS9FNUZtaLQAAvgRAALl91IDDREFfdhVJAIkBA",
		"CAACAgEAAx0CY12kCQACBfhioLLbL9fnFNg4aUde6hYLD_kv9gACYyQAAnj8xgU19GKxPrA5XCQE",
		"CAACAgEAAx0CY12kCQACBflioLLbz4915lu-YHs8MqNBSQ5vhAACCh0AAnj8xgV4GJ2ZXZ1CESQE",
		"CAACAgEAAx0CY12kCQACBfpioLLbOrH4w2MYm3VN6LCXUz0eAQACGR0AAnj8xgX0TmpI2tKNGCQE",
		"CAACAgEAAx0CY12kCQACBglioLLbSS3FGtJt5dabBBBrcxtkYQACDCoAAnj8xgUOBPMhFGzv_iQE",
		"CAACAgEAAx0CY12kCQACBfxioLLb6HzIOCcuQzuIWp-62XgS4AACkiQAAnj8xgWGON-bfv8gOiQE",
		"CAACAgEAAx0CY12kCQACBf1ioLLbk9zTv54xNLl10W5uPNa0vAACqiQAAnj8xgX2vSDcAni6-yQE",
		"CAACAgEAAx0CY12kCQACBf5ioLLb8KXXsZ1UqslsN_tLCX4WHQACyh0AAnj8xgXaId5ToeN52iQE",
		"CAACAgEAAx0CY12kCQACBf9ioLLb_Eys-SK9SWIMgCt82ZKNvgACZCkAAnj8xgXAyoLAaTnLJiQE",
		"CAACAgEAAx0CY12kCQACBhJioLLbKNGSovUZREsL_KtRGt0fVQAC8ygAAnj8xgUiPk_ZNFmHlyQE",
		"CAACAgEAAx0CY12kCQACBhtioLLbmyW8NGl_z9-ovOpYwhioPwACMCgAAnj8xgUEbbgxHJdctiQE",
		"CAACAgEAAx0CY12kCQACBkBioLLbhwZP12DfvTVzFs8-TG6JPgACAh4AAnj8xgXbqzNj0BpkUiQE",
		"CAACAgEAAx0CY12kCQACBghioLLbIiD64BuOAaqYElsvWKQguQACZSkAAnj8xgVg94bnbwLNVSQE",
		"CAACAgEAAx0CY12kCQACBgFioLLbXchNbfOAKr7b1tV1uLTlzwACgCcAAnj8xgVuTyVK5P-rhSQE",
		"CAACAgEAAx0CY12kCQACBgJioLLbGKSUizm3if73iRLDkiEZlwACbiIAAnj8xgV8RScRuZet6iQE",
		"CAACAgEAAx0CY12kCQACBgNioLLbcXu8nRKfPRO5HglmAe8VQAACcCUAAnj8xgUayxJ64_kwVyQE",
		"CAACAgEAAx0CY12kCQACBgRioLLbdkt63qBk6CfnQylnx6OW3QAC8xwAAnj8xgW_ViJVNnqmFCQE",
		"CAACAgEAAx0CY12kCQACBgVioLLbQvc-D_5hLiq3_-0IxhCDEwACaiQAAnj8xgVHzNm8na_xxCQE",
		"CAACAgEAAx0CY12kCQACBlBioLLbHZE0ahzAYjdRX_Kf2A9U0AACoR8AAnj8xgUUDNXITh0CHiQE",
		"CAACAgEAAx0CY12kCQACBgdioLLbmxU71Da5tgJzNhpcpwNlCwACgCkAAnj8xgXp4TeNURp0fSQE",
		"CAACAgIAAx0CY12kCQACBgABYqCy2_ZTSVSaR0xmkVoAATnGqoSGAAI7YwAC4KOCB-1nufkt6X47JAQ",
		"CAACAgEAAx0CY12kCQACBgpioLLbaF_DPgABqXgOJJ6H3DO2TD8AAnwkAAJ4_MYF5PaMTRSH08kkBA",
		"CAACAgEAAx0CY12kCQACBg5ioLLbZT6leedAm_S1aG3ArJTPQQACsiQAAnj8xgXvzXhoTMAAARYkBA",
		"CAACAgEAAx0CY12kCQACBg9ioLLbgBPvV-1WitF56pFTL4st9wACsSQAAnj8xgUw0hY-NMy1aiQE",
		"CAACAgEAAx0CY12kCQACBgtioLLbKwjlT1jdzDY3qUtQvVZGuQACbykAAnj8xgVMUVRjDhLa1iQE",
		"CAACAgEAAx0CY12kCQACBgxioLLbvo05nstr3lSAOTChumpxwAACeCkAAnj8xgXhHvm_Uxlw3iQE",
		"CAACAgEAAx0CY12kCQACBhBioLLbUqryI1VSoURhbodGLoaYkAAClSQAAnj8xgXnT2PzwzOGwiQE",
		"CAACAgEAAx0CY12kCQACBg1ioLLbbd7VdqIkWWt4-yRfiAx2VAACBSUAAnj8xgV7Qr6ouLh80yQE",
		"CAACAgEAAx0CY12kCQACBhZioLLb0NzFJIL0Uxj95zkjDnI5ygAC5ygAAnj8xgWWXQodP3SrNyQE",
		"CAACAgEAAx0CY12kCQACBhFioLLbOgZC40vKvMC2Z1fGA25XygAC7igAAnj8xgX5Vpq1ObIVrSQE",
		"CAACAgEAAx0CY12kCQACBhNioLLbDp4WBrGqhK9164YnHkyaFAACBCUAAnj8xgWBethYUsAHxSQE",
		"CAACAgEAAx0CY12kCQACBhRioLLb9KgCSTB-GPfV-EWaVEwyiwAC5iQAAnj8xgW_Py9b2MD31yQE",
		"CAACAgEAAx0CY12kCQACBhVioLLbCgAB9NdoHclTC9D_M-wXLC4AAuwoAAJ4_MYFnLM8pxnioYUkBA",
		"CAACAgEAAx0CY12kCQACBiVioLLbemdpmYIl-Icg0mUX2T6RAwACGSYAAnj8xgVu9vihrFsWmiQE",
		"CAACAgEAAx0CY12kCQACBhxioLLbudLAyHEV9Nlw28jV_hKtpgACqCkAAnj8xgVANB-YMqqQXyQE",
		"CAACAgEAAx0CY12kCQACBhhioLLb7W51bOcPH2fr2xHQCVYLhAAC5SgAAnj8xgXv7UxsUgSYcyQE",
		"CAACAgEAAx0CY12kCQACBhdioLLbRbID2rA9ZoN9lwQ60iVHnQAC5igAAnj8xgXvC8IS4By8yiQE",
		"CAACAgEAAx0CY12kCQACBh1ioLLbJWudK5bN-TBUJI97tIQTewACpykAAnj8xgW-Ev1cTA0jLiQE",
		"CAACAgEAAx0CY12kCQACBh5ioLLb50UQoakRvPAGmZ7smIhwlgACpSkAAnj8xgWKlfJTwESHgSQE",
		"CAACAgEAAx0CY12kCQACBhpioLLbqs7MGGOCnYYDPdbNjimkJAACSygAAnj8xgWJ0mimUjuytCQE",
		"CAACAgEAAx0CY12kCQACBh9ioLLbX71ct8esNMh8coAezk7uwwAClykAAnj8xgV8YMLdbacsQyQE",
		"CAACAgEAAx0CY12kCQACBiJioLLbBVKqYE2q0qMyQqI-BSs54AACHCAAAnj8xgW8yd1_ip2XjSQE",
		"CAACAgEAAx0CY12kCQACBiNioLLboH7TDS4Y8Nx02IMVmmn9JwACtigAAnj8xgVb_kgAAdseUrckBA",
		"CAACAgEAAx0CY12kCQACBiBioLLb495IoN_CMFr-h77FWC1dYwACZSQAAnj8xgWYJ8PBOHfWkSQE",
		"CAACAgEAAx0CY12kCQACBiFioLLbKDpjzeL1vL8i5oKMyyQzXgACZCQAAnj8xgXgHsmS12Z1KSQE",
		"CAACAgEAAx0CY12kCQACBiRioLLbA4JDRVpJpx_NuC1BvhqUiAACLycAAnj8xgXJzLAofkMp7iQE",
		"CAACAgEAAx0CY12kCQACBi5ioLLb2Y7mnqg-OmzZrTXbN-yTnwACkygAAnj8xgVHpp6Xa5n4EyQE",
		"CAACAgEAAx0CY12kCQACBiZioLLb0_zUbJEZyfeC0oygILof6wACQyIAAnj8xgXDrJku7Tpv4yQE",
		"CAACAgEAAx0CY12kCQACBjNioLLbUc1VgyeuImrH37m4bavMKgAC4hoAAnj8xgVhVkUWppbG5yQE",
		"CAACAgEAAx0CY12kCQACBidioLLbZN8x7TrB2EDjkyaHOg4PIgACZCcAAnj8xgVSrGWLwOrh9yQE",
		"CAACAgEAAx0CY12kCQACBihioLLbdiuAB6Rt-BuWvs9dFRVHBwACWycAAnj8xgWpKYS62PqPhCQE",
		"CAACAgEAAx0CY12kCQACBi9ioLLb7nTA2BeUktZs3goOTFwSWAACFRsAAnj8xgWIQ3VYr4Ya4iQE",
		"CAACAgEAAx0CY12kCQACBilioLLbIEGvfY9Ebw3rQ2JSx99DMwACfycAAnj8xgVLFRnnjNoXgyQE",
		"CAACAgEAAx0CY12kCQACBjBioLLb0cad90_UUE90eR7_eBvSHAACEBsAAnj8xgXwU4yuKxMpmyQE",
		"CAACAgEAAx0CY12kCQACBhlioLLbmAl2aztE6dDJCZCeY6GzNwAC4igAAnj8xgWLkvab4S_jZSQE",
		"CAACAgEAAx0CY12kCQACBjFioLLbnR-00rPfCFEG0cCbKFLD4QACDxsAAnj8xgU1Qg-_a04K7iQE",
		"CAACAgEAAx0CY12kCQACBipioLLbT9hw67SGL4mBzTFkC1399AACdycAAnj8xgUryAoC5zs3SyQE",
		"CAACAgEAAx0CY12kCQACBjJioLLbmw8IYDUWT8OznH-UhYOn5gAC0RoAAnj8xgXJpk2pE9H1ByQE",
		"CAACAgEAAx0CY12kCQACBjVioLLbaUKz2sIz9bl-t1czzDKffgACNyYAAnj8xgVi8Edhl2MJ-yQE",
		"CAACAgEAAx0CY12kCQACBjdioLLb8JHj3mPXy8dRViwHkL9PUQACPyMAAnj8xgWagvz_4SqYYiQE",
		"CAACAgEAAx0CY12kCQACBitioLLb6hgiF1reF-fgxVCmYs3Y9AACnygAAnj8xgUiBiaOQLIKkCQE",
		"CAACAgEAAx0CY12kCQACBixioLLb7c2U5JadoT6j1vr45s1FMQACOCIAAnj8xgWP9VhL_eWLeSQE",
		"CAACAgEAAx0CY12kCQACBjZioLLbZX98_FuwD9zgh5yvCAzThAACJyYAAnj8xgXSFPxFG-qb1SQE",
		"CAACAgEAAx0CY12kCQACBi1ioLLbAfOaSWqw0bGk4r-LZj5ZuwACXyIAAnj8xgWP5M3gEtmJiiQE",
		"CAACAgEAAx0CY12kCQACBjxioLLbWqiW6zK32JD8X47pjhLoVwAC1B0AAnj8xgUDLoSE_7OdEyQE",
		"CAACAgEAAx0CY12kCQACBjRioLLbuCao0FfLdY6qWKcLfvAlZwACRCYAAnj8xgXELAejzUeBfiQE",
		"CAACAgEAAx0CY12kCQACBjlioLLbjdZ0kJhyfiVNDgjkXLxrpwACLCMAAnj8xgXXcxDhyK4kxSQE",
		"CAACAgEAAx0CY12kCQACBj9ioLLbLLQ128NVw6uLZK9R65BciwAC_B0AAnj8xgWmy7QBxStKvSQE",
		"CAACAgEAAx0CY12kCQACBjpioLLbmsPiYKgBXmluAcnTm0V3ogACKB4AAnj8xgWVq8WL0NVffiQE",
		"CAACAgEAAx0CY12kCQACBkFioLLbrsTL51LfC4XSvDdhCxWklAAC8iAAAnj8xgVGO2vN4fuDviQE",
		"CAACAgEAAx0CY12kCQACBj1ioLLb6wnX-8efkmYDly0GqQMfDQAC2x0AAnj8xgXvPY6oOfbsPSQE",
		"CAACAgEAAx0CY12kCQACBjtioLLbUr_qT6rVzZty6Q3-WkKzmQACDx4AAnj8xgWg3XzQbDeqBSQE",
		"CAACAgEAAx0CY12kCQACBjhioLLb3xpw4VLnfl7e7kXVOdLXHwACIhsAAnj8xgV9AzR-nQGjaSQE",
		"CAACAgEAAx0CY12kCQACBj5ioLLbHAmoTRV_txMmKypg99aI1gAC8x0AAnj8xgXXEWuh1quA2yQE",
		"CAACAgEAAx0CY12kCQACBkRioLLbE0J4-t4qkSYD0RdKL0FBcQAC7CAAAnj8xgU_39Qh3y5Q_iQE",
		"CAACAgEAAx0CY12kCQACBkZioLLbQ1P0J85bRzNcaYHd2amdiAACGSQAAnj8xgXJeT-_s4_KwiQE",
		"CAACAgEAAx0CY12kCQACBkJioLLbMYp5Aleiqr3rqvKmVPAIYwACTR4AAnj8xgXIM0sTGpsH4SQE",
		"CAACAgEAAx0CY12kCQACBkdioLLbvrIbK63Waie_UvGrLj0I3AACTiYAAnj8xgUYEnZwxXfx0SQE",
		"CAACAgEAAx0CY12kCQACBkhioLLbKSMAATvMelCOtVi8AAHfb-xLAAJYKQACePzGBSjyl-SI3zwwJAQ",
		"CAACAgEAAx0CY12kCQACBkVioLLbWTdTGtT-W-ub50qiugjrsAACFiQAAnj8xgUe0-qBFzw1SSQE",
		"CAACAgEAAx0CY12kCQACBkNioLLbBGzvN_Bxjt9Btu7k4ele-gACLCAAAnj8xgWvnHYCXqLd3yQE",
		"CAACAgEAAx0CY12kCQACBkpioLLb7kYOKPSwHLiQqn6HI2BpCQAC6B4AAnj8xgUVkRoG84GT4yQE",
		"CAACAgEAAx0CY12kCQACBktioLLbAAEfHnnW6gH74-MjWO8XIDUAAuseAAJ4_MYFqGpe-YxMizUkBA",
		"CAACAgEAAx0CY12kCQACBkxioLLbSBq8rUe0CBQDIt0kBKCdkQAC5B4AAnj8xgUIEUF2l0HjbiQE",
		"CAACAgEAAx0CY12kCQACBk5ioLLbaP9hgRMwjcUcM-F7PS_s1gAC4x8AAnj8xgXAH2FrM6WXfCQE",
		"CAACAgEAAx0CY12kCQACBk1ioLLbP5gVAAHouS_SW6-ZbMJpnVAAAvMeAAJ4_MYF8uAS3j0sIrQkBA",
		"CAACAgEAAx0CY12kCQACBlFioLLbOh39IjmJMf8QwhUgM5-aXQAC3x8AAnj8xgXecBx5-WgIlyQE",
		"CAACAgEAAx0CY12kCQACBk9ioLLbQ2Tw0dL6-LOAUsZ7ultgIQAC6B8AAnj8xgXnaOGa-kuyaSQE",
		"CAACAgEAAx0CY12kCQACBlNioLLbXa_kuzH5Hwlxjx5T9OUYwQACSSYAAnj8xgXa8JYrRy7SQiQE",
		"CAACAgEAAx0CY12kCQACBldioLLbEVb8ndP0QQ46jY4K9aC_RQACTiYAAnj8xgUYEnZwxXfx0SQE",
		"CAACAgEAAx0CY12kCQACBlRioLLb0HNZrIxEzT-Ax-N5q9aeqQACRSQAAnj8xgVPdt76K_0AAW4kBA",
		"CAACAgEAAx0CY12kCQACBlVioLLbedJoMrLIPKhT49kucHqvYgACIx8AAnj8xgWziVW0lKA0MCQE",
		"CAACAgEAAx0CY12kCQACBlJioLLbf_gAAWbfpRqPMVzhTHJcCZsAAsofAAJ4_MYF3ozIO9JcAQABJAQ",
		"CAACAgEAAx0CY12kCQACBlZioLLbh5ibLRbqengxrUfULd8tZAACfCIAAnj8xgXeTbxXJTF1xyQE",
		"CAACAgEAAx0CY12kCQACBklioLLbEkYGBn-ersWhGN6xNU45dwACXCkAAnj8xgV5tmdodbGO3yQE",
		"CAACAgEAAx0CY12kCQACBgZioLLbk851Kgb8vO5y5Vq4SLg4PgACbSQAAnj8xgVdNqWpON0XLCQE",
		"CAACAgIAAx0CY12kCQACBlhioLMSWu2xFBdCjw3B1eWEBelElAACFxcAAv01EUqGUPLkGNe7xiQE",
		"CAACAgEAAx0CY12kCQACBl1ioLMSJ8lbBLFD8tKxG5o3ehd-9QACMScAAnj8xgUvB6Y3E2zUUCQE",
		"CAACAgEAAx0CY12kCQACBlpioLMS9bgeIAbej3vKAnor6q552gACZSIAAnj8xgWHscMnMMGf8iQE",
		"CAACAgIAAx0CY12kCQACBllioLMSr3KVoT21z3Q-cf2RLWVtwAACnGMAAuCjggeR7QqGgNhoNyQE",
		"CAACAgEAAx0CY12kCQACBlxioLMSAtKO9CIZ-OiDX0wIRtC-HQAChyIAAnj8xgXdaRU3rgzNoiQE",
		"CAACAgEAAx0CY12kCQACBltioLMSFdXhOVCEo2Y1ZS1usnbONQACfiIAAnj8xgV6713m0nudKCQE",
		"CAACAgEAAx0CY12kCQACBl5ioLMS4Bso5aJVVuvbVnMpXFLZeAACWScAAnj8xgU876zYyjoTmyQE",
		"CAACAgEAAx0CY12kCQACBmFioLMS8UIvqnpu-KU9frbMiME6EAACGSYAAnj8xgVu9vihrFsWmiQE",
		"CAACAgIAAx0CY12kCQACBmNioLMSgaGFDeGMBLBp4Hp9-T0qPAACRiUAAp7OCwABFv3CYkdW_PYkBA",
		"CAACAgIAAx0CY12kCQACBmJioLMSONHYLayxqQABskQK0rcR3OYAAmojAAKezgsAAWqBNMG9Ivx8JAQ",
		"CAACAgEAAx0CY12kCQACBl9ioLMShbrsIFXiG5fB4IEDpVdU6gACbCcAAnj8xgX1n5pITyBiYiQE",
		"CAACAgEAAx0CY12kCQACBmBioLMSrjF46de1W-9Ewi_oMC5mlwACQyIAAnj8xgXDrJku7Tpv4yQE",
	}
	var sticker tele.Sticker
	sticker.FileID = ids[rand.N(len(ids))]
	return c.Send(&sticker)
}
