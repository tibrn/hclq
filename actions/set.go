package actions

import (
	"bytes"
	"io"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
	"github.com/tibrn/hclq/hclq"
)

func Set(filename string, reader io.Reader, queryString, newValue string) ([]byte, error) {

	doc, err := GetDocument(filename, reader)

	if err != nil {
		return nil, err
	}

	err = doc.Set(queryString,
		func(list *ast.ListType) error {
			listNode, err := hclq.HclListFromJSON(newValue)
			if err != nil {
				return err
			}
			list.List = listNode.List
			return nil
		}, func(tok *token.Token) error {
			tok.Text = `"` + newValue + `"`
			tok.Type = getTokenType(newValue)
			return nil
		})
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer([]byte{})
	doc.Print(buff)

	return buff.Bytes(), nil
}

func trimToken(tok string) string {
	return strings.Trim(tok, `"`)
}

func getTokenType(val string) token.Type {
	_, err := strconv.ParseInt(val, 0, 64)
	if err == nil {
		return token.NUMBER
	}
	_, err = strconv.ParseFloat(val, 64)
	if err == nil {
		return token.FLOAT
	}
	_, err = strconv.ParseBool(val)
	if err == nil {
		return token.BOOL
	}
	return token.STRING
	// TODO: support HEREDOC
}
