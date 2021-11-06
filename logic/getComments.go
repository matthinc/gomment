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

func (logic *BusinessLogic) GetCommentsTree(threadPath string, parent int, depth int, max int, offset int) (model.CommentsResponse, error) {
	comments, err := logic.DB.GetNewestCommentsByPath(threadPath, max)
	if err != nil {
		return model.CommentsResponse{}, fmt.Errorf("unable to get comments from database: %w", err)
	}

	trees, _ := constructTreeDepthFirst(comments, 0, depth)

	total := len(trees)

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

	return model.CommentsResponse{
		Comments: trees,
		Total:    total,
	}, nil
}

func (logic *BusinessLogic) GetNewestComments(threadPath string, parentId int, maxDepth int, maxCount int) (model.CommentsResponse, error) {
	orderedComments, err := logic.DB.GetNewestCommentsByPath(threadPath, maxCount)
	if err != nil {
		return model.CommentsResponse{}, fmt.Errorf("unable to get comments from database: %w", err)
	}

	subtrees, total := constructTreeDepthFirst(orderedComments, parentId, maxDepth)

	return model.CommentsResponse{
		Comments: subtrees,
		Total:    total,
	}, nil
}
