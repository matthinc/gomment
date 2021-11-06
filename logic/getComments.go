package logic

import (
	"fmt"

	"github.com/matthinc/gomment/model"
)

func constructTreeDepthFirst(comments []model.Comment, parentId int, depthLeft int) ([]model.CommentTree, int) {
	subtrees := make([]model.CommentTree, 0)
	total := 0

	for _, comment := range comments {
		if comment.ParentId == parentId {
			var children []model.CommentTree

			if depthLeft > 0 {
				subtotal := 0
				children, subtotal = constructTreeDepthFirst(comments, comment.Id, depthLeft-1)
				total = total + subtotal
			} else {
				children = []model.CommentTree{}
			}

			subtrees = append(subtrees, model.CommentTree{Comment: comment, Children: children})
		}
	}

	return subtrees, total + len(subtrees)
}

func (logic *BusinessLogic) GetNewestComments(threadPath string, parentId int, maxDepth int, maxCount int) (model.CommentsResponse, error) {
	orderedComments, err := logic.DB.GetNewestCommentsByPath(threadPath, maxCount)
	if err != nil {
		return model.CommentsResponse{}, fmt.Errorf("unable to get comments from database: %w", err)
	}

	subtrees, total := constructTreeDepthFirst(orderedComments, parentId, maxDepth)

	return model.CommentsResponse{
		Comments:    subtrees,
		NumChildren: len(subtrees),
		Total:       total,
	}, nil
}
