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

func easyjsonF5aeb05cDecodeMailInternalModelsDeliveryModels(in *jlexer.Lexer, out *OtherLabel) {
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
			out.ID = string(in.String())
		case "profileId":
			out.ProfileId = uint32(in.Uint32())
		case "name":
			out.Name = string(in.String())
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
func easyjsonF5aeb05cEncodeMailInternalModelsDeliveryModels(out *jwriter.Writer, in OtherLabel) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.ProfileId != 0 {
		const prefix string = ",\"profileId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint32(uint32(in.ProfileId))
	}
	{
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v OtherLabel) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF5aeb05cEncodeMailInternalModelsDeliveryModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v OtherLabel) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF5aeb05cEncodeMailInternalModelsDeliveryModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *OtherLabel) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF5aeb05cDecodeMailInternalModelsDeliveryModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *OtherLabel) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF5aeb05cDecodeMailInternalModelsDeliveryModels(l, v)
}
func easyjsonF5aeb05cDecodeMailInternalModelsDeliveryModels1(in *jlexer.Lexer, out *LabelEmail) {
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
		case "folderId":
			out.LabelID = string(in.String())
		case "emailId":
			out.EmailID = string(in.String())
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
func easyjsonF5aeb05cEncodeMailInternalModelsDeliveryModels1(out *jwriter.Writer, in LabelEmail) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"folderId\":"
		out.RawString(prefix[1:])
		out.String(string(in.LabelID))
	}
	{
		const prefix string = ",\"emailId\":"
		out.RawString(prefix)
		out.String(string(in.EmailID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v LabelEmail) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonF5aeb05cEncodeMailInternalModelsDeliveryModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v LabelEmail) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF5aeb05cEncodeMailInternalModelsDeliveryModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *LabelEmail) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonF5aeb05cDecodeMailInternalModelsDeliveryModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *LabelEmail) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF5aeb05cDecodeMailInternalModelsDeliveryModels1(l, v)
}
