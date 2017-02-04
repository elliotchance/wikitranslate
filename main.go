package main

import (
    "fmt"
    "regexp"
    "strings"
    "encoding/base64"
    "strconv"
    "net/http"
    "bytes"
    "os"
    "bufio"
    "html"
    "io/ioutil"
)

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, str string, repl func([]string) string) string {
    result := ""
    lastIndex := 0

    for _, v := range re.FindAllSubmatchIndex([]byte(str), -1) {
        groups := []string{}
        for i := 0; i < len(v); i += 2 {
            groups = append(groups, str[v[i]:v[i+1]])
        }

        result += str[lastIndex:v[0]] + repl(groups)
        lastIndex = v[1]
    }

    return result + str[lastIndex:]
}

func BalanceHtmlTags(html string) string {
    parts := []string{}
    result := ""
    split1 := strings.Split(html, "<")
    for _, p := range split1 {
        parts = append(parts, strings.Split(p, ">")...)
    }

    stack := []string{}
    for i := 0; i < len(parts) - 1; i += 2 {
        result += parts[i]

        if /*len(parts[i + 1]) > 0 &&*/ parts[i + 1][0] == '/' {
            for j := len(stack) - 1; j >= 0; j-- {
                s := stack[j]
                stack = stack[:len(stack) - 1]
                result += fmt.Sprintf("</%v>", s)
                if parts[i + 1] == ("/" + s) {
                    break
                }
            }
        } else {
            tagParts := strings.Split(parts[i + 1], " ")
            stack = append(stack, tagParts[0])
            result += fmt.Sprintf("<%v>", parts[i + 1])
        }
    }

    result += parts[len(parts) - 1]

    // Anything left on the stack has to be closed.
    for _, s := range stack {
        result += fmt.Sprintf("</%v>", s)
    }

    return result
}

func HtmlToWiki(html string) string {
    re := regexp.MustCompile(`<img src="(.*?)" options="(.*?)" link="(.*?)">(.*?)</img>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        r := fmt.Sprintf(`[[File:%v`, groups[1])
        if groups[2] != "" || groups[3] != "" {
            r += "|" + groups[2]
        }
        if groups[3] != "" || groups[4] != "" {
            r += "|" + groups[4]

            if groups[3] != "" {
                r += "link=" + groups[3]
            }
        }
        return r + "]]"
    })

    re = regexp.MustCompile(`<a href="(.*?)">(.+?)</a>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        if groups[1] == groups[2] {
            return fmt.Sprintf(`[[%v]]`, groups[1])
        }

        return fmt.Sprintf(`[[%v|%v]]`, groups[1], groups[2])
    })

    html = strings.Replace(html, "<strong><em>", "'''''", -1)
    html = strings.Replace(html, "</strong></em>", "'''''", -1)

    html = strings.Replace(html, "<strong>", "'''", -1)
    html = strings.Replace(html, "</strong>", "'''", -1)
    html = strings.Replace(html, "<em>", "''", -1)
    html = strings.Replace(html, "</em>", "''", -1)

    html = strings.Replace(html, "<li>", "*", -1)
    html = strings.Replace(html, "</li>", "", -1)
    html = strings.Replace(html, "<oli>", "#", -1)
    html = strings.Replace(html, "</oli>", "", -1)

    re = regexp.MustCompile(`<h(.)>(.+?)</h.>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        level, err := strconv.Atoi(groups[1])
        if err != nil {
            panic(err)
        }

        return strings.Repeat("=", level) + groups[2] + strings.Repeat("=", level)
    })

    re = regexp.MustCompile(`<template name="(.+?)">(.*?)</template>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        if groups[2] == "" {
            return fmt.Sprintf(`{{%v}}`, groups[1])
        }

        re = regexp.MustCompile(`<arg name="(.*?)">(.*?)</arg>`)
        result := ReplaceAllStringSubmatchFunc(re, groups[2], func(groups []string) string {
            if groups[1] == "" {
                return fmt.Sprintf(`|%v`, groups[2])
            }

            return fmt.Sprintf(`|%v=%v`, groups[1], groups[2])
        })

        return fmt.Sprintf(`{{%v%v}}`, groups[1], result)
    })

    re = regexp.MustCompile(`(?s)<table(.*?)>(.*?)</table>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        result := "{|"

        re = regexp.MustCompile(`(?s)<tr(.*?)>\n(.*?)\n</tr>`)
        result += ReplaceAllStringSubmatchFunc(re, groups[2], func(groups []string) string {
            re = regexp.MustCompile(`(?s)<t([dh])(.*?)>(.*?)</t[dh]>`)
            return "|-" + groups[1] + "\n" + ReplaceAllStringSubmatchFunc(re, groups[2], func(groups []string) string {
                if groups[1] == "d" {
                    return "|" + groups[3]
                }
                return "!" + groups[3]
            })
        })

        return result + "|}"
    })

    re = regexp.MustCompile(`<ref data="(.*?)"(.*?)></ref>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        decoded, err := base64.StdEncoding.DecodeString(groups[1])
        if err != nil {
            panic(err)
        }

        return fmt.Sprintf(`<ref%v>%v</ref>`, groups[2], string(decoded))
    })

    re = regexp.MustCompile(`<nowiki data="(.*?)"(.*?)></nowiki>`)
    html = ReplaceAllStringSubmatchFunc(re, html, func(groups []string) string {
        decoded, err := base64.StdEncoding.DecodeString(groups[1])
        if err != nil {
            panic(err)
        }

        return fmt.Sprintf(`<nowiki%v>%v</nowiki>`, groups[2], string(decoded))
    })

    return html
}

