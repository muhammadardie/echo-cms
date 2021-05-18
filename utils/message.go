package utils

import "reflect"

type Message struct {
	Saved   string
	Updated string
	Deleted string
}

func GetMessage(m string) string {
	message := Message{}
	message.Saved = "Record saved successfully"
	message.Updated = "Record updated successfully"
	message.Deleted = "Record deleted successfully"

	r := reflect.ValueOf(message)
	rv := r.FieldByName(m)

	if rv.IsValid() {
		v := rv.Interface().(string)

		return v
	}

	return m
}
