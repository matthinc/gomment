package sqlite_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// integration test suite for comment creation

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func createTestDatabase(testName string) (*sqlite.DB, func(), error) {
	db := sqlite.New()

	dbFilename := "integration-test-" + testName + ".sqlite"
	os.Remove(dbFilename)

	err := db.Open(dbFilename)
	if err != nil {
		return nil, nil, err
	}

	return &db, func() {
		os.Remove(dbFilename)
	}, nil
}

func assertDepth(t *testing.T, db *sqlite.DB, commentId int, depth int) {
	commentRow, err := db.GetCommentRow(int64(commentId))
	require.NoError(t, err)
	assert.Equal(t, depth, commentRow.DepthLevel, fmt.Sprintf("expected comment id %d to have depth %d, was %d", commentId, depth, commentRow.DepthLevel))
}

func TestNbfRootComment(t *testing.T) {
	db, deleter, err := createTestDatabase("01")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-01",
		ParentId:   0,
	}, 0)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	comments, metainfo, err := db.GetCommentsNbf("/test-01", 100, 100)
	require.NoError(t, err)

	assert.Equal(t, 1, metainfo.NumTotal, "expected the thread to have 1 total comment")
	assert.Equal(t, 1, metainfo.NumRoot, "expected the thread to have 1 root comment")
	require.Equal(t, 1, len(comments), "expected 1 comment to be returned")

	assert.Equal(t, comments[0].TouchedAt, comments[0].CreatedAt, "expected the touched_at time to be equal to created_at")
	assert.Zero(t, comments[0].NumChildren)

	assertDepth(t, db, comments[0].Id, 0)
}

func TestCreateNonExistingParent(t *testing.T) {
	db, deleter, err := createTestDatabase("02")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-02",
		ParentId:   999,
	}, 0)

	assert.Zero(t, commentId)
	require.Error(t, err)
}

func TestNbfEmptyThread(t *testing.T) {
	db, deleter, err := createTestDatabase("05")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-05",
		ParentId:   0,
	}, 0)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	comments, metainfo, err := db.GetCommentsNbf("/foobar", 100, 100)
	require.NoError(t, err, "expected no error even if the path does not exist yet")

	assert.Equal(t, 0, metainfo.NumTotal, "expected the thread to have 0 total comments")
	assert.Equal(t, 0, metainfo.NumRoot, "expected the thread to have 0 root comments")

	assert.Zero(t, len(comments), "expected comment list to be empty for non-existant path")
}

func TestNbfChildComment(t *testing.T) {
	db, deleter, err := createTestDatabase("03")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-03",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	commentId, err = db.CreateComment(&model.CommentCreation{
		Author:     "Alfred Peterson",
		Email:      "alfred@peterson.se",
		Text:       "I disagree!",
		ThreadPath: "/test-03",
		ParentId:   1,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	comments, metainfo, err := db.GetCommentsNbf("/test-03", 100, 100)
	require.NoError(t, err)

	assert.Equal(t, 2, metainfo.NumTotal, "expected the thread to have 2 total comments")
	assert.Equal(t, 1, metainfo.NumRoot, "expected the thread to have 1 root comment")

	require.Equal(t, 2, len(comments), "expected 2 comment to be in the database")

	// leaf comment is at the bottom of the list due to intra-branch ordering
	leafComment := comments[1]
	rootComment := comments[0]

	assert.Equal(t, leafComment.TouchedAt, leafComment.CreatedAt, "expected the touched_at time to be equal to created_at in the leaf comment")
	assert.NotEqual(t, rootComment.TouchedAt, rootComment.CreatedAt, "expected the touched_at time not to be equal to created_at in the root comment")

	assert.GreaterOrEqual(t, leafComment.CreatedAt, rootComment.CreatedAt, "expected the creation time of the most recent comment to be smaller than the one of the second one")

	assert.Equal(t, 1, rootComment.NumChildren, "expected root comment to have one child")
	assert.Equal(t, 0, leafComment.NumChildren, "expected leaf comment to have no children")

	assertDepth(t, db, leafComment.Id, 1)
	assertDepth(t, db, rootComment.Id, 0)
}

