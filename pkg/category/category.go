package category

import (
	pb "github.com/andreymgn/RSOI-category/pkg/category/proto"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusNoPostTitle    = status.Error(codes.InvalidArgument, "post title is required")
	statusNotFound       = status.Error(codes.NotFound, "post not found")
	statusInvalidUUID    = status.Error(codes.InvalidArgument, "invalid UUID")
	statusNoCategoryName = status.Error(codes.InvalidArgument, "category name is required")
)

func internalError(err error) error {
	return status.Error(codes.Internal, err.Error())
}

// SingleCategory converts Category to SingleCategory
func (c *Category) SingleCategory() *pb.SingleCategory {
	res := new(pb.SingleCategory)
	res.Uid = c.UID.String()
	res.UserUid = c.UserUID.String()
	res.Name = c.Name

	return res
}

// ListCategories returns categories
func (s *Server) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	categories, err := s.db.getAllCategories(pageSize, req.PageNumber)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.ListCategoriesResponse)
	for _, category := range categories {
		categoryResponse := category.SingleCategory()

		res.Categories = append(res.Categories, categoryResponse)
	}

	res.PageSize = pageSize
	res.PageNumber = req.PageNumber

	return res, nil
}

// GetCategoryAdmin returns admin of category
func (s *Server) GetCategoryInfo(ctx context.Context, req *pb.GetCategoryInfoRequest) (*pb.SingleCategory, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	category, err := s.db.getCategoryInfo(uid)
	switch err {
	case nil:
		return category.SingleCategory(), nil
	case errNotFound:
		return nil, statusNotFound
	default:
		return nil, internalError(err)
	}
}

// CreateCategory creates a new post category
func (s *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.SingleCategory, error) {
	if req.Name == "" {
		return nil, statusNoCategoryName
	}

	userUID, err := uuid.Parse(req.UserUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	category, err := s.db.createCategory(req.Name, userUID)
	if err != nil {
		return nil, internalError(err)
	}

	return category.SingleCategory(), nil
}
