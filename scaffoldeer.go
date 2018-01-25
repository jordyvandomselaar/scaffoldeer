package main

import (
    "gopkg.in/urfave/cli.v1"
    "io/ioutil"
    "log"
    "os"
    "path"
    "path/filepath"
    "strings"
)

func main() {
    app := cli.NewApp()
    app.Name = "Scaffoldeer"
    app.Usage = "Scaffold stubs with ease."
    app.Version = "1.0.0"
    app.Commands = []cli.Command{
        {
            Name:   "make",
            Usage:  "Scaffold a new template",
            Action: scaffoldAction,
            Flags: []cli.Flag{
                cli.StringFlag{
                    Name:  "fields",
                    Usage: "Replacements. Syntax; --fields key:value,foo:bar",
                },
            },
        },
    }

    app.Run(os.Args)
}

// Holds all data we need to parse, and write stubs.
type Stub struct {
    FullPath, RelativePath, Name string
    Content                      []byte
}

// The scaffold action to be called from cli.
func scaffoldAction(c *cli.Context) {
    fields := c.String("fields")
    templateName := c.Args().Get(0)

    err := scaffold(templateName, fields)
    handleError(err)
}

// Scaffold a template.
func scaffold(templateName string, fields string) error {
    scriptPath, err := os.Executable()

    if err != nil {
        return err
    }

    scriptPath = path.Join(scriptPath, "../")
    templatePath := path.Join(scriptPath, "templates", templateName)
    stubsPath := path.Join(templatePath, "stubs")

    stubs, err := getStubs(stubsPath)

    if err != nil {
        return err
    }

    replacements := parseReplacements(fields)

    stubs = parseStubs(stubs, replacements)

    writeStubs(stubs)

    return nil
}

// Parse argument string into map. --fields name:foo -> map[name] = "foo"
func parseReplacements(fields string) map[string]string {
    replacements := make(map[string]string)

    parsedFields := strings.Split(fields, ",")

    for _, value := range parsedFields {
        parsedFieldsSeparated := strings.Split(value, ":")

        replacements[parsedFieldsSeparated[0]] = parsedFieldsSeparated[1]
    }

    return replacements
}

// Copy stubs to another the current folder.
func writeStubs(stubs []Stub) {
    for _, stub := range stubs {
        os.MkdirAll(stub.RelativePath, 0777)

        ioutil.WriteFile(
            path.Join(stub.RelativePath, stub.Name),
            stub.Content,
            0777,
        )
    }
}

// Parse stubs
func parseStubs(stubs []Stub, replacements map[string]string) []Stub {
    parsedStubs := []Stub{}

    for _, stub := range stubs {
        stubNewName := replacePlaceholders(stub.Name, replacements)
        stubNewName = strings.Replace(stubNewName, ".stub", "", -1)
        stubNewRelativePath := replacePlaceholders(stub.RelativePath, replacements)
        stubNewContent := parseFileContent(stub.Content, replacements)

        stub.Name = stubNewName
        stub.RelativePath = stubNewRelativePath
        stub.Content = stubNewContent

        parsedStubs = append(parsedStubs, stub)
    }

    return parsedStubs
}

// Get all stub files recursively.
func getStubs(stubsPath string) ([]Stub, error) {
    stubs := []Stub{}

    err := filepath.Walk(stubsPath, func(currentPath string, fileInfo os.FileInfo, err error) error {
        if err != nil || fileInfo.IsDir() {
            return err
        }

        fileContent, _ := ioutil.ReadFile(currentPath)
        relativePath := strings.Replace(currentPath, stubsPath, "", -1)
        relativePath = strings.Replace(relativePath, fileInfo.Name(), "", -1)
        relativePath = strings.Trim(relativePath, "/")

        stubs = append(stubs, Stub{
            FullPath:     currentPath,
            RelativePath: relativePath,
            Name:         fileInfo.Name(),
            Content:      fileContent,
        })

        return err
    })

    return stubs, err
}

// Replace placeholders in file content.
func parseFileContent(content []byte, replacements map[string]string) []byte {
    return []byte(replacePlaceholders(string(content), replacements))
}

// Replace placeholders with variables.
func replacePlaceholders(input string, replacements map[string]string) string {
    for key, value := range replacements {
        input = strings.Replace(input, "__"+key+"__", value, -1)
    }

    return input
}

// Log any errors to console.
func handleError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
