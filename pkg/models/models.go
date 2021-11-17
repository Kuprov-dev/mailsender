package models

//TODO add field for sender
type Template struct {
	UUID         string   `bson:"_id,omitempty"`
	Status       string   `bson:"status"`
	TemplateUUID string   `bson:"template_uuid"`
	Params       []string `bson:"params"`
	Receivers    []string `bson:"receivers"`
	Message      string   `bson:"message"`
}
type Message struct {
	Key   []byte
	Value []byte
}
