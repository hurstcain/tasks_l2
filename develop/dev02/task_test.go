package dev02

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpackString(t *testing.T) {
	validTestData := []struct {
		s        string
		expected string
	}{
		{
			s:        `a4bc2d5e`,
			expected: `aaaabccddddde`,
		},
		{
			s:        `abcd`,
			expected: `abcd`,
		},
		{
			s:        `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			s:        `qwe\45`,
			expected: `qwe44444`,
		},
		{
			s:        `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			s:        `\8\9\0\\`,
			expected: `890\`,
		},
		{
			s:        `g0j4`,
			expected: `jjjj`,
		},
		{
			s:        `\80\22ffg0`,
			expected: `22ff`,
		},
		{
			s:        `абв3о2\74`,
			expected: `абвввоо7777`,
		},
		{
			s:        `а1б1\74`,
			expected: `аб7777`,
		},
		{
			s:        ``,
			expected: ``,
		},
	}

	invalidTestData := []struct {
		s        string
		expected string
	}{
		{
			s:        `45`,
			expected: ``,
		},
		{
			s:        `4`,
			expected: ``,
		},
		{
			s:        `\666`,
			expected: ``,
		},
		{
			s:        `vbn56vbn`,
			expected: ``,
		},
		{
			s:        `\\\`,
			expected: ``,
		},
		{
			s:        `gp5\`,
			expected: ``,
		},
		{
			s:        `\`,
			expected: ``,
		},
	}

	for _, data := range validTestData {
		res, err := UnpackString(data.s)
		assert.Equal(t, data.expected, res)
		assert.NoError(t, err)
	}

	for _, data := range invalidTestData {
		res, err := UnpackString(data.s)
		assert.Equal(t, data.expected, res)
		assert.Error(t, err)
	}
}
