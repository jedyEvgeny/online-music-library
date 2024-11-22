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
		Lyrics: `[Verse 1: Matt Bellamy]
Ooh, baby, don't you know I suffer?
Ooh, baby, can't you hear me moan?
You caught me under false pretenses
How long before you let me go?

[Pre-Chorus: Matt Bellamy]
Ooh, you set my soul alight
Ooh, you set my soul alight

[Chorus: Matt Bellamy, Chris Wolstenholme, Dom Howard]
Glaciers melting in the dead of night (Ooh)
And the superstar's sucked into the supermassive (You set my soul alight)
Glaciers melting in the dead of night (Ooh)
And the superstar's sucked into the (You set my soul)`,
		Link: "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	},
	{
		Group:       "Nirvana",
		Song:        "Smells Like Teen Spirit",
		ReleaseDate: "10.09.1991",
		Lyrics: `[Verse 1]
Load up on guns, bring your friends
It's fun to lose and to pretend
She's over-bored and self-assured
Oh no, I know a dirty word

[Pre-Chorus]
Hello, hello, hello, how low
Hello, hello, hello, how low
Hello, hello, hello, how low
Hello, hello, hello

[Chorus]
With the lights out, it's less dangerous
Here we are now, entertain us
I feel stupid and contagious
Here we are now, entertain us
A mulatto, an albino
A mosquito, my libido, yeah`,
		Link: "https://rutube.ru/video/5ad78cc50603dedcaeba3ead4dbe0cfc/",
	},
	{
		Group:       "The Beatles",
		Song:        "Hey Jude",
		ReleaseDate: "26.08.1968",
		Lyrics: `[Verse 1]
Hey, Jude, don't make it bad
Take a sad song and make it better
Remember to let her into your heart
Then you can start to make it better

[Verse 2]
Hey, Jude, don't be afraid
You were made to go out and get her
The minute you let her under your skin
Then you begin to make it better

[Bridge]
And anytime you feel the pain, hey, Jude, refrain
Don't carry the world upon your shoulders
For well you know that it's a fool who plays it cool
By making his world a little colder
Na-na-na-na-na, na-na-na-na

[Verse 3]
Hey, Jude, don't let me down
You have found her, now go and get her
(Let it out and let it in)
Remember (Hey, Jude) to let her into your heart
Then you can start to make it better
See upcoming rock shows
Get tickets for your favorite artists
You might also like
Family Matters
Drake
So Long, London
Taylor Swift
loml
Taylor Swift
[Bridge]
So let it out and let it in, hey, Jude, begin
You're waiting for someone to perform with
And don't you know that it's just you, hey, Jude, you'll do
The movement you need is on your shoulder
Na-na-na-na-na, na-na-na-na, yeah

[Verse 4]
Hey, Jude, don't make it bad
Take a sad song and make it better
Remember to let her under your skin
Then you'll begin to make it (Woah, fucking hell!)
Better, better, better, better, better, better, oh`,
		Link: "https://rutube.ru/video/6c6b701206f28fd2767d14f9b495e674/",
	},
}
