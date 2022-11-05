package app

import (
	"fmt"
	"nechego/model"
	"time"

	tele "gopkg.in/telebot.v3"
)

const (
	kickVotesNeeded = 5
	kickDuration    = 5 * time.Minute
)

type kickEvent int

const (
	kickInit kickEvent = iota
	kickVote
	kickCancel
	kickDuplicate
	kickWrong
	kickNoUser
	kickSuccess
)

type kickStatus struct {
	event          kickEvent
	kicked         int64
	votesRemaining int
}

type kickRequest struct {
	gid     int64
	voter   int64
	kicked  int64
	timeout func()
	reply   chan kickStatus
}

type kickSession struct {
	kicked int64
	voters map[int64]bool
	timer  *time.Timer
}

func kickSystem(c chan kickRequest) {
	m := map[int64]*kickSession{}
	cancel := make(chan int64)
	for {
		select {
		case r := <-c:
			session, ok := m[r.gid]
			switch {
			case !ok && r.kicked == 0:
				r.reply <- kickStatus{event: kickNoUser}
			case !ok:
				session = &kickSession{
					kicked: r.kicked,
					voters: map[int64]bool{
						r.voter: true,
					},
					timer: time.AfterFunc(kickDuration, func() {
						r.timeout()
						cancel <- r.gid
					}),
				}
				m[r.gid] = session
				r.reply <- kickStatus{
					event:          kickInit,
					kicked:         session.kicked,
					votesRemaining: kickVotesNeeded - 1,
				}
			case r.kicked != 0 && r.kicked != session.kicked:
				r.reply <- kickStatus{event: kickWrong, kicked: session.kicked}
			case session.voters[r.voter]:
				r.reply <- kickStatus{event: kickDuplicate, kicked: session.kicked}
			default:
				session.voters[r.voter] = true
				rem := kickVotesNeeded - len(session.voters)
				event := kickVote
				if rem == 0 {
					event = kickSuccess
					session.timer.Stop()
					delete(m, r.gid)
				}
				r.reply <- kickStatus{
					event:          event,
					kicked:         session.kicked,
					votesRemaining: rem,
				}
			}
			close(r.reply)
		case i := <-cancel:
			delete(m, i)
		}
	}
}

var kickChannel = func() chan kickRequest {
	c := make(chan kickRequest)
	go kickSystem(c)
	return c
}()

func voteKick(gid, voter, kicked int64, timeout func()) kickStatus {
	reply := make(chan kickStatus, 1)
	kickChannel <- kickRequest{gid, voter, kicked, timeout, reply}
	return <-reply
}

var (
	responseKickInit = Response(`<b>✍️ Началось голосование за исключение пользователя %v.</b>

<i>Необходимо набрать еще %v.</i>
<i>Голосование закончится через 5 минут.</i>`)
	responseKickVote = Response(`<b>✍️ Вы проголосовали за исключение пользователя %v.</b>

<i>Необходимо набрать еще %v.</i>`)
	responseKickSuccess = Response("<b>😵 Пользователь <s>%v</s> исключен из беседы.</b>")

	responseKickDuplicate = UserError("Вы уже проголосовали.")
	responseKickWrong     = UserError("Подождите пока закончится голосование, прежде чем начать новое.")
	responseKickNoUser    = UserError("Перешлите сообщение пользователя, которого вы хотите исключить.")
	responseKickCancel    = UserError("Голосование за исключение пользователя %v истекло.")
)

func formatVote(n int) string {
	n0 := n % 10
	var s string
	if n0 == 1 {
		s = "голос"
	} else if n0 >= 2 && n0 <= 4 {
		s = "голоса"
	} else {
		s = "голосов"
	}
	return fmt.Sprintf("%d %s", n, s)
}

func (a *App) handleKick(c tele.Context) error {
	sender := getUser(c)
	reply, ok := maybeGetReplyUser(c)
	var status kickStatus
	if ok {
		status = voteKick(sender.GID, sender.UID, reply.UID, func() {
			respondUserError(c, responseKickCancel.Fill(a.mustMention(reply)))
		})
	} else {
		status = voteKick(sender.GID, sender.UID, 0, func() {})
	}
	if status.event == kickNoUser {
		return respondUserError(c, responseKickNoUser)
	}
	kicked, err := a.service.FindUser(model.User{GID: sender.GID, UID: status.kicked})
	if err != nil {
		return respondInternalError(c, err)
	}
	switch status.event {
	case kickInit:
		return respond(c, responseKickInit.Fill(a.mustMention(kicked), formatVote(status.votesRemaining)))
	case kickVote:
		return respond(c, responseKickVote.Fill(a.mustMention(kicked), formatVote(status.votesRemaining)))
	case kickCancel:
		return respondUserError(c, responseKickCancel.Fill(a.mustMention(kicked)))
	case kickDuplicate:
		return respondUserError(c, responseKickDuplicate)
	case kickWrong:
		return respondUserError(c, responseKickWrong)
	case kickSuccess:
		return respond(c, responseKickSuccess.Fill(a.mustMention(kicked)))
	}
	return nil
}
