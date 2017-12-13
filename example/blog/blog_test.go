package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ezbuy/ezorm/db"
)

func init() {
	conf := new(db.MongoConfig)
	conf.DBName = "ezorm"
	conf.MongoDB = "mongodb://127.0.0.1"
	db.Setup(conf)
}

func TestBlogSave(t *testing.T) {
	p := BlogMgr.NewBlog()
	p.Title = "I like ezorm"
	p.Slug = fmt.Sprintf("ezorm_%d", time.Now().Nanosecond())

	_, err := p.Save()
	if err != nil {
		t.Fatal(err)
	}

	id := p.Id()

	b, err := BlogMgr.FindByID(id)
	if err != nil {
		t.Fatal(err.Error())
	}

	fmt.Printf("get blog ok: %#v", b)
}

func TestBlogIter(t *testing.T) {
	testBlogIterForeach(t)
	testBlogIterAppend(t)
	testBlogIterFilter(t)
	testBlogIterShow(t)
}

func testBlogIterForeach(t *testing.T) {
	blogs := []*Blog{
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
	}

	iter := NewBlogIter(blogs)

	iter.Foreach(func(i int, b *Blog) {
		t.Logf("%s", b.Id())
	})
}

func testBlogIterAppend(t *testing.T) {
	blogs := []*Blog{
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
	}

	iter := NewBlogIter(blogs)

	iter.Append(BlogMgr.NewBlog()).
		Append(BlogMgr.NewBlog()).
		Append(BlogMgr.NewBlog()).
		Append(BlogMgr.NewBlog())

	t.Logf("%d", iter.Len())

}

func testBlogIterFilter(t *testing.T) {
	blogs := []*Blog{
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
	}

	iter := NewBlogIter(blogs)

	newIter := iter.Filter(func(b *Blog) bool {
		return b.ID.Counter()%2 == 0
	})

	indexIter := iter.FilterByIndex(func(i int) bool {
		return i%2 == 1
	})

	t.Logf("%d:%d:%d", iter.Len(), newIter.Len(), indexIter.Len())

}

func testBlogIterShow(t *testing.T) {

	blogs := []*Blog{
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
		BlogMgr.NewBlog(),
	}

	iter := NewBlogIter(blogs)
	filIndex := func(i int) bool { return i%2 == 0 }
	fileId := func(b *Blog) bool { return b.ID.Counter()%2 == 0 }
	foreach := func(i int, b *Blog) {}

	length := iter.Append(BlogMgr.NewBlog()).Foreach(foreach).
		Append(BlogMgr.NewBlog()).Foreach(foreach).
		Append(BlogMgr.NewBlog()).Foreach(foreach).
		Append(BlogMgr.NewBlog()).Foreach(foreach).
		Append(BlogMgr.NewBlog()).Foreach(foreach).
		FilterByIndex(filIndex).Foreach(foreach).
		Filter(fileId).Foreach(foreach).
		Len()

	t.Log(length)
}
