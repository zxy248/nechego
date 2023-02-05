// Code generated by "stringer -type=HandlerID -output=handlers_string.go"; DO NOT EDIT.

package handlers

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NoHandler-0]
	_ = x[InfaHandler-1]
	_ = x[WhoHandler-2]
	_ = x[ListHandler-3]
	_ = x[TopHandler-4]
	_ = x[MouseHandler-5]
	_ = x[TiktokHandler-6]
	_ = x[GameHandler-7]
	_ = x[WeatherHandler-8]
	_ = x[CatHandler-9]
	_ = x[AnimeHandler-10]
	_ = x[FurryHandler-11]
	_ = x[FlagHandler-12]
	_ = x[PersonHandler-13]
	_ = x[HorseHandler-14]
	_ = x[ArtHandler-15]
	_ = x[CarHandler-16]
	_ = x[SoyHandler-17]
	_ = x[DanbooruHandler-18]
	_ = x[FapHandler-19]
	_ = x[MasyunyaHandler-20]
	_ = x[PoppyHandler-21]
	_ = x[SimaHandler-22]
	_ = x[HelloHandler-23]
	_ = x[BasiliHandler-24]
	_ = x[CasperHandler-25]
	_ = x[ZeusHandler-26]
	_ = x[PicHandler-27]
	_ = x[AvatarHandler-28]
	_ = x[TurnOnHandler-29]
	_ = x[TurnOffHandler-30]
	_ = x[BanHandler-31]
	_ = x[UnbanHandler-32]
	_ = x[CalculatorHandler-33]
	_ = x[DailyEblanHandler-34]
	_ = x[DailyAdminHandler-35]
	_ = x[DailyPairHandler-36]
	_ = x[NameHandler-37]
	_ = x[InventoryHandler-38]
	_ = x[SortHandler-39]
	_ = x[CatchHandler-40]
	_ = x[DropHandler-41]
	_ = x[PickHandler-42]
	_ = x[FloorHandler-43]
	_ = x[MarketHandler-44]
	_ = x[NameMarketHandler-45]
	_ = x[BuyHandler-46]
	_ = x[EatHandler-47]
	_ = x[EatQuickHandler-48]
	_ = x[FishHandler-49]
	_ = x[CastNetHandler-50]
	_ = x[DrawNetHandler-51]
	_ = x[NetHandler-52]
	_ = x[FishingRecordsHandler-53]
	_ = x[CraftHandler-54]
	_ = x[StatusHandler-55]
	_ = x[SellHandler-56]
	_ = x[StackHandler-57]
	_ = x[CashoutHandler-58]
	_ = x[FightHandler-59]
	_ = x[ProfileHandler-60]
	_ = x[DiceHandler-61]
	_ = x[RollHandler-62]
	_ = x[TopStrongHandler-63]
	_ = x[TopRatingHandler-64]
	_ = x[TopRichHandler-65]
	_ = x[CapitalHandler-66]
	_ = x[BalanceHandler-67]
	_ = x[EnergyHandler-68]
	_ = x[NamePetHandler-69]
	_ = x[ReceiveSMSHandler-70]
	_ = x[SendSMSHandler-71]
	_ = x[ContactsHandler-72]
	_ = x[SpamHandler-73]
}

const _HandlerID_name = "NoHandlerInfaHandlerWhoHandlerListHandlerTopHandlerMouseHandlerTiktokHandlerGameHandlerWeatherHandlerCatHandlerAnimeHandlerFurryHandlerFlagHandlerPersonHandlerHorseHandlerArtHandlerCarHandlerSoyHandlerDanbooruHandlerFapHandlerMasyunyaHandlerPoppyHandlerSimaHandlerHelloHandlerBasiliHandlerCasperHandlerZeusHandlerPicHandlerAvatarHandlerTurnOnHandlerTurnOffHandlerBanHandlerUnbanHandlerCalculatorHandlerDailyEblanHandlerDailyAdminHandlerDailyPairHandlerNameHandlerInventoryHandlerSortHandlerCatchHandlerDropHandlerPickHandlerFloorHandlerMarketHandlerNameMarketHandlerBuyHandlerEatHandlerEatQuickHandlerFishHandlerCastNetHandlerDrawNetHandlerNetHandlerFishingRecordsHandlerCraftHandlerStatusHandlerSellHandlerStackHandlerCashoutHandlerFightHandlerProfileHandlerDiceHandlerRollHandlerTopStrongHandlerTopRatingHandlerTopRichHandlerCapitalHandlerBalanceHandlerEnergyHandlerNamePetHandlerReceiveSMSHandlerSendSMSHandlerContactsHandlerSpamHandler"

var _HandlerID_index = [...]uint16{0, 9, 20, 30, 41, 51, 63, 76, 87, 101, 111, 123, 135, 146, 159, 171, 181, 191, 201, 216, 226, 241, 253, 264, 276, 289, 302, 313, 323, 336, 349, 363, 373, 385, 402, 419, 436, 452, 463, 479, 490, 502, 513, 524, 536, 549, 566, 576, 586, 601, 612, 626, 640, 650, 671, 683, 696, 707, 719, 733, 745, 759, 770, 781, 797, 813, 827, 841, 855, 868, 882, 899, 913, 928, 939}

func (i HandlerID) String() string {
	if i < 0 || i >= HandlerID(len(_HandlerID_index)-1) {
		return "HandlerID(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _HandlerID_name[_HandlerID_index[i]:_HandlerID_index[i+1]]
}
