package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var csvs []string = []string{"assets/goblin/goblin.csv", "assets/gopher/gopher.csv"}
var imgs []string = []string{"assets/goblin/goblin.png", "assets/gopher/gopher.png"}

func TestAssert(t *testing.T) {
	require.Equal(t, 123, 123, "they should be equal")
	require.NotEqual(t, 2, 456, "they should not be equal")
	require.FileExists(t, csvs[0])
}

func TestLoadSprite(t *testing.T) {
	require.Equal(t, len(csvs), len(imgs), "make sure imgs and csvs equal len")
	for idx := range csvs {
		resource := LoadSprite(imgs[idx], csvs[idx])
		require.Nil(t, resource.err)
	}
}

func TestLoadSpriteCsv(t *testing.T) {
	definition, err := LoadSpriteCsvDefinition(csvs[0])
	require.Nil(t, err)
	require.Equal(t, definition.gridSize, GridSize{5, 11})
	val, ok := definition.animationMap["WalkDown"]
	require.Equal(t, ok, true)
	require.Equal(t, "WalkDown", val.name)
	require.Equal(t, 0, val.rowIdx)
	require.Equal(t, 0, val.colIdxStart)
	require.Equal(t, 11, val.colIdxEnd)
}
