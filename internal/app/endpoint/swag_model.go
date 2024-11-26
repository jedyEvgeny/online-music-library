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
	Message    string `json:"message" example:"ошибка создания json-объекта"`
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
	Message    string `json:"message" example:"не смогли прочитать параметр запроса"`
	StatusCode int    `json:"statusCode" example:"400"`
}

type ResponseLirycs404 struct {
	Sucsess    bool   `json:"sucsess" example:"false"`
	Message    string `json:"message" example:"не смогли прочитать параметр запроса"`
	StatusCode int    `json:"statusCode" example:"400"`
}
