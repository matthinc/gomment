package logic

import (
    "github.com/matthinc/gomment/model"
)

func (logic* BusinessLogic) GetComments(thread int) []model.Comment {
    comments := logic.DB.QueryComments(thread)
    return comments
}

func commentHasChildren(comments []model.Comment, id int) bool {
    for _, comment := range comments {
        if comment.ParentId == id {
            return true
        }
    }
    return false
}

func commentsToTree(comments []model.Comment, parent int, depthLeft int) []model.CommentTree {
    tree := make([]model.CommentTree, 0)

    for _, comment := range comments {
        if comment.ParentId == parent {
            var children []model.CommentTree
            var hasChildren bool

            if depthLeft > 0 {
                children = commentsToTree(comments, comment.Id, depthLeft - 1)
                hasChildren = len(children) > 0
            } else {
                children = nil
                hasChildren = commentHasChildren(comments, comment.Id)
            }

            tree = append(tree,
                model.CommentTree { Comment: comment, Children: children, HasChildren: hasChildren })
        }
    }

    return tree
}

func (logic* BusinessLogic) GetCommentsTree(thread int, depth int, max int, offset int, tree int) []model.CommentTree {
    comments := logic.GetComments(thread)
    trees := commentsToTree(comments, 0, depth)

    // Offset
    if len(trees) > offset {
        trees = trees[offset:]
    } else {
        trees = []model.CommentTree{}
    }

    // Max
    if len(trees) > max && max > 0 {
        trees = trees[:max]
    }

    return trees
}
