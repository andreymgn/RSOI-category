package category

import (
	"errors"
	"testing"

	pb "github.com/andreymgn/RSOI-category/pkg/category/proto"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var (
	errDummy     = errors.New("dummy")
	dummyUID     = uuid.New()
	nilUIDString = uuid.Nil.String()
)

type mockdb struct{}

func (mdb *mockdb) getAllCategories(pageSize, pageNumber int32) ([]*Category, error) {
	result := make([]*Category, 0)
	uid1 := uuid.New()
	uid2 := uuid.New()
	uid3 := uuid.New()

	result = append(result, &Category{uid1, uid2, "aaa"})
	result = append(result, &Category{uid2, uid3, "bbb"})
	result = append(result, &Category{uid3, uid1, "ccc"})
	return result, nil
}

func (mdb *mockdb) getCategoryInfo(uid uuid.UUID) (*Category, error) {
	return &Category{uuid.Nil, uuid.Nil, "aaa"}, nil
}

func (mdb *mockdb) createCategory(name string, userUID uuid.UUID) (*Category, error) {
	if name == "success" {
		uid := uuid.New()

		return &Category{uid, userUID, name}, nil
	}

	return nil, errDummy
}

func TestListCategories(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.ListCategoriesRequest{}
	_, err := s.ListCategories(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateCategory(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreateCategoryRequest{Name: "success", UserUid: nilUIDString}
	_, err := s.CreateCategory(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateCategoryFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreateCategoryRequest{Name: ""}
	_, err := s.CreateCategory(context.Background(), req)
	if err != statusNoCategoryName {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.CreateCategoryRequest{Name: "fail"}
	_, err = s.CreateCategory(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestGetCategoryAdminByPost(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetCategoryInfoRequest{Uid: nilUIDString}
	_, err := s.GetCategoryInfo(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestGetCategoryAdminByPostFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.GetCategoryInfoRequest{Uid: ""}
	_, err := s.GetCategoryInfo(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}
