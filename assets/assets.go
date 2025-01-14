package assets

import "embed"

//go:embed css/*
var CSS embed.FS

//go:embed templates/* templates/**/*
var Templates embed.FS
