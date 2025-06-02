package monitor

import (
	"github.com/go-playground/validator/v10"
)

func CreateUpdateDtoStructLevelValidation(sl validator.StructLevel) {
	cfg := sl.Current().Interface().(CreateUpdateDto)

	// if float64(monitor.Timeout)*0.8 >= float64(monitor.Interval) {
	// 	ctx.JSON(http.StatusBadRequest, utils.NewFailResponse("Timeout cannot be greater than 80% of interval"))
	// 	return
	// }

	if float64(cfg.Timeout)*0.8 >= float64(cfg.Interval) {
		sl.ReportError(cfg.Timeout, "Timeout", "timeout", "timeout", "")
	}
}
