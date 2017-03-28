package reactor

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func MustParseDisplayModel(src string) *DisplayModel {
	ret, err := ParseDisplayModel(src)
	if err != nil {
		panic(err)
	}
	return ret
}

func ParseDisplayModel(src string) (*DisplayModel, error) {
	src = strings.TrimSpace(src)
	reader := strings.NewReader(src)

	decoder := xml.NewDecoder(reader)

	stack := []*DisplayModel{}

	for {
		token, err := decoder.Token()
		// fmt.Println("token", token)

		if token == nil && err == io.EOF {
			var result *DisplayModel

			if len(stack) > 0 {
				result = stack[0]
			}
			return result, nil
		}

		if err != nil {
			return nil, err
		}

		switch t := token.(type) {
		case xml.StartElement:
			model := &DisplayModel{
				Element:    t.Name.Local,
				Attributes: map[string]interface{}{},
			}
			for _, attribute := range t.Attr {
				if attribute.Name.Local == "id" {
					model.ID = attribute.Value
				} else if attribute.Name.Local == "reportEvents" {
					re := []ReportEvent{}

					evts := strings.Split(attribute.Value, " ")

					for _, evtString := range evts {

						parts := strings.Split(evtString, ":")

						rest := parts[1:]

						sp := false
						pd := false

						extraValues := []string{}

						for _, v := range rest {
							switch v {
							case "SP":
								sp = true
							case "PD":
								pd = true
							default:
								if strings.HasPrefix(v, "X-") {
									extraValues = append(extraValues, v[2:])
								} else {
									return nil, fmt.Errorf("Unknown reportEvent parameter %#v", v)
								}
							}
						}

						evt := ReportEvent{
							Name:            parts[0],
							StopPropagation: sp,
							PreventDefault:  pd,
							ExtraValues:     extraValues,
						}

						re = append(re, evt)
					}

					model.ReportEvents = re
				} else {
					var value interface{} = attribute.Value
					if attribute.Name.Space == "bool" {
						value, err = strconv.ParseBool(attribute.Value)
						if err != nil {
							return nil, err
						}
					} else if attribute.Name.Space == "int" {
						value, err = strconv.Atoi(attribute.Value)
						if err != nil {
							return nil, err
						}
					}
					if attribute.Name.Local == "htmlID" {
						model.Attributes["id"] = value
					} else {
						model.Attributes[attribute.Name.Local] = value
					}
				}
			}
			if len(stack) != 0 {
				prev := stack[len(stack)-1]
				prev.Children = append(prev.Children, model)
			}
			stack = append(stack, model)

		case xml.CharData:
			if len(stack) != 0 {

				text := string(t)

				prev := stack[len(stack)-1]

				if text != "" && strings.TrimSpace(text) != "" {
					model := &DisplayModel{
						Text: text,
					}
					prev.Children = append(prev.Children, model)

				}
			}
		case xml.EndElement:
			if len(stack) > 1 {
				stack = stack[:len(stack)-1]
			}
		}
	}

}