func TestNbfTwoChildComments(t *testing.T) {
	db, deleter, err := createTestDatabase("04")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-04",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	commentId, err = db.CreateComment(&model.CommentCreation{
		Author:     "Alfred Peterson",
		Email:      "alfred@peterson.se",
		Text:       "I disagree!",
		ThreadPath: "/test-04",
		ParentId:   int(commentId),
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	commentId, err = db.CreateComment(&model.CommentCreation{
		Author:     "Child Childison",
		Email:      "child@childison.dk",
		Text:       "I am the child.",
		ThreadPath: "/test-04",
		ParentId:   int(commentId),
	}, 3)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	comments, metainfo, err := db.GetCommentsNbf("/test-04", 100, 100)
	require.NoError(t, err)

	assert.Equal(t, 3, metainfo.NumTotal, "expected the thread to have 3 total comments")
	assert.Equal(t, 1, metainfo.NumRoot, "expected the thread to have 1 root comment")

	require.Equal(t, 3, len(comments), "expected 3 comment to be in the database")

	// leaf comment is at the bottom of the list due to intra-branch ordering
	leafComment := comments[2]
	midComment := comments[1]
	rootComment := comments[0]

	assert.Equal(t, leafComment.TouchedAt, leafComment.CreatedAt, "expected the touched_at time to be equal to created_at in the leaf comment")
	assert.NotEqual(t, midComment.TouchedAt, midComment.CreatedAt, "expected the touched_at time not to be equal to created_at in the mid comment")
	assert.NotEqual(t, rootComment.TouchedAt, rootComment.CreatedAt, "expected the touched_at time not to be equal to created_at in the root comment")

	assert.GreaterOrEqual(t, midComment.CreatedAt, rootComment.CreatedAt, "expected creationtime(midComment) >= creationtime(rootComment)")
	assert.GreaterOrEqual(t, leafComment.CreatedAt, midComment.CreatedAt, "expected creationtime(leafComment) >= creationtime(midComment)")
	assert.GreaterOrEqual(t, leafComment.CreatedAt, rootComment.CreatedAt, "expected creationtime(leafComment) >= creationtime(rootComment)")

	assert.Equal(t, 1, rootComment.NumChildren, "expected root comment to have one child")
	assert.Equal(t, 1, midComment.NumChildren, "expected mid comment to have one child")
	assert.Equal(t, 0, leafComment.NumChildren, "expected leaf comment to have no children")

	assertDepth(t, db, rootComment.Id, 0)
	assertDepth(t, db, midComment.Id, 1)
	assertDepth(t, db, leafComment.Id, 2)
}

func TestAvailableThreads(t *testing.T) {
	db, deleter, err := createTestDatabase("06")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-06-a",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	commentId, err = db.CreateComment(&model.CommentCreation{
		Author:     "Alfred Peterson",
		Email:      "alfred@peterson.se",
		Text:       "I disagree!",
		ThreadPath: "/test-06-a",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	commentId, err = db.CreateComment(&model.CommentCreation{
		Author:     "Child Childison",
		Email:      "child@childison.dk",
		Text:       "I am the child.",
		ThreadPath: "/test-06-b",
		ParentId:   0,
	}, 3)
	require.NoError(t, err)

	threads, err := db.GetThreads()
	require.NoError(t, err)

	assert.Equal(t, 2, len(threads), "expected two threads to exist")
	assert.NotZero(t, threads[0].Id, "expected thread id to be non-zero")
	assert.NotZero(t, threads[1].Id, "expected thread id to be non-zero")
	assert.NotEqual(t, threads[0].Id, threads[1].Id, "expected thread ids to be different")

	assert.ElementsMatch(t, []string{"/test-06-a", "/test-06-b"}, []string{threads[0].Path, threads[1].Path}, "expected two thread names to match")
}

func TestNbfNewestLimit(t *testing.T) {
	db, deleter, err := createTestDatabase("07")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	commentId, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-07",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	commentId, err = db.CreateComment(&model.CommentCreation{
		Author:     "Alfred Peterson",
		Email:      "alfred@peterson.se",
		Text:       "I disagree!",
		ThreadPath: "/test-07",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, commentId, "expected comment id to be not zero")

	for i := range [4]int{} {
		comments, metainfo, err := db.GetCommentsNbf("/test-07", 100, i)
		assert.NoError(t, err)

		assert.Equal(t, 2, metainfo.NumTotal, "expected the thread to have 2 total comments")
		assert.Equal(t, 2, metainfo.NumRoot, "expected the thread to have 2 root comment")

		assert.Equal(t, min(i, 2), len(comments))
	}
}

