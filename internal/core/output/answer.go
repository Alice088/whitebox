package output

import (
	"encoding/json"
	"errors"
)

func ToAnswer(bytes []byte) (Answer, error) {
	answer := Answer{}

	var f Final
	if err := json.Unmarshal(bytes, &f); err == nil {
		answer.Struct = f
		answer.Type = FinalType
		return answer, nil
	}

	var p Plan
	if err := json.Unmarshal(bytes, &p); err == nil {
		answer.Struct = p
		answer.Type = PlanType
		return answer, nil
	}

	var t Tool
	if err := json.Unmarshal(bytes, &t); err == nil {
		answer.Struct = t
		answer.Type = ToolType
		return answer, nil
	}

	return answer, errors.New("invalid format")
}
