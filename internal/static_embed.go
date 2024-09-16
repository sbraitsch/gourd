package internal

import "embed"

// StaticAssets embeds the static folder into the binary so paths remain valid.
//
//go:embed static/*
var StaticAssets embed.FS
