package csv

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const filename = "../data/fng.1000.csv.rot128"

func TestReadFileShouldGet1000Records(t *testing.T) {
	t.Run("test read file", func(t *testing.T) {
		transations, err := Read(filename)
		require.NoError(t, err)
		require.Equal(t, 1000, len(transations))
	})
}

func TestTransFormInvalid(t *testing.T) {
	t.Run("test transform length < 6 should be nil", func(t *testing.T) {
		str := []string{"name"}
		actual := transform(str)
		require.Nil(t, actual)
	})

	t.Run("test transform empty name should be nil", func(t *testing.T) {
		str := []string{"", "500", "card", "cvv", "month", "year"}
		actual := transform(str)
		require.Nil(t, actual)
	})

	t.Run("test transform amount <= 0 should be nil", func(t *testing.T) {
		str := []string{"name", "amount", "card", "cvv", "month", "year"}
		actual := transform(str)
		require.Nil(t, actual)
	})

	t.Run("test transform empty card should be nil", func(t *testing.T) {
		str := []string{"name", "500", "", "cvv", "month", "year"}
		actual := transform(str)
		require.Nil(t, actual)
	})

	t.Run("test transform empty cvv should be nil", func(t *testing.T) {
		str := []string{"name", "500", "card", "", "month", "year"}
		actual := transform(str)
		require.Nil(t, actual)
	})

	t.Run("test transform empty month should be nil", func(t *testing.T) {
		str := []string{"name", "500", "card", "cvv", "", "year"}
		actual := transform(str)
		require.Nil(t, actual)
	})

	t.Run("test transform empty year should be nil", func(t *testing.T) {
		str := []string{"name", "500", "card", "cvv", "month", ""}
		actual := transform(str)
		require.Nil(t, actual)
	})
}