// create two branches with one child each, get latest two comments
func TestNbfComplexTwoBranches(t *testing.T) {
	db, deleter, err := createTestDatabase("08")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		Author:     "Peter Müller",
		Email:      "peter@mueller.de",
		Text:       "This is a great integration test.",
		ThreadPath: "/test-08",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	rootComment2, err := db.CreateComment(&model.CommentCreation{
		Author:     "Alfred Peterson",
		Email:      "alfred@peterson.se",
		Text:       "I disagree!",
		ThreadPath: "/test-08",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, rootComment2, "expected comment id to be not zero")

	leafComment1, err := db.CreateComment(&model.CommentCreation{
		Author:     "Müller's Child",
		Email:      "child@mueller.de",
		Text:       "This is a great child.",
		ThreadPath: "/test-08",
		ParentId:   int(rootComment1),
	}, 3)
	require.NoError(t, err)
	assert.NotZero(t, leafComment1, "expected comment id to be not zero")

	leafComment2, err := db.CreateComment(&model.CommentCreation{
		Author:     "Alfred's Child",
		Email:      "child@peterson.se",
		Text:       "No, dad!",
		ThreadPath: "/test-08",
		ParentId:   int(rootComment2),
	}, 4)
	require.NoError(t, err)
	assert.NotZero(t, leafComment2, "expected comment id to be not zero")

	// expected order (most recent first)
	orderedIds := []int64{rootComment2, leafComment2, rootComment1, leafComment1}

	for i := range [6]int{} {
		comments, metainfo, err := db.GetCommentsNbf("/test-08", 100, i)
		require.NoError(t, err)

		assert.Equal(t, 4, metainfo.NumTotal, "expected the thread to have 4 total comments")
		assert.Equal(t, 2, metainfo.NumRoot, "expected the thread to have 2 root comments")

		require.Equal(t, min(4, i), len(comments))

		for idx, comment := range comments {
			assert.Equal(t, orderedIds[idx], int64(comment.Id), "expected id to be in the right order")
		}
	}
}

func TestNbfMoreSiblingsSimple(t *testing.T) {
	db, deleter, err := createTestDatabase("09")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-09",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	rootComment2, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-09",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, rootComment2, "expected comment id to be not zero")

	comments, err := db.GetMoreCommentsNbf(1, 0, 2, []int64{rootComment2}, 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(comments), "expected one sibling comment to be returned")

	assert.Equal(t, 1, comments[0].Id, "expected loaded sibling to have id 1")
}

func TestNbfMoreSiblingsUnordered(t *testing.T) {
	db, deleter, err := createTestDatabase("10")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	comments, err := db.GetMoreCommentsNbf(1, 0, 2, []int64{2, 1}, 1)
	assert.Error(t, err, "expected error when excludeIds parameter is not ordered")
	assert.Equal(t, 0, len(comments))
}

func TestNbfMoreSiblingsExcludeMiddle(t *testing.T) {
	db, deleter, err := createTestDatabase("11")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-11",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	rootComment2, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-11",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, rootComment2, "expected comment id to be not zero")

	rootComment3, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-11",
		ParentId:   0,
	}, 3)
	require.NoError(t, err)
	assert.NotZero(t, rootComment3, "expected comment id to be not zero")

	comments, err := db.GetMoreCommentsNbf(1, 0, 3, []int64{rootComment2}, 99)
	require.NoError(t, err)
	require.Equal(t, 2, len(comments), "expected two sibling comments to be returned")

	assert.Equal(t, 3, comments[0].Id, "expected first loaded sibling to have id 3")
	assert.Equal(t, 1, comments[1].Id, "expected first loaded sibling to have id 1")
}

func TestNbfMoreSiblingsIgnoreYounger(t *testing.T) {
	db, deleter, err := createTestDatabase("12")
	if err != nil {
		t.Fatal(err)
	}
	defer deleter()
	defer db.Close()

	err = db.Setup()
	if err != nil {
		t.Fatal(err)
	}

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-12",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	rootComment2, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-12",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, rootComment2, "expected comment id to be not zero")

	rootComment3, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/test-12",
		ParentId:   0,
	}, 3)
	require.NoError(t, err)
	assert.NotZero(t, rootComment3, "expected comment id to be not zero")

	comments, err := db.GetMoreCommentsNbf(1, 0, 2, []int64{rootComment2}, 1)
	require.NoError(t, err)
	require.Equal(t, 1, len(comments), "expected one sibling comment to be returned")

	assert.Equal(t, 1, comments[0].Id, "expected loaded sibling to have id 1")
}