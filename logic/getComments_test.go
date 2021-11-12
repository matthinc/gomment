package logic_test

import (
	"testing"

	"github.com/matthinc/gomment/logic"
	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDb struct {
	comments  []model.Comment
	creations []model.CommentCreation
}

func (db *MockDb) Open(path string) error { return nil }
func (db *MockDb) Setup() error           { return nil }
func (db *MockDb) Close()                 {}
func (db *MockDb) CreateComment(commentCreation *model.CommentCreation, createdAt int64) (int64, error) {
	db.creations = append(db.creations, *commentCreation)
	return 0, nil
}
func (db *MockDb) GetCommentsNbf(path string, limit int) ([]model.Comment, persistence.ThreadMetaInfo, error) {
	numRoot := 0

	for _, comment := range db.comments {
		if comment.ParentId == 0 {
			numRoot = numRoot + 1
		}
	}

	return db.comments, persistence.ThreadMetaInfo{
		NumTotal: len(db.comments),
		NumRoot:  numRoot,
	}, nil
}
func (db *MockDb) GetMoreCommentsNbf(threadId int64, parentId int64, newestCreatedAt int64, excludeIds []int64, limit int) ([]model.Comment, error) {
	return db.comments, nil
}
func (db *MockDb) GetThreads() ([]model.Thread, error) { return []model.Thread{}, nil }

func TestSimple(t *testing.T) {
	db := MockDb{
		comments: []model.Comment{{
			Id:          1,
			Author:      "gandalf",
			Email:       "gandi@dalfi.com",
			Text:        "Power to me",
			ParentId:    0,
			CreatedAt:   1,
			TouchedAt:   1,
			NumChildren: 0,
		}},
	}
	sut := logic.Create(&db, "")

	commentResponse, err := sut.GetCommentsNbf("", 0, 99, 99)
	require.NoError(t, err)

	assert.Equal(t, 1, commentResponse.NumTotal, "expected the total number of comments to be 1")
	require.Equal(t, 1, len(commentResponse.Comments), "expected the tree to have 1 root comment")

	firstTree := commentResponse.Comments[0]

	assert.Equal(t, 1, firstTree.Comment.Id, "expected the first root comment to have the id 1")
	assert.Equal(t, 0, len(firstTree.Children), "expected the first root comment to not have any children")
}

func TestTwoRootComments(t *testing.T) {
	db := MockDb{
		comments: []model.Comment{{
			Id:          1,
			Author:      "gandalf",
			Email:       "gandi@dalfi.com",
			Text:        "Power to me",
			ParentId:    0,
			CreatedAt:   1,
			TouchedAt:   1,
			NumChildren: 0,
		}, {
			Id:          2,
			Author:      "peter",
			Email:       "peter@sad.de",
			Text:        "I am a spammer",
			ParentId:    0,
			CreatedAt:   2,
			TouchedAt:   2,
			NumChildren: 0,
		}},
	}
	sut := logic.Create(&db, "")

	commentResponse, err := sut.GetCommentsNbf("", 0, 99, 99)
	require.NoError(t, err)

	assert.Equal(t, 2, commentResponse.NumTotal, "expected the total number of comments to be 2")
	require.Equal(t, 2, len(commentResponse.Comments), "expected the tree to have 2 root comments")

	firstTree := commentResponse.Comments[0]
	secondTree := commentResponse.Comments[1]

	assert.Equal(t, 1, firstTree.Comment.Id, "expected the first comment to have the id 1")
	assert.Equal(t, 2, secondTree.Comment.Id, "expected the second comment to have the id 2")
	assert.Equal(t, 0, len(firstTree.Children), "expected the first root comment to not have any children")
	assert.Equal(t, 0, len(secondTree.Children), "expected the second root comment to not have any children")
}

func TestTwoChainedComments(t *testing.T) {
	db := MockDb{
		comments: []model.Comment{{
			Id:          1,
			Author:      "gandalf",
			Email:       "gandi@dalfi.com",
			Text:        "Power to me",
			ParentId:    0,
			CreatedAt:   1,
			TouchedAt:   1,
			NumChildren: 1,
		}, {
			Id:          2,
			Author:      "peter",
			Email:       "peter@sad.de",
			Text:        "I am a spammer",
			ParentId:    1,
			CreatedAt:   2,
			TouchedAt:   2,
			NumChildren: 0,
		}},
	}
	sut := logic.Create(&db, "")

	commentResponse, err := sut.GetCommentsNbf("", 0, 99, 99)
	require.NoError(t, err)

	assert.Equal(t, 2, commentResponse.NumTotal, "expected the total number of comments to be 2")

	require.Equal(t, 1, len(commentResponse.Comments), "expected the tree to have 1 root comment")
	firstTree := commentResponse.Comments[0]

	assert.Equal(t, 1, firstTree.Comment.Id, "expected the first root comment to have the id 1")

	require.Equal(t, 1, len(firstTree.Children), "expected the first root comment to have 1 child")
	firstChild := firstTree.Children[0]

	assert.Equal(t, 2, firstChild.Comment.Id, "expected first child to have the id 2")
	assert.Equal(t, 0, len(firstChild.Children), "expected the first child to not have any children")
}

func TestTwoChains(t *testing.T) {
	db := MockDb{
		comments: []model.Comment{{
			Id:          1,
			Author:      "gandalf",
			Email:       "gandi@dalfi.com",
			Text:        "Power to me",
			ParentId:    0,
			CreatedAt:   1,
			TouchedAt:   1,
			NumChildren: 1,
		}, {
			Id:          2,
			Author:      "peter",
			Email:       "peter@sad.de",
			Text:        "I am a spammer",
			ParentId:    1,
			CreatedAt:   3,
			TouchedAt:   3,
			NumChildren: 0,
		}, {
			Id:          3,
			Author:      "",
			Email:       "",
			Text:        "",
			ParentId:    0,
			CreatedAt:   2,
			TouchedAt:   2,
			NumChildren: 1,
		}, {
			Id:          4,
			Author:      "",
			Email:       "",
			Text:        "",
			ParentId:    3,
			CreatedAt:   4,
			TouchedAt:   4,
			NumChildren: 0,
		}},
	}
	sut := logic.Create(&db, "")

	commentResponse, err := sut.GetCommentsNbf("", 0, 99, 99)
	require.NoError(t, err)

	assert.Equal(t, 4, commentResponse.NumTotal, "expected the total number of comments to be 4")

	require.Equal(t, 2, len(commentResponse.Comments), "expected the tree to have 2 root comments")
	firstTree := commentResponse.Comments[0]
	secondTree := commentResponse.Comments[1]

	assert.Equal(t, 1, firstTree.Comment.Id, "expected the first root comment to have the id 1")
	assert.Equal(t, 3, secondTree.Comment.Id, "expected the second root comment to have the id 3")

	require.Equal(t, 1, len(firstTree.Children), "expected the first root comment to have 1 child")
	require.Equal(t, 1, len(secondTree.Children), "expected the second root comment to have 1 child")
	firstChild := firstTree.Children[0]
	secondChild := secondTree.Children[0]

	assert.Equal(t, 2, firstChild.Comment.Id, "expected first child to have the id 2")
	assert.Equal(t, 0, len(firstChild.Children), "expected the first child to not have any children")
	assert.Equal(t, 4, secondChild.Comment.Id, "expected second child to have the id 4")
	assert.Equal(t, 0, len(secondChild.Children), "expected the second child to not have any children")
}

func TestLeafyChain(t *testing.T) {
	db := MockDb{
		comments: []model.Comment{{
			Id:       1,
			ParentId: 0,
		}, {
			Id:       2,
			ParentId: 1,
		}, {
			Id:       3,
			ParentId: 2,
		}, {
			Id:       4,
			ParentId: 2,
		}},
	}
	sut := logic.Create(&db, "")

	commentResponse, err := sut.GetCommentsNbf("", 0, 99, 99)
	require.NoError(t, err)

	assert.Equal(t, 4, commentResponse.NumTotal, "expected the total number of comments to be 4")

	// require.Equal(t, 2, len(commentResponse.Comments), "expected the tree to have 2 root comments")
	// firstTree := commentResponse.Comments[0]
	// secondTree := commentResponse.Comments[1]

	// assert.Equal(t, 1, firstTree.Comment.Id, "expected the first root comment to have the id 1")
	// assert.Equal(t, 3, secondTree.Comment.Id, "expected the second root comment to have the id 3")

	// require.Equal(t, 1, len(firstTree.Children), "expected the first root comment to have 1 child")
	// require.Equal(t, 1, len(secondTree.Children), "expected the second root comment to have 1 child")
	// firstChild := firstTree.Children[0]
	// secondChild := secondTree.Children[0]

	// assert.Equal(t, 2, firstChild.Comment.Id, "expected first child to have the id 2")
	// assert.Equal(t, 0, len(firstChild.Children), "expected the first child to not have any children")
	// assert.Equal(t, 4, secondChild.Comment.Id, "expected second child to have the id 4")
	// assert.Equal(t, 0, len(secondChild.Children), "expected the second child to not have any children")
}
