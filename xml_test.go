package reactor

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	m := MustParseDisplayModel(`<test id="abc" htmlID="def" bool:testt="true" bool:testf="false" reportEvents="click input"> test </test>`)
	require.NotNil(t, m)
	require.Equal(t, "abc", m.ID)
	require.Equal(t, true, m.Attributes["testt"])
	require.Equal(t, false, m.Attributes["testf"])
	require.Equal(t, false, m.Attributes["testf"])
	require.Equal(t, "def", m.Attributes["id"])
	require.Equal(t, []ReportEvent{ReportEvent{Name: "click", ExtraValues: []string{}}, ReportEvent{Name: "input", ExtraValues: []string{}}}, m.ReportEvents)
	require.Equal(t, 1, len(m.Children))
	require.Equal(t, " test ", m.Children[0].Text)
}

func TestParserWithIntermittentTags(t *testing.T) {
	m := MustParseDisplayModel(`<test id="abc" htmlID="def" bool:testt="true" bool:testf="false" reportEvents="click input"> test <span>1</span> 2 </test>`)
	require.NotNil(t, m)
	require.Equal(t, "abc", m.ID)
	require.Equal(t, true, m.Attributes["testt"])
	require.Equal(t, false, m.Attributes["testf"])
	require.Equal(t, false, m.Attributes["testf"])
	require.Equal(t, "def", m.Attributes["id"])
	require.Equal(t, []ReportEvent{
		ReportEvent{
			Name:        "click",
			ExtraValues: []string{},
		},
		ReportEvent{
			Name:        "input",
			ExtraValues: []string{},
		},
	}, m.ReportEvents)
	require.Equal(t, 3, len(m.Children))
	require.Equal(t, " 2 ", m.Children[2].Text)
	require.Equal(t, " test ", m.Children[0].Text)

}

func TestParseWhitespace(t *testing.T) {
	cases := []struct {
		Source         string
		ExpectedStruct *DisplayModel
	}{
		{"<x/>",
			&DisplayModel{
				Element:    "x",
				Attributes: map[string]interface{}{},
			},
		},

		{"<x></x>",
			&DisplayModel{
				Element:    "x",
				Attributes: map[string]interface{}{},
			},
		},

		{"<x> <y/> </x>",
			&DisplayModel{
				Element:    "x",
				Attributes: map[string]interface{}{},
				Children: []*DisplayModel{
					{
						Element:    "y",
						Attributes: map[string]interface{}{},
					},
				},
			},
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %d", i+1), func(t *testing.T) {
			require.EqualValues(t, c.ExpectedStruct, MustParseDisplayModel(c.Source))
		})
	}
}

func TestParserWithReportEvent(t *testing.T) {
	m := MustParseDisplayModel(`<test id="abc" htmlID="def" bool:testt="true" bool:testf="false" reportEvents="click:PD input:SP:X-screenX"/>`)
	require.NotNil(t, m)
	require.Equal(t, []ReportEvent{
		ReportEvent{
			PreventDefault: true,
			Name:           "click",
			ExtraValues:    []string{},
		},
		ReportEvent{
			StopPropagation: true,
			Name:            "input",
			ExtraValues:     []string{"screenX"},
		},
	}, m.ReportEvents)

}