func WikiToHtml(wikimarkup string) string {
    re := regexp.MustCompile(`<nowiki(.*?)>(.*?)</nowiki>`)
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        encoded := base64.StdEncoding.EncodeToString([]byte(groups[2]))
        return fmt.Sprintf(`<nowiki data="%v"%v></nowiki>`, encoded, groups[1])
    })

    re = regexp.MustCompile(`<ref(.*?)>(.*?)</ref>`)
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        encoded := base64.StdEncoding.EncodeToString([]byte(groups[2]))
        return fmt.Sprintf(`<ref data="%v"%v></ref>`, encoded, groups[1])
    })

    re = regexp.MustCompile("{{([^|}]+)\\|?(.*?)}}")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        r := `<template name="` + groups[1] + `">`

        if groups[2] != "" {
            params := strings.Split(groups[2], "|")
            for _, param := range params {
                if strings.Contains(param, "=") {
                    kv := strings.Split(param, "=")
                    r += fmt.Sprintf(`<arg name="%v">%v</arg>`, kv[0], kv[1])
                } else {
                    r += fmt.Sprintf(`<arg name="">%v</arg>`, param)
                }
            }
        }

        return r + "</template>"
    })

    re = regexp.MustCompile("'''(.+?)'''")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        return "<strong>" + groups[1] + "</strong>"
    })

    re = regexp.MustCompile("''(.+?)''")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        return "<em>" + groups[1] + "</em>"
    })

    re = regexp.MustCompile("\\[\\[(.+?)\\]\\]")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        // File:
        if strings.HasPrefix(groups[1], "File:") {
            parts := strings.SplitN(groups[1], "|", 3)
            if len(parts) == 1 {
                parts = append(parts, "", "")
            } else if len(parts) == 2 {
                parts = append(parts, "")
            }

            link := ""
            if strings.HasPrefix(parts[2], "link=") {
                link = parts[2][5:]
                parts[2] = ""
            }

            return fmt.Sprintf(`<img src="%v" options="%v" link="%v">%v</a>`, parts[0][5:], parts[1], link, parts[2])
        }
        
        // Else
        parts := strings.SplitN(groups[1], "|", 2)
        if len(parts) == 1 {
            return fmt.Sprintf(`<a href="%v">%v</a>`, parts[0], parts[0])
        } else {
            return fmt.Sprintf(`<a href="%v">%v</a>`, parts[0], parts[1])
        }
    })

    re = regexp.MustCompile("\\[(.{10,}?)\\]")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        parts := strings.SplitN(groups[1], " ", 2)
        if len(parts) == 1 {
            return "<a>" + groups[1] + "</a>"
        } else {
            return fmt.Sprintf(`<a href="%v">%v</a>`, parts[0], parts[1])
        }
    })

    // Headings
    for i := 6; i >= 1; i-- {
        re = regexp.MustCompile("(^|\\s)" + strings.Repeat("=", i) + "(.+?)" + strings.Repeat("=", i))
        wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
            return fmt.Sprintf("%v<h%v>%v</h%v>", groups[1], i, groups[2], i)
        })
    }

    // Bullet points/lists
    re = regexp.MustCompile("(?m)^([*#])([^\\n]+)")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(groups []string) string {
        if groups[1] == "*" {
            return fmt.Sprintf("<li>%v</li>", groups[2])
        }

        return fmt.Sprintf("<oli>%v</oli>", groups[2])
    })

    // Images
    // Raw URLs
    
    // Tables
    re = regexp.MustCompile("(?s){\\|([^\\n]*)(.*?)\\|}")
    wikimarkup = ReplaceAllStringSubmatchFunc(re, wikimarkup, func(table_groups []string) string {
        table := "<table " + table_groups[1] + ">\n"

        lines := strings.Split(table_groups[2], "\n")

        inRow := false
        printedTr := false
        for _, line := range lines {
            if strings.HasPrefix(line, "|-") {
                printedTr = true
                if inRow {
                    table += "</tr>\n"
                }
                table += "<tr " + line[2:] + ">\n"
                inRow = true
                continue
            }

            if strings.HasPrefix(line, "|") {
                if !printedTr {
                    printedTr = true
                    inRow = true
                    table += "<tr>\n"
                }

                parts := strings.Split(line[1:], "|")
                style := ""
                body := ""

                if len(parts) > 1 {
                    style = parts[0]
                    body = parts[1]
                } else {
                    body = parts[0]
                }

                bodyParts := strings.Split(body, "!!")

                for _, bodyPart := range bodyParts {
                    table += "<td " + style + ">" + bodyPart + "</td>\n"
                }

                continue
            }

            if strings.HasPrefix(line, "!") {
                if !printedTr {
                    printedTr = true
                    inRow = true
                    table += "<tr>\n"
                }

                parts := strings.Split(line[1:], "|")
                style := ""
                body := ""

                if len(parts) > 1 {
                    style = parts[0]
                    body = parts[1]
                } else {
                    body = parts[0]
                }

                bodyParts := strings.Split(body, "!!")

                for _, bodyPart := range bodyParts {
                    table += "<th " + style + ">" + bodyPart + "</th>\n"
                }

                continue
            }
        }
        table += "</tr>\n"

        return table + "</table>"
    })

    // Random unbalanced left overs
    wikimarkup = strings.Replace(wikimarkup, "'''", "<strong>", -1)
    wikimarkup = strings.Replace(wikimarkup, "''", "<em>", -1)

    return BalanceHtmlTags(wikimarkup)
}

