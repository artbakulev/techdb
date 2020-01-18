package models

const (
	CreateError     = "ошибка создания объекта"
	UpdateError     = "ошибка обновления объекта"
	DBParsingError  = "ошибка парсинга результатов"
	NotFoundError   = "объект не найден"
	InternalError   = "внутренняя ошибка"
	ConflictError   = "военный конфликт"
	BadRequestError = "невалидный запрос"
	DBError         = "конфикт данных в базе"
)

type Error struct {
	StatusCode int16  `json:"status_code"`
	Message    string `json:"message"`
}

func NewError(statusCode int16, message string, extraMessages ...string) *Error {
	if len(extraMessages) > 0 {
		for _, item := range extraMessages {
			message += ": " + item
		}
	}
	return &Error{
		StatusCode: statusCode,
		Message:    message,
	}
}
