/*
Package bahn implements methods for querying live trip information from within ICE trains of Deutsche Bahn.

See the LICENSE file for licensing details.
*/
package bahn // import "github.com/octo/icestat/bahn"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func unmarshalJSONP(r io.ReadCloser, v interface{}) error {
	defer r.Close()

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	b = bytes.TrimLeft(b, "(")
	b = bytes.TrimRight(b, ");\r\n")

	return json.Unmarshal(b, v)
}