func main() {
    fmt.Printf("Wiki URL or File: ")
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')

    if input[0] == 'h' {
        fmt.Printf("Downloading page...")

        tokens := strings.Split(input, "/")
        title := strings.TrimSpace(tokens[len(tokens) - 1])

        url := "https://en.wikipedia.org/w/index.php?title=" + title + "&action=edit"
        response, err := http.Get(url)
        if err != nil {
            panic(err)
        }
        defer response.Body.Close()

        buf := new(bytes.Buffer)
        buf.ReadFrom(response.Body)
        content := buf.String()

        re := regexp.MustCompile("(?s)<textarea.*?>(.*)</textarea>")
        wikimarkup := html.UnescapeString(re.FindStringSubmatch(content)[1])

        fmt.Printf(" Done\n")

        fileHandle, _ := os.Create(title + ".html")
        writer := bufio.NewWriter(fileHandle)
        defer fileHandle.Close()

        writer.WriteString(WikiToHtml(wikimarkup))

        writer.Flush()
    } else {
        fileName := strings.TrimSpace(input)
        html, err := ioutil.ReadFile(fileName)
        if err != nil {
            panic(err)
        }

        fileHandle, _ := os.Create(fileName + ".txt")
        writer := bufio.NewWriter(fileHandle)
        defer fileHandle.Close()

        writer.WriteString(HtmlToWiki(string(html)))

        writer.Flush()

        fmt.Printf("Done\n")
    }
}