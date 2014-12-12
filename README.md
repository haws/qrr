qrr (query-replace-regexp)
==========================

This is a command line utillity to interactive find and replace in files. Its inspired in codemod.py (https://github.com/facebookarchive/codemod).

## Installation

If you have Go installed, run:

`go get -u github.com/hugows/qrr`

## Usage

`qrr <from> <to>`

or, if you have defined projects in your ~/.qrr-config, do:

`qrr -p <project> <from> <to>` 

This was created to allow quick replace without changing directories.
qrr
