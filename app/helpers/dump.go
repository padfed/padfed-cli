package helpers

import (
	"bytes"
	"fmt"
	"io"
)

func Dump(w io.Writer, v interface{}, pretty bool) error {
	if w != nil {
		var s []byte
		if pretty {
			s = MustMarshalPretty(v)
		} else {
			s = MustMarshal(v)
		}
		_, err := io.Copy(w, bytes.NewReader(s))
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w)
		return err
	}
	return nil
}
