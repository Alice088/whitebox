package output

import (
	"encoding/json"
	"errors"
)

func ToAnswer(bytes []byte) (Answer, error) {
	answer := Answer{}

	var meta struct {
		Type string `json:"type"`
	}

	if err := json.Unmarshal(bytes, &meta); err != nil {
		return answer, err
	}

	switch meta.Type {

	case string(FinalType):
		var v Final
		if err := json.Unmarshal(bytes, &v); err != nil {
			return answer, err
		}
		answer.Type = FinalType
		answer.Struct = v
		return answer, nil

	case string(PlanType):
		var v Plan
		if err := json.Unmarshal(bytes, &v); err != nil {
			return answer, err
		}
		answer.Type = PlanType
		answer.Struct = v
		return answer, nil

	case string(ToolType):
		var v Tool
		if err := json.Unmarshal(bytes, &v); err != nil {
			return answer, err
		}
		answer.Type = ToolType
		answer.Struct = v
		return answer, nil

	case string(AskType):
		var v Ask
		if err := json.Unmarshal(bytes, &v); err != nil {
			return answer, err
		}
		answer.Type = AskType
		answer.Struct = v
		return answer, nil
	}

	return answer, errors.New("unknown type")
}
