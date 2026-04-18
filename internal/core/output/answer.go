package output

import (
	"encoding/json"
	"errors"
)

func ToAnswer[T any](bytes []byte) (Answer[T], error) {
	var meta struct {
		Type Type `json:"type"`
	}

	if err := json.Unmarshal(bytes, &meta); err != nil {
		return Answer[T]{}, err
	}

	var result Answer[T]
	result.Type = meta.Type

	switch meta.Type {

	case FinalType:
		var v Final
		if err := json.Unmarshal(bytes, &v); err != nil {
			return result, err
		}
		casted, ok := any(v).(T)
		if !ok {
			return Answer[T]{}, errors.New("unexpected type: not Final")
		}
		result.Struct = casted
		return result, nil

	case PlanType:
		var v Plan
		if err := json.Unmarshal(bytes, &v); err != nil {
			return result, err
		}
		casted, ok := any(v).(T)
		if !ok {
			return Answer[T]{}, errors.New("unexpected type: not Plan")
		}
		result.Struct = casted
		return result, nil

	case ToolType:
		var v Tool
		if err := json.Unmarshal(bytes, &v); err != nil {
			return result, err
		}
		casted, ok := any(v).(T)
		if !ok {
			return Answer[T]{}, errors.New("unexpected type: not Tool")
		}
		result.Struct = casted
		return result, nil

	case AskType:
		var v Ask
		if err := json.Unmarshal(bytes, &v); err != nil {
			return result, err
		}
		casted, ok := any(v).(T)
		if !ok {
			return Answer[T]{}, errors.New("unexpected type: not Ask")
		}
		result.Struct = casted
		return result, nil
	}

	return Answer[T]{}, errors.New("unknown type")
}
