// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package delivery_models

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA4a3e1daDecodeMailInternalModelsDeliveryModels(in *jlexer.Lexer, out *Answer) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = uint32(in.Uint32())
		case "question_id":
			out.QuestionId = uint32(in.Uint32())
		case "login":
			out.Login = string(in.String())
		case "mark":
			out.Mark = uint32(in.Uint32())
		case "text":
			out.Text = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA4a3e1daEncodeMailInternalModelsDeliveryModels(out *jwriter.Writer, in Answer) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.ID))
	}
	if in.QuestionId != 0 {
		const prefix string = ",\"question_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint32(uint32(in.QuestionId))
	}
	if in.Login != "" {
		const prefix string = ",\"login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Login))
	}
	if in.Mark != 0 {
		const prefix string = ",\"mark\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint32(uint32(in.Mark))
	}
	if in.Text != "" {
		const prefix string = ",\"text\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Text))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Answer) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA4a3e1daEncodeMailInternalModelsDeliveryModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Answer) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA4a3e1daEncodeMailInternalModelsDeliveryModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Answer) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA4a3e1daDecodeMailInternalModelsDeliveryModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Answer) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA4a3e1daDecodeMailInternalModelsDeliveryModels(l, v)
}
