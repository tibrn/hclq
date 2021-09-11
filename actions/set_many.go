package actions

import (
	"bytes"
	"io"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/tibrn/hclq/hclq"
)

type Replace struct {
	Query    string
	NewValue string
}

func SetMany(filename string, reader io.Reader, replacements []Replace) ([]byte, error) {

	doc, err := GetDocument(filename, reader)

	if err != nil {
		return nil, err
	}

	for _, replacement := range replacements {

		err := doc.Set(replacement.Query,
			func(list *ast.ListType) error {
				listNode, err := hclq.HclListFromJSON(replacement.NewValue)
				if err != nil {
					return err
				}
				list.List = listNode.List
				return nil
			}, func(tok *token.Token) error {
				tok.Text = `"` + replacement.NewValue + `"`
				tok.Type = getTokenType(replacement.NewValue)
				return nil
			})

		if err != nil {
			return nil, err
		}
	}

	buff := bytes.NewBuffer([]byte{})
	doc.Print(buff)

	return buff.Bytes(), nil
}
