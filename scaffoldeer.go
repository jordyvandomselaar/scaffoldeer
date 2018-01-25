package main

import (
    "bytes"
    "gopkg.in/urfave/cli.v1"
    "io"
    "io/ioutil"
    "os"
    "path"
    "strings"
    "log"
    "path/filepath"
)

func main() {
    app := cli.NewApp()
    app.Name = "Scaffoldeer"
    app.Usage = "Scaffold stubs with ease."
    app.Version = "0.2"
    app.Commands = []cli.Command{
        {
            Name:   "make",
            Usage:  "Scaffold a new template",
            Action: scaffold,
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

func scaffold(c *cli.Context) error {
    fields := c.String("fields")

    templateName := c.Args().Get(0)
    scriptPath, err := os.Executable()

    if err != nil {
        return err
    }

    scriptPath = path.Join(scriptPath, "../")

    templatePath := path.Join(scriptPath, "templates", templateName)

    stubsPath := path.Join(templatePath, "stubs")

    stubs, err := getStubs(stubsPath)

    handleError(err)

    parsedFields := strings.Split(fields, ",")
    replacements := make(map[string]string)

    for _, value := range parsedFields {
        parsedFieldsSeparated := strings.Split(value, ":")

        replacements[parsedFieldsSeparated[0]] = parsedFieldsSeparated[1]
    }

    copyStubs(stubs, replacements)

    return nil
}

func copyStubs(stubs []Stub, replacements map[string]string) {
    for _, stub := range stubs {
        stubNewName := parseReplacements(stub.Name, replacements)
        stubNewRelativePath := parseReplacements(stub.RelativePath, replacements)

        stubNewName = strings.Replace(stubNewName, ".stub", "", -1)

        os.MkdirAll(stubNewRelativePath, 0777)
        copyFile(stub.FullPath, path.Join(stubNewRelativePath, stubNewName), replacements)
    }
}

func getStubs(stubsPath string) ([]Stub, error) {
    stubs := []Stub{}

    err := filepath.Walk(stubsPath, func(currentPath string, fileInfo os.FileInfo, err error) error {
        if fileInfo.IsDir() {
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

func copyFile(source string, destination string, replacements map[string]string) error {
    sourceFileContent, err := ioutil.ReadFile(source)
    sourceFileContent = parseFileContent(sourceFileContent, replacements)

    sourceFileContentReader := bytes.NewReader(sourceFileContent)

    destinationFile, err := os.Create(destination)
    defer destinationFile.Close()

    _, err = io.Copy(destinationFile, sourceFileContentReader)

    return err
}

func parseFileContent(content []byte, replacements map[string]string) []byte {
    fileContent := string(content)

    fileContent = parseReplacements(fileContent, replacements)

    return []byte(fileContent)
}

type Stub struct {
    FullPath, RelativePath, Name string
    Content                      []byte
}

func parseReplacements(input string, replacements map[string]string) string {
    for key, value := range replacements {
        input = strings.Replace(input, "__"+key+"__", value, -1)
    }

    return input
}

func handleError(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
