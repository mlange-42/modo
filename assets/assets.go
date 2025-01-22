package assets

import "embed"

//go:embed config/*
var Config embed.FS

//go:embed css/*
var CSS embed.FS

//go:embed templates/* templates/**/*
var Templates embed.FS
