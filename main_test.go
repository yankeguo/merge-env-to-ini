package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeModifies(t *testing.T) {
	envs := []string{
		"DEMO_SectionA__test1=key_a=val_a",
		"DEMO_SectionA__test2=key_b=val_b",
		"DEMO_SectionB__test3=key_c=val_c",
		"DEMO___test4=key_d=val_d",
	}
	modifies := decodeModifies(envs, "DEMO_")
	assert.Equal(t, "val_a", modifies["SectionA"]["key_a"])
	assert.Equal(t, "val_b", modifies["SectionA"]["key_b"])
	assert.Equal(t, "val_c", modifies["SectionB"]["key_c"])
	assert.Equal(t, "val_d", modifies[""]["key_d"])
}

func TestApplyModifies(t *testing.T) {
	modifies := map[string]map[string]string{
		"section1": {
			"key_a": "val_a",
		},
		"section2": {
			"key_b": "val_b",
		},
		"section3": {
			"key_b": "val_b",
		},
	}
	lines := [][]byte{
		[]byte("[section1]"),
		[]byte(";; whatafdafdfad"),
		[]byte(";;; key_a   = val_9"),
		[]byte(" key_a   = val_0"),
		[]byte("[section2]"),
		[]byte(";; whatafdafdfad"),
		[]byte(" key_a   = val_0"),
	}
	expected := [][]byte{
		[]byte("[section1]"),
		[]byte(";; whatafdafdfad"),
		[]byte("key_a = val_a"),
		[]byte("; key_a   = val_0"),
		[]byte("[section2]"),
		[]byte(";; whatafdafdfad"),
		[]byte(" key_a   = val_0"),
		[]byte("key_b = val_b"),
		[]byte("[section3]"),
		[]byte("key_b = val_b"),
	}

	expectedRaw := string(bytes.Join(expected, []byte("\n")))
	resultRaw := string(bytes.Join(applyModifies(modifies, lines), []byte("\n")))

	assert.Equal(t, expectedRaw, resultRaw)
}
