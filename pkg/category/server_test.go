package category

import (
	"errors"
	"testing"
	"time"

	pb "github.com/andreymgn/RSOI-category/pkg/category/proto"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

var (
	errDummy     = errors.New("dummy")
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

func (mdb *mockdb) getAllReports(categoryUID uuid.UUID, pageSize, pageNumber int32) ([]*Report, error) {
	result := make([]*Report, 0)
	uid1 := uuid.New()
	uid2 := uuid.New()
	uid3 := uuid.New()

	result = append(result, &Report{uid1, uid2, uid3, uuid.Nil, "aaa", time.Now()})
	result = append(result, &Report{uid2, uid3, uid1, uid2, "bbb", time.Now()})
	result = append(result, &Report{uid3, uid2, uid2, uuid.Nil, "ccc", time.Now()})
	return result, nil
}

func (mdb *mockdb) createReport(categoryUID, postUID, commentUID uuid.UUID, reason string) (*Report, error) {
	if reason == "success" {
		uid := uuid.New()

		return &Report{uid, categoryUID, postUID, commentUID, reason, time.Now()}, nil
	}

	return nil, errDummy
}

func (mdb *mockdb) deleteReport(uid uuid.UUID) error {
	if uid == uuid.Nil {
		return nil
	}

	return errDummy
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

func TestListReports(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.ListReportsRequest{CategoryUid: nilUIDString}
	_, err := s.ListReports(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateReport(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.CreateReportRequest{Reason: "success", CategoryUid: nilUIDString, PostUid: nilUIDString, CommentUid: nilUIDString}
	_, err := s.CreateReport(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestCreateReportFail(t *testing.T) {
	s := &Server{&mockdb{}}

	req := &pb.CreateReportRequest{Reason: "", CategoryUid: nilUIDString, PostUid: nilUIDString, CommentUid: nilUIDString}
	_, err := s.CreateReport(context.Background(), req)
	if err != statusNoReportReason {
		t.Errorf("unexpected error %v", err)
	}

	req = &pb.CreateReportRequest{Reason: "fail"}
	_, err = s.CreateReport(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}

func TestDeleteReport(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeleteReportRequest{Uid: nilUIDString}
	_, err := s.DeleteReport(context.Background(), req)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestDeleteReportFail(t *testing.T) {
	s := &Server{&mockdb{}}
	req := &pb.DeleteReportRequest{Uid: ""}
	_, err := s.DeleteReport(context.Background(), req)
	if err == nil {
		t.Errorf("expected error, got nothing")
	}
}
