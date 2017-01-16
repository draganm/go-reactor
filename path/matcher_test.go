package path_test

import (
	"testing"

	"github.com/draganm/go-reactor/path"
	"github.com/stretchr/testify/assert"
)

func TestNewMatcher(t *testing.T) {
	t.Run("fixed path", func(t *testing.T) {
		m, err := path.NewMatcher("/x")
		if err != nil {
			t.Fatal(err)
		}
		t.Run("exact match", func(t *testing.T) {
			assert.NotNil(t, m("/x"))
		})
		t.Run("not matching when unmached", func(t *testing.T) {
			assert.Nil(t, m("/y"))
		})
		t.Run("not matching prefix", func(t *testing.T) {
			assert.Nil(t, m("/xy"))
		})
		t.Run("not matching postfix", func(t *testing.T) {
			assert.Nil(t, m("/y/x"))
		})
	})
	t.Run("path with one parameter", func(t *testing.T) {
		m, err := path.NewMatcher("/x/:test")
		if err != nil {
			t.Fatal(err)
		}
		t.Run("match empty", func(t *testing.T) {
			assert.NotNil(t, m("/x/"))
		})
		t.Run("match non-empty", func(t *testing.T) {
			assert.Equal(t, m("/x/y"), map[string]string{"test": "y"})
		})
	})

	t.Run("path with two parameters", func(t *testing.T) {
		m, err := path.NewMatcher("/x/:t1/:t2")
		if err != nil {
			t.Fatal(err)
		}
		t.Run("match empty", func(t *testing.T) {
			assert.Equal(t, m("/x//"), map[string]string{"t1": "", "t2": ""})
		})
		t.Run("match non-empty", func(t *testing.T) {
			assert.Equal(t, m("/x/y/z"), map[string]string{"t1": "y", "t2": "z"})
		})
	})

}
