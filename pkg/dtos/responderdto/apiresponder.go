package respondertdto

import "github.com/gofiber/fiber/v2"

type apiresponse struct {
	Code    int
	Data    interface{}
	Message string
	Ctx     *fiber.Ctx
}

func NewResponse(code int, data interface{}, message string, ctx *fiber.Ctx) *apiresponse {
	response := &apiresponse{
		Code:    code,
		Data:    data,
		Message: message,
		Ctx:     ctx,
	}
	return response
}

func (r *apiresponse) Responder() error {
	mapResults := map[string]interface{}{
		"status":  r.Code,
		"data":    r.Data,
		"message": r.Message,
	}
	return r.Ctx.Status(r.Code).JSON(mapResults)
}
