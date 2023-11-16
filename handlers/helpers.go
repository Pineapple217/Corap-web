package handlers

import (
	"Corap-web/models"
	"errors"
)

func FormatDataType(plotTypeStr string) (models.DataType, error) {
	switch plotTypeStr {
	case string(models.Temp):
		return models.Temp, nil
	case string(models.CO2):
		return models.CO2, nil
	case string(models.Humidity):
		return models.Humidity, nil
	default:
		return "", errors.New("string is not a valid datatype")
	}
}
