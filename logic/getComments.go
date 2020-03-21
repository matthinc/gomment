package logic

import (
    "github.com/matthinc/gomment/model"
)

func (logic* BusinessLogic) GetComments(thread int) []*model.Comment {
    comments := logic.DB.QueryComments(thread)
    return comments
}

func commentsToTree(comments []*model.Comment, parent int) []*model.CommentTree {
    var tree []*model.CommentTree
    
    for _, comment := range comments {
        if comment.ParentId == parent {
            tree = append(
                tree,
                &model.CommentTree {
                    Comment: comment,
                    Children: commentsToTree(comments, comment.Id),
                })
        }
    }
    
    return tree
}

func (logic* BusinessLogic) GetCommentsTree(thread int) []*model.CommentTree {
    comments := logic.GetComments(thread)
    return commentsToTree(comments, -1)
}

