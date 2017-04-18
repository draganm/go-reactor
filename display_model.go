package reactor

import "reflect"

type DisplayModel struct {
	ID           string                 `json:"id,omitempty"`
	Element      string                 `json:"el,omitempty"`
	Text         string                 `json:"te,omitempty"`
	Children     []*DisplayModel        `json:"ch,omitempty"`
	Attributes   map[string]interface{} `json:"at,omitempty"`
	ReportEvents []ReportEvent          `json:"ev,omitempty"`
}

type ReportEvent struct {
	Name            string   `json:"name,omitempty"`
	StopPropagation bool     `json:"sp,omitempty"`
	PreventDefault  bool     `json:"pd,omitempty"`
	ExtraValues     []string `json:"xv,omitempty"`
}

type DisplayUpdate struct {
	Model    *DisplayModel `json:"model,omitempty"`
	Eval     string        `json:"eval,omitempty"`
	Title    string        `json:"title,omitempty"`
	Location string        `json:"location,omitempty"`
}

func (d *DisplayUpdate) DeepEqual(other *DisplayUpdate) bool {

	return reflect.DeepEqual(d, other)
}

func (m *DisplayModel) FindElementByID(id string) *DisplayModel {
	if m.ID == id {
		return m
	}
	for _, child := range m.Children {
		found := child.FindElementByID(id)
		if found != nil {
			return found
		}
	}
	return nil
}

func (m *DisplayModel) ReplaceElementWithPath(path []int, replacement *DisplayModel) {
	if len(path) == 0 {
		panic("can't replace myself")
	}
	if len(path) == 1 {
		m.Children[path[0]] = replacement
	} else {
		m.Children[path[0]].ReplaceElementWithPath(path[1:], replacement)
	}
}

func (m *DisplayModel) findElementPathByID(id string, path []int) *[]int {
	if m.ID == id {
		return &path
	}
	for i, child := range m.Children {
		p := append(path, i)
		found := child.findElementPathByID(id, p)
		if found != nil {
			return found
		}
	}
	return nil
}

func (m *DisplayModel) FindElementPathByID(id string) *[]int {
	return m.findElementPathByID(id, []int{})
}

func (m *DisplayModel) SetElementAttribute(id, name string, value interface{}) {
	element := m.FindElementByID(id)
	if element != nil {
		element.Attributes[name] = value
	}
}

func (m *DisplayModel) SetElementText(id, text string) *DisplayModel {

	element := m.FindElementByID(id)
	if element != nil && text != "" {
		element.Children = []*DisplayModel{
			&DisplayModel{
				Text: text,
			},
		}
	}
	return element
}

func (m *DisplayModel) ReplaceChild(id string, replacement *DisplayModel) *DisplayModel {
	for index, child := range m.Children {
		if child.ID == id {
			m.Children[index] = replacement
		} else {
			child.ReplaceChild(id, replacement)
		}
	}
	return m
}

func (m *DisplayModel) DeleteChild(id string) {
	indexToDelete := -1
	for index, child := range m.Children {
		if child.ID == id {
			indexToDelete = index
			break
		} else {
			child.DeleteChild(id)
		}
	}
	if indexToDelete >= 0 {
		m.Children = append(m.Children[:indexToDelete], m.Children[indexToDelete+1:]...)
	}
}

func (m *DisplayModel) AppendChild(id string, child *DisplayModel) {
	element := m.FindElementByID(id)
	if element != nil {
		element.Children = append(element.Children, child)
	}
}

func (m *DisplayModel) DeepCopy() *DisplayModel {
	result := *m
	result.Children = []*DisplayModel{}
	for _, c := range m.Children {
		result.Children = append(result.Children, c.DeepCopy())
	}
	result.Attributes = map[string]interface{}{}
	for k, v := range m.Attributes {
		result.Attributes[k] = v
	}
	return &result
}

func (m *DisplayModel) DeepEqual(other *DisplayModel) bool {

	if m == other {
		return true
	}

	if m == nil {
		return false
	}

	if other == nil {
		return false
	}

	if m.ID != other.ID ||
		m.Element != other.Element ||
		m.Text != other.Text ||
		!reflect.DeepEqual(m.Attributes, other.Attributes) ||
		!reflect.DeepEqual(m.ReportEvents, other.ReportEvents) ||
		len(m.Children) != len(other.Children) {
		return false
	}

	for i, child := range m.Children {
		if !child.DeepEqual(other.Children[i]) {
			return false
		}
	}
	return true

}
