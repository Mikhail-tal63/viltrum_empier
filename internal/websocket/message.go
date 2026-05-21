package websocket

import "encoding/json"

// MarshalEvent produces a flat JSON payload matching the frontend ServerEvent
// shape: {"type": "<eventType>", ...fields}.
func MarshalEvent(eventType string, fields map[string]any) ([]byte, error) {
	out := make(map[string]any, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	out["type"] = eventType
	return json.Marshal(out)
}
