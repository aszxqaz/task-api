package fromjson

// json.Decoder парсит в map[string]any целые числа как float64.
// Вместе с тем, по соображениям гибкости точно типизировать
// некоторые структуры публичного апи представляется плохим решением,
// как и выносить парсинг JSON за пределы внешнего слоя приложения в область
// бизнес-логики.
func ParseInt(raw any) (int, bool) {
	var r *int
	if f, ok := raw.(float64); ok {
		d := int(f)
		r = &d
	}
	if i, ok := raw.(int); ok {
		r = &i
	}
	if r == nil {
		return 0, false
	}
	return *r, true
}
