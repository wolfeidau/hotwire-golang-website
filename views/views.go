package views

import "embed"

// Content holds our static web server content.
//
//go:embed templates/* layouts/* messages/*
var Content embed.FS
