package ruffe

import (
	"encoding/json"
	"encoding/xml"
	"io"
)

var Content = ContentUtil{}

type ContentUtil struct{}

func (u ContentUtil) JSON(ctx Context) error {
	err := u.SetJSONMarshaller(ctx)
	if err != nil {
		return err
	}
	return u.SetJSONUnmarshaller(ctx)
}

func (u ContentUtil) SetJSONMarshaller(ctx Context) error {
	return u.setMarshaller(ctx, jsonContent{})
}

func (u ContentUtil) SetJSONUnmarshaller(ctx Context) error {
	return u.setMarshaller(ctx, jsonContent{})
}

func (u ContentUtil) XML(ctx Context) error {
	err := u.SetXMLMarshaller(ctx)
	if err != nil {
		return err
	}
	return u.SetXMLUnmarshaller(ctx)
}

func (u ContentUtil) SetXMLMarshaller(ctx Context) error {
	return u.setMarshaller(ctx, xmlContent{})
}

func (u ContentUtil) SetXMLUnmarshaller(ctx Context) error {
	return u.setMarshaller(ctx, xmlContent{})
}

func (ContentUtil) setMarshaller(ctx Context, rm responseMarshaller) error {
	ictx, _ := ctx.(*innerContext)
	if ictx == nil {
		return ErrUnsupportableContext
	}
	ictx.rm = rm
	return nil
}

func (ContentUtil) setUnmarshaller(ctx Context, ru requestUnmarshaller) error {
	ictx, _ := ctx.(*innerContext)
	if ictx == nil {
		return ErrUnsupportableContext
	}
	ictx.ru = ru
	return nil
}

type jsonContent struct{}

func (c jsonContent) Unmarshal(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func (c jsonContent) ContentType() string {
	return jsonContentTypeValue
}

func (c jsonContent) Marshal(w io.Writer, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

type xmlContent struct{}

func (c xmlContent) Unmarshal(r io.Reader, v interface{}) error {
	return xml.NewDecoder(r).Decode(v)
}

func (c xmlContent) ContentType() string {
	return xmlContentTypeValue
}

func (c xmlContent) Marshal(w io.Writer, v interface{}) error {
	return xml.NewEncoder(w).Encode(v)
}
