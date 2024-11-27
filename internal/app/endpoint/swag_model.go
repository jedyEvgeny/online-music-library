package endpoint

type Song struct {
	Group string `json:"group" example:"Muse"`
	Song  string `json:"song" example:"Supermassive Black Hole"`
}

type ResponsePost201 struct {
	Sucsess    bool   `json:"sucsess" example:"true"`
	Message    string `json:"message" example:"Ресурс создан"`
	StatusCode int    `json:"statusCode" example:"201"`
	ID         int    `json:"resourceID" example:"77"`
}

type ResponsePost400 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"|пустое поле 'group'|"`
	StatusCode int    `json:"statusCode" example:"400"`
}

type ResponsePost405 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка метода. Ожидался: POST, имеется: GET"`
	StatusCode int    `json:"statusCode" example:"405"`
}

type ResponsePost500 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ресурс в стороннем хранилище не найден, код ответа: 404"`
	StatusCode int    `json:"statusCode" example:"500"`
}

type ResponseDel400 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"не смогли прочитать параметр запроса"`
	StatusCode int    `json:"statusCode" example:"400"`
}

type ResponseDel405 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка метода. Ожидался: DELETE, имеется: POST"`
	StatusCode int    `json:"statusCode" example:"405"`
}

type ResponseDel500 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка подключения к БД"`
	StatusCode int    `json:"statusCode" example:"500"`
}

type ResponseLirycs200 struct {
	Sucsess    bool     `json:"sucsess" example:"true"`
	Message    string   `json:"message" example:"Ресурс существует"`
	StatusCode int      `json:"statusCode" example:"200"`
	Lirycs     []string `json:"lyrics" example:"[\"[Verse 1]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\n...\", \"[Verse 2]\\nHey, Jude, don't be afraid\\nYou were made to go out and get her\\n...\"]"`
}

type ResponseLirycs400 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка: strconv.Atoi: parsing \"f\": invalid syntax. 'id_song' не число. Имеется: f"`
	StatusCode int    `json:"statusCode" example:"400"`
}

type ResponseLirycs404 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"отсутствуют данные в БД: sql: no rows in result set"`
	StatusCode int    `json:"statusCode" example:"404"`
}

type ResponseLirycs405 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка метода. Ожидался: GET, имеется: POST"`
	StatusCode int    `json:"statusCode" example:"405"`
}

type ResponseLirycs500 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка подключения к БД"`
	StatusCode int    `json:"statusCode" example:"500"`
}

type ResponseLibrary200 struct {
	Sucsess    bool   `json:"sucsess" example:"true"`
	Message    string `json:"message" example:"Ресурс существует"`
	StatusCode int    `json:"statusCode" example:"200"`
	Library    []struct {
		Group       string `json:"group" example:"The Beatles"`
		Song        string `json:"song" example:"Hey Jude"`
		ReleaseDate string `json:"releaseDate" example:"26.08.1968"`
		Lirycs      string `json:"lyrics" example:"[\"[Verse 1]\\nHey, Jude, don't make it bad\\nTake a sad song and make it better\\n...\", \"[Verse 2]\\nHey, Jude, don't be afraid\\nYou were made to go out and get her\\n...\"]"`
		Link        string `json:"link" example:"https://rutube.ru/video/6c6b701206f28fd2767d14f9b495e674/"`
	}
}

type ResponseLibrary400 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"не смогли прочитать параметр запроса"`
	StatusCode int    `json:"statusCode" example:"400"`
}

type ResponseLibrary404 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"отсутствуют данные в БД: sql: no rows in result set"`
	StatusCode int    `json:"statusCode" example:"404"`
}

type ResponseLibrary405 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка метода. Ожидался: GET, имеется: POST"`
	StatusCode int    `json:"statusCode" example:"405"`
}

type ResponseLibrary500 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка подключения к БД"`
	StatusCode int    `json:"statusCode" example:"500"`
}

type EnrichedSong struct {
	Group       string `json:"group" example:"Mobile"`
	Song        string `json:"song" example:"Hey Effective"`
	ReleaseDate string `json:"releaseDate" example:"27.11.2024"`
}

type ResponseUpdate200 struct {
	Sucsess    bool   `json:"sucsess" example:"true"`
	Message    string `json:"message" example:"Ресурс обновлён"`
	StatusCode int    `json:"statusCode" example:"200"`
}

type ResponseUpdate405 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"ошибка метода. Ожидался: PATCH, имеется: POST"`
	StatusCode int    `json:"statusCode" example:"405"`
}
