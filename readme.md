# Scaffoldeer

## Installation Instructions

1. Visit the the [releases page](https://github.com/jordyvandomselaar/scaffoldeer/releases).
2. Download the right binary for your OS.
   1. .exe for windows, the rest has a suffix with the os name.
3. Create a `templates` folder next to the binary.

## Collaboration

Feel free to send in PR's if there's anything you're missing. You can also add new issues ofcourse!

## Creating new templates

1. Create a new folder inside the templates folder, the name of the folder will be the name of the template.
2. Inside the new folder create a `stubs` folder with the structure you want in terms of files.
3. Create files and folders the way you want them to scaffold.
   1. You can use variables in the naming as well!

## Syntax

**Adding variables to your stub**

Using a variable is as simple as writing `__varName__`.

**Defining a variable at the CLI**

When calling Scaffoldeer, you can add a `--fields` flag with the values. Example; `scaffoldeer make template --fields foo:bar,baz:end`