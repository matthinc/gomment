package sqlite_test

import (
	"os"
	"testing"

	"github.com/matthinc/gomment/model"
	"github.com/matthinc/gomment/persistence"
	"github.com/matthinc/gomment/persistence/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// integration test suite for testing the order of comments retrieved by NSF/OSF

type getCommentsXsfT func(path string, maxDepth int, limit int) ([]model.Comment, persistence.ThreadMetaInfo, error)

func createTestDb(t *testing.T, testName string) (*sqlite.DB, func()) {
	db := sqlite.New()

	dbFilename := "integration-test-" + testName + ".sqlite"
	os.Remove(dbFilename)

	err := db.Open(dbFilename)
	if err != nil {
		t.Fatal(err, "failed to open test database")
		return nil, nil
	}

	err = db.Setup()
	if err != nil {
		t.Fatal(err, "failed to setup database")
		return nil, nil
	}

	return &db, func() {
		defer db.Close()
		os.Remove(dbFilename)
	}
}

func TestXsfNoComments(t *testing.T) {
	db, deleter := createTestDb(t, "xsf-01")
	defer deleter()

	var (
		comments []model.Comment
		meta     persistence.ThreadMetaInfo
		err      error
	)

	testFns := []getCommentsXsfT{db.GetCommentsNsf, db.GetCommentsOsf}

	for _, testFn := range testFns {
		comments, meta, err = testFn("/", 99, 99)
		require.NoError(t, err)

		assert.Equal(t, 0, meta.NumRoot, "expected number of root comments to be 0")
		assert.Equal(t, 0, meta.NumTotal, "expected number of total comments to be 0")

		require.Equal(t, 0, len(comments), "expected no comments to be returned")
	}
}

func TestXsfOneComment(t *testing.T) {
	db, deleter := createTestDb(t, "xsf-02")
	defer deleter()

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	testFns := []getCommentsXsfT{db.GetCommentsNsf, db.GetCommentsOsf}

	for _, testFn := range testFns {
		comments, meta, err := testFn("/", 99, 99)
		require.NoError(t, err)

		assert.Equal(t, 1, meta.NumRoot, "expected number of root comments to be 1")
		assert.Equal(t, 1, meta.NumTotal, "expected number of total comments to be 1")

		require.Equal(t, 1, len(comments), "expected one comment to be returned")

		assert.Equal(t, int64(1), comments[0].CreatedAt, "expected loaded comment to have created_at 1")
	}
}

func TestXsfTwoSiblings(t *testing.T) {
	db, deleter := createTestDb(t, "xsf-03")
	defer deleter()

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	rootComment2, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, rootComment2, "expected comment id to be not zero")

	testFns := []getCommentsXsfT{db.GetCommentsNsf, db.GetCommentsOsf}

	for idx, testFn := range testFns {
		comments, meta, err := testFn("/", 99, 99)
		require.NoError(t, err)

		assert.Equal(t, 2, meta.NumRoot, "expected number of root comments to be 2")
		assert.Equal(t, 2, meta.NumTotal, "expected number of total comments to be 2")

		require.Equal(t, 2, len(comments), "expected two comments to be returned")

		if idx == 0 {
			assert.Equal(t, int64(2), comments[0].CreatedAt, "expected newest comment at position 0")
			assert.Equal(t, int64(1), comments[1].CreatedAt, "expected oldest comment at position 1")
		} else {
			assert.Equal(t, int64(1), comments[0].CreatedAt, "expected oldest comment at position 0")
			assert.Equal(t, int64(2), comments[1].CreatedAt, "expected newest comment at position 1")
		}
	}
}

func TestXsfParentChild(t *testing.T) {
	db, deleter := createTestDb(t, "xsf-04")
	defer deleter()

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	childComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   rootComment1,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, childComment1, "expected comment id to be not zero")

	testFns := []getCommentsXsfT{db.GetCommentsNsf, db.GetCommentsOsf}

	for idx, testFn := range testFns {
		comments, meta, err := testFn("/", 99, 99)
		require.NoError(t, err)

		assert.Equal(t, 1, meta.NumRoot, "expected number of root comments to be 1")
		assert.Equal(t, 2, meta.NumTotal, "expected number of total comments to be 2")

		require.Equal(t, 2, len(comments), "expected two comments to be returned")

		if idx == 0 {
			assert.Equal(t, int64(2), comments[0].CreatedAt, "expected newest comment at position 0")
			assert.Equal(t, int64(1), comments[1].CreatedAt, "expected oldest comment at position 1")
		} else {
			assert.Equal(t, int64(1), comments[0].CreatedAt, "expected oldest comment at position 0")
			assert.Equal(t, int64(2), comments[1].CreatedAt, "expected newest comment at position 1")
		}
	}
}

func TestXsfLimit(t *testing.T) {
	db, deleter := createTestDb(t, "xsf-05")
	defer deleter()

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	rootComment2, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, rootComment2, "expected comment id to be not zero")

	rootComment3, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 3)
	require.NoError(t, err)
	assert.NotZero(t, rootComment3, "expected comment id to be not zero")

	testFns := []getCommentsXsfT{db.GetCommentsNsf, db.GetCommentsOsf}

	for limit := range [4]int{} {
		for fnIdx, testFn := range testFns {
			comments, meta, err := testFn("/", 99, limit)
			require.NoError(t, err)

			assert.Equal(t, 3, meta.NumRoot, "expected number of root comments to be 3")
			assert.Equal(t, 3, meta.NumTotal, "expected number of total comments to be 3")

			numComments := min(3, limit)
			require.Equal(t, numComments, len(comments), "expected number of returned comments to match")

			for i := 0; i < numComments; i++ {
				if fnIdx == 0 {
					assert.Equal(t, int64(3-i), comments[i].CreatedAt, "expected comment to be in right order")
				} else {
					assert.Equal(t, int64(1+i), comments[i].CreatedAt, "expected comment to be in right order")
				}
			}
		}
	}
}

func TestXsfDepth(t *testing.T) {
	db, deleter := createTestDb(t, "xsf-06")
	defer deleter()

	rootComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   0,
	}, 1)
	require.NoError(t, err)
	assert.NotZero(t, rootComment1, "expected comment id to be not zero")

	childComment1, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   rootComment1,
	}, 2)
	require.NoError(t, err)
	assert.NotZero(t, childComment1, "expected comment id to be not zero")

	childComment2, err := db.CreateComment(&model.CommentCreation{
		ThreadPath: "/",
		ParentId:   childComment1,
	}, 3)
	require.NoError(t, err)
	assert.NotZero(t, childComment2, "expected comment id to be not zero")

	testFns := []getCommentsXsfT{db.GetCommentsNsf, db.GetCommentsOsf}

	for depth := range [4]int{} {
		for fnIdx, testFn := range testFns {
			comments, meta, err := testFn("/", depth, 99)
			require.NoError(t, err)

			assert.Equal(t, 1, meta.NumRoot, "expected number of root comments to be 1")
			assert.Equal(t, 3, meta.NumTotal, "expected number of total comments to be 3")

			numComments := min(3, depth+1)
			require.Equal(t, numComments, len(comments), "expected number of returned comments to match")

			for i := 0; i < numComments; i++ {
				if fnIdx == 0 {
					assert.Equal(t, int64(numComments-i), comments[i].CreatedAt, "expected comment to be in right order")
				} else {
					assert.Equal(t, int64(1+i), comments[i].CreatedAt, "expected comment to be in right order")
				}
			}
		}
	}
}
