package app

import (
	"encoding/json"
	"jedyEvgeny/online-music-library/internal/app/service"
	"log"
	"net/http"
)

// emulateResponseFromRemoteService моделирует ответ отудалённой БД, хранящей первичную информацию о песнях
func emulateResponseFromRemoteService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	query := r.URL.Query()
	groupReq := query.Get("group")
	songReq := query.Get("song")

	for _, song := range emulateStorageSongs {
		if song.Group == groupReq && song.Song == songReq {
			resp := struct {
				ReleaseDate string `json:"releaseDate"`
				Lyrics      string `json:"text"`
				Link        string `json:"link"`
			}{
				ReleaseDate: song.ReleaseDate,
				Lyrics:      song.Lyrics,
				Link:        song.Link,
			}
			respJson, _ := json.Marshal(resp)
			w.WriteHeader(http.StatusOK)
			w.Write(respJson)
			return
		}
	}
	log.Printf("Песня `%s` исполнителя `%s` не найдена в эмуляторе хранилища:", songReq, groupReq)
	w.WriteHeader(http.StatusNotFound)
}

var emulateStorageSongs = []service.EnrichedSong{
	{
		Group:       "Muse",
		Song:        "Supermassive Black Hole",
		ReleaseDate: "16.07.2006",
		Lyrics:      "[Verse 1: Matt Bellamy]\\nOoh, baby, don't you know I suffer?\\nOoh, baby, can't you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\n[Pre-Chorus: Matt Bellamy]\\nOoh, you set my soul alight\\nOoh, you set my soul alight\\n\\n[Chorus: Matt Bellamy, Chris Wolstenholme, Dom Howard]\\nGlaciers melting in the dead of night (Ooh)\\nAnd the superstar's sucked into the supermassive (You set my soul alight)\\nGlaciers melting in the dead of night (Ooh)\\nAnd the superstar's sucked into the (You set my soul)",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	},
	{
		Group:       "Muse",
		Song:        "test 1",
		ReleaseDate: "16.07.2006",
		Lyrics:      "[Verse 1: Matt Bellamy]\\nOoh, baby, don't you know I suffer?\\nOoh, baby, can't you hear me moan?\\nYou caught me under false pretenses\\nHow long before you let me go?\\n\\n[Pre-Chorus: Matt Bellamy]\\nOoh, you set my soul alight\\nOoh, you set my soul alight\\n\\n[Chorus: Matt Bellamy, Chris Wolstenholme, Dom Howard]\\nGlaciers melting in the dead of night (Ooh)\\nAnd the superstar's sucked into the supermassive (You set my soul alight)\\nGlaciers melting in the dead of night (Ooh)\\nAnd the superstar's sucked into the (You set my soul)",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	},
	{
		Group:       "Nirvana",
		Song:        "Smells Like Teen Spirit",
		ReleaseDate: "10.09.1991",
		Lyrics:      "[Verse 1]\\nLoad up on guns, bring your friends\\nIt's fun to lose and to pretend\\nShe's over-bored and self-assured\\nOh no, I know a dirty word\\n\\n[Pre-Chorus]\\nHello, hello, hello, how low\\nHello, hello, hello, how low\\nHello, hello, hello, how low\\nHello, hello, hello\\n\\n[Chorus]\\nWith the lights out, it's less dangerous\\nHere we are now, entertain us\\nI feel stupid and contagious\\nHere we are now, entertain us\\nA mulatto, an albino\\nA mosquito, my libido, yeah",
		Link:        "https://rutube.ru/video/5ad78cc50603dedcaeba3ead4dbe0cfc/",
	},
	{
		Group:       "The Beatles",
		Song:        "Hey Jude",
		ReleaseDate: "26.08.1968",
		Lyrics:      "[Verse 1]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\nRemember to let her into your heart\\nThen you can start to make it better\\n\\n[Verse 2]\\nHey, Jude, don't be afraid\\nYou were made to go out and get her\\nThe minute you let her under your skin\\nThen you begin to make it better\\n[Bridge]\\nAnd anytime you feel the pain, hey, Jude, refrain\\nDon't carry the world upon your shoulders\\nFor well you know that it's a fool who plays it cool\\nBy making his world a little colder\\nNa-na-na-na-na, na-na-na-na\\n\\n[Verse 3]\\nHey, Jude, don't let me down\\nYou have found her, now go and get her\\n(Let it out and let it in)\\nRemember (Hey, Jude) to let her into your heart\\nThen you can start to make it better\\nSee upcoming rock shows\\nGet tickets for your favorite artists\\nYou might also like\\nFamily Matters\\nDrake\\nSo Long, London\\nTaylor Swift\\nloml\\nTaylor Swift\\n[Bridge]\\nSo let it out and let it in, hey, Jude, begin\\nYou're waiting for someone to perform with\\nAnd don't you know that it's just you, hey, Jude, you'll do\\nThe movement you need is on your shoulder\\nNa-na-na-na-na, na-na-na-na, yeah\\n\\n[Verse 4]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\nRemember to let her under your skin\\nThen you'll begin to make it (Woah, fucking hell!)\\nBetter, better, better, better, better, better, oh",
		Link:        "https://rutube.ru/video/6c6b701206f28fd2767d14f9b495e674/",
	},
	{
		Group:       "The Beatles",
		Song:        "test 1",
		ReleaseDate: "26.08.1968",
		Lyrics:      "[Verse 1]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\nRemember to let her into your heart\\nThen you can start to make it better\\n\\n[Verse 2]\\nHey, Jude, don't be afraid\\nYou were made to go out and get her\\nThe minute you let her under your skin\\nThen you begin to make it better\\n[Bridge]\\nAnd anytime you feel the pain, hey, Jude, refrain\\nDon't carry the world upon your shoulders\\nFor well you know that it's a fool who plays it cool\\nBy making his world a little colder\\nNa-na-na-na-na, na-na-na-na\\n\\n[Verse 3]\\nHey, Jude, don't let me down\\nYou have found her, now go and get her\\n(Let it out and let it in)\\nRemember (Hey, Jude) to let her into your heart\\nThen you can start to make it better\\nSee upcoming rock shows\\nGet tickets for your favorite artists\\nYou might also like\\nFamily Matters\\nDrake\\nSo Long, London\\nTaylor Swift\\nloml\\nTaylor Swift\\n[Bridge]\\nSo let it out and let it in, hey, Jude, begin\\nYou're waiting for someone to perform with\\nAnd don't you know that it's just you, hey, Jude, you'll do\\nThe movement you need is on your shoulder\\nNa-na-na-na-na, na-na-na-na, yeah\\n\\n[Verse 4]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\nRemember to let her under your skin\\nThen you'll begin to make it (Woah, fucking hell!)\\nBetter, better, better, better, better, better, oh",
		Link:        "https://rutube.ru/video/6c6b701206f28fd2767d14f9b495e674/",
	},
	{
		Group:       "The Beatles",
		Song:        "test 2",
		ReleaseDate: "26.08.1968",
		Lyrics:      "[Verse 1]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\nRemember to let her into your heart\\nThen you can start to make it better\\n\\n[Verse 2]\\nHey, Jude, don't be afraid\\nYou were made to go out and get her\\nThe minute you let her under your skin\\nThen you begin to make it better\\n[Bridge]\\nAnd anytime you feel the pain, hey, Jude, refrain\\nDon't carry the world upon your shoulders\\nFor well you know that it's a fool who plays it cool\\nBy making his world a little colder\\nNa-na-na-na-na, na-na-na-na\\n\\n[Verse 3]\\nHey, Jude, don't let me down\\nYou have found her, now go and get her\\n(Let it out and let it in)\\nRemember (Hey, Jude) to let her into your heart\\nThen you can start to make it better\\nSee upcoming rock shows\\nGet tickets for your favorite artists\\nYou might also like\\nFamily Matters\\nDrake\\nSo Long, London\\nTaylor Swift\\nloml\\nTaylor Swift\\n[Bridge]\\nSo let it out and let it in, hey, Jude, begin\\nYou're waiting for someone to perform with\\nAnd don't you know that it's just you, hey, Jude, you'll do\\nThe movement you need is on your shoulder\\nNa-na-na-na-na, na-na-na-na, yeah\\n\\n[Verse 4]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\nRemember to let her under your skin\\nThen you'll begin to make it (Woah, fucking hell!)\\nBetter, better, better, better, better, better, oh",
		Link:        "https://rutube.ru/video/6c6b701206f28fd2767d14f9b495e674/",
	},
}
