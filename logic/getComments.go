package logic

import (
	"fmt"

	"github.com/matthinc/gomment/model"
)

func constructTreeDepthFirst(comments []model.Comment, parentId int, depthLeft int) []model.CommentTree {
	subtrees := make([]model.CommentTree, 0)

	for _, comment := range comments {
		if comment.ParentId == parentId {
			var children []model.CommentTree

			if depthLeft > 0 {
				children = constructTreeDepthFirst(comments, comment.Id, depthLeft-1)
			} else {
				children = []model.CommentTree{}
			}

			subtrees = append(subtrees, model.CommentTree{Comment: comment, Children: children})
		}
	}

	return subtrees
}

func (logic *BusinessLogic) GetNewestComments(threadPath string, parentId int, maxDepth int, maxCount int) (model.CommentsResponse, error) {
	orderedComments, metadata, err := logic.DB.GetNewestCommentsByPath(threadPath, maxCount)
	if err != nil {
		return model.CommentsResponse{}, fmt.Errorf("unable to get comments from database: %w", err)
	}

	subtrees := constructTreeDepthFirst(orderedComments, parentId, maxDepth)

	return model.CommentsResponse{
		Comments:       subtrees,
		NumRoot:        metadata.NumRoot,
		NumTotal:       metadata.NumTotal,
		NumRootPayload: len(subtrees),
	}, nil
}
