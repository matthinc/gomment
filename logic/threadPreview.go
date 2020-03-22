package logic

import (
    "github.com/matthinc/gomment/model"
    "strings"
    "strconv"
)

var previewColors =  [...] string {
    "#ffbdbd",
    "#ffe7a3",
    "#c6ffa3",
    "#bdfeff",
    "#f0cfff",
}

func generateTreeHtmlPreview(tree model.CommentTree, sb *strings.Builder, depth int) {
    comment := tree.Comment
    sb.WriteString(`<div style="border: 1px solid black;padding:5px;background-color:`)
    sb.WriteString(previewColors[depth % len(previewColors)])
    sb.WriteString(`">`)
    // <id>
    sb.WriteString(`<div style="font-style:italic">ID:&nbsp;`)
    sb.WriteString(strconv.Itoa(comment.Id))
    sb.WriteString(`</div>`)
    // </id>
    // <author>
    sb.WriteString(`<div style="font-weight: bold;padding-bottom: 3px">`)
    sb.WriteString(comment.Author)
    sb.WriteString(`</div>`)
    // </author>
    sb.WriteString(comment.Text)
    // <children>
    sb.WriteString(`<div style="padding-left: 10px">`)
    for _, c := range tree.Children {
        generateTreeHtmlPreview(c, sb, depth + 1)
    }
    sb.WriteString(`</div>`)
    // </children>
    sb.WriteString(`</div>`)
}

func (logic* BusinessLogic) GenerateHTMLThreadPreview(thread int) string {
    tree := logic.GetCommentsTree(thread)
    
    var sb strings.Builder
    sb.WriteString("<h1>Thread Preview</h1>")

    for _, topLevelNode := range tree {
        generateTreeHtmlPreview(topLevelNode, &sb, 0)
    }
    
    return sb.String()
}
