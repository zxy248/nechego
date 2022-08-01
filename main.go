package main

import (
	"fmt"
	"math/rand"
	"nechego/app"
	"nechego/dice"
	"nechego/fight"
	"nechego/model"
	"nechego/numbers"
	"nechego/service"
	"nechego/statistics"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

const (
	tokenEnv = "NECHEGO_TOKEN" // bot token
	dsnEnv   = "NECHEGO_DSN"   // data source name
	ownerEnv = "NECHEGO_OWNER" // bot owner's id
	debugEnv = "NECHEGO_DEBUG" // set to 1 for debug output
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	b, err := tele.NewBot(tele.Settings{
		Token:  botToken(),
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	fail(err)

	var l *zap.Logger
	if debugFlag() {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	fail(err)
	defer l.Sync()

	m := model.New(database())

	t := statistics.New(m, statistics.Settings{
		EnergyRange:   numbers.MakeInterval(0, 5),
		PoorThreshold: 1000,
	})

	d := dice.New(dice.Settings{
		RollTime: time.Second * 25,
	})

	s := service.New(m, t, d, service.Config{
		EatEnergyRestore:   numbers.MakeInterval(1, 2),
		FishingRodPrice:    2990,
		FishingEnergyDrain: 1,
		FightSettings: fight.Settings{
			ChanceRatio:  0.5,
			StrengthFunc: t.Strength,
		},
		FightEnergyDrain:   1,
		WinReward:          numbers.MakeInterval(50, 2000),
		ParliamentMembers:  5,
		ParliamentMajority: 4,
		DepositFee:         50,
		WithdrawFee:        0,
		MinDebt:            1000,
		DebtPercentage:     0.2,
		InitialBalance:     1500,
		MinBet:             50,
		EnergyRestoreDelta: 1,
	})

	a := app.New(b, l, t, s, app.Preferences{
		DataPath:     "data",
		EnergyPeriod: time.Hour * 2 / 5,
		ListLength:   numbers.MakeInterval(3, 5),
		HelloChance:  0.2,
	})
	fail(a.InitStickers())
	a.Start()
}

func botToken() string {
	v := os.Getenv(tokenEnv)
	if v == "" {
		fail(fmt.Errorf("%s not set", tokenEnv))
	}
	return v
}

func database() *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", databaseDSN())
	return db
}

func databaseDSN() string {
	v := os.Getenv(dsnEnv)
	if v == "" {
		fail(fmt.Errorf("%s not set", dsnEnv))
	}
	return v
}

func ownerUID() int64 {
	owner := os.Getenv(ownerEnv)
	if owner == "" {
		fail(fmt.Errorf("%s not set", ownerEnv))
	}
	id, err := strconv.ParseInt(owner, 10, 64)
	if err != nil {
		fail(err)
	}
	return id
}

func debugFlag() bool {
	v := os.Getenv(debugEnv)
	return v == "1"
}

func fail(err error) {
	if err != nil {
		panic(err)
	}
}
