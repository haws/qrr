qrr (query-replace-regexp)
==========================

This is a command line utillity to interactive find and replace in files. Its a rewrite of codemod.py (https://github.com/facebookarchive/codemod) in Go, with extra features I always wanted.

## Installation

If you have Go installed, run:

`go get -u github.com/haws/qrr`

## Usage

`qrr <from> <to>`

or, if you have defined projects in your ~/.qrr-config, do:

`qrr -p <project> <from> <to>` 

This was created to allow quick replace without changing directories.
