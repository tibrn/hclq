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

func Set(reader io.Reader, queryString, newValue string) ([]byte, error) {

	doc, err := hclq.FromReader(reader)
	if err != nil {
		return nil, err
	}

	err = doc.Set(queryString,
		func(list *ast.ListType) error {
			node, err := hclq.HclFromJSON(newValue)
			if err != nil {
				return err
			}
			list.List = append(node.(*ast.ListType).List, list.List...)
			return nil
		}, func(tok *token.Token) error {
			tok.Text = `"` + newValue + trimToken(tok.Text) + `"`
			tok.Type = token.STRING
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
