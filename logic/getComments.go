package logic

import (
	"fmt"

	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence"
)

type orderType int

const (
	orderNbf orderType = iota
	orderNsf orderType = iota
	orderOsf orderType = iota
)

func constructTreeDepthFirst(comments []model.Comment, parentId int64, depthLeft int) []CommentTree {
	subtrees := make([]CommentTree, 0)

	for _, comment := range comments {
		if comment.ParentId == parentId {
			var children []CommentTree

			if depthLeft > 0 {
				children = constructTreeDepthFirst(comments, comment.Id, depthLeft-1)
			} else {
				children = []CommentTree{}
			}

			subtrees = append(subtrees, CommentTree{Comment: comment, Children: children})
		}
	}

	return subtrees
}

func (logic *BusinessLogic) getComments(order orderType, threadPath string, parentId int64, maxDepth int, maxCount int) (CommentResult, error) {
	var (
		orderedComments []model.Comment
		metadata        persistence.ThreadMetaInfo
		err             error
	)

	switch order {
	case orderNbf:
		orderedComments, metadata, err = logic.DB.GetCommentsNbf(threadPath, maxDepth, maxCount)
	case orderNsf:
		orderedComments, metadata, err = logic.DB.GetCommentsNsf(threadPath, maxDepth, maxCount)
	case orderOsf:
		orderedComments, metadata, err = logic.DB.GetCommentsOsf(threadPath, maxDepth, maxCount)
	}
	if err != nil {
		return CommentResult{}, fmt.Errorf("unable to get comments from database: %w", err)
	}

	subtrees := constructTreeDepthFirst(orderedComments, parentId, maxDepth)

	return CommentResult{
		Comments:       subtrees,
		NumRoot:        metadata.NumRoot,
		NumTotal:       metadata.NumTotal,
		NumRootPayload: len(subtrees),
		ThreadId:       metadata.ThreadId,
	}, nil
}

func (logic *BusinessLogic) GetCommentsNbf(threadPath string, parentId int64, maxDepth int, maxCount int) (CommentResult, error) {
	return logic.getComments(orderNbf, threadPath, parentId, maxDepth, maxCount)
}

func (logic *BusinessLogic) GetCommentsNsf(threadPath string, parentId int64, maxDepth int, maxCount int) (CommentResult, error) {
	return logic.getComments(orderNsf, threadPath, parentId, maxDepth, maxCount)
}

func (logic *BusinessLogic) GetCommentsOsf(threadPath string, parentId int64, maxDepth int, maxCount int) (CommentResult, error) {
	return logic.getComments(orderOsf, threadPath, parentId, maxDepth, maxCount)
}

func (logic *BusinessLogic) getMoreComments(order orderType, threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error) {
	var (
		orderedComments []model.Comment
		err             error
	)
	switch order {
	case orderNbf:
		orderedComments, err = logic.DB.GetMoreCommentsNbf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	case orderNsf:
		orderedComments, err = logic.DB.GetMoreCommentsNsf(threadId, parentId, newestCreatedAt, limit)
	case orderOsf:
		orderedComments, err = logic.DB.GetMoreCommentsOsf(threadId, parentId, newestCreatedAt, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to get comments from database: %w", err)
	}

	return orderedComments, nil
}

func (logic *BusinessLogic) GetMoreCommentsNbf(threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error) {
	return logic.getMoreComments(orderNbf, threadId, parentId, newestCreatedAt, excludeIds, limit)
}

func (logic *BusinessLogic) GetMoreCommentsNsf(threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error) {
	return logic.getMoreComments(orderNsf, threadId, parentId, newestCreatedAt, excludeIds, limit)
}

func (logic *BusinessLogic) GetMoreCommentsOsf(threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error) {
	return logic.getMoreComments(orderOsf, threadId, parentId, newestCreatedAt, excludeIds, limit)
}
