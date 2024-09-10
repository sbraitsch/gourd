package internal

import "embed"

//go:embed static/*
var StaticAssets embed.FS
