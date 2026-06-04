package websocket

import "encoding/json"

func MarshalEvent(eventType string, fields map[string]any) ([]byte, error) {
	out := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	out["type"] = eventType
	return json.Marshal(out)
}
