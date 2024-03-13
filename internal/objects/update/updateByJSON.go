package update

type MetricRequestObj struct {
	ID    string `json:"id"`              // имя метрики
	MType string `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta string `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value string `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
