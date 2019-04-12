package category

import (
	pb "github.com/andreymgn/RSOI-category/pkg/category/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	statusCategoryNotFound = status.Error(codes.NotFound, "category not found")
	statusReportNotFound   = status.Error(codes.NotFound, "report not found")
	statusInvalidUUID      = status.Error(codes.InvalidArgument, "invalid UUID")
	statusNoCategoryName   = status.Error(codes.InvalidArgument, "category name is required")
	statusNoReportReason   = status.Error(codes.InvalidArgument, "report reason is required")
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

// SingleReport converts Report to SingleReport
func (r *Report) SingleReport() (*pb.SingleReport, error) {
	createdAtProto, err := ptypes.TimestampProto(r.CreatedAt)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.SingleReport)
	res.Uid = r.UID.String()
	res.CategoryUid = r.CategoryUID.String()
	res.PostUid = r.PostUID.String()
	res.CommentUid = r.CommentUID.String()
	res.Reason = r.Reason
	res.CreatedAt = createdAtProto

	return res, nil
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
		return nil, statusCategoryNotFound
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

// ListReports returns list of reports in some category
func (s *Server) ListReports(ctx context.Context, req *pb.ListReportsRequest) (*pb.ListReportsResponse, error) {
	var pageSize int32
	if req.PageSize == 0 {
		pageSize = 10
	} else {
		pageSize = req.PageSize
	}

	uid, err := uuid.Parse(req.CategoryUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	reports, err := s.db.getAllReports(uid, pageSize, req.PageNumber)
	if err != nil {
		return nil, internalError(err)
	}

	res := new(pb.ListReportsResponse)
	for _, report := range reports {
		reportResponse, err := report.SingleReport()
		if err != nil {
			return nil, err
		}

		res.Reports = append(res.Reports, reportResponse)
	}

	res.PageSize = pageSize
	res.PageNumber = req.PageNumber

	return res, nil
}

// CreateReport creates new report
func (s *Server) CreateReport(ctx context.Context, req *pb.CreateReportRequest) (*pb.SingleReport, error) {
	if req.Reason == "" {
		return nil, statusNoReportReason
	}

	categoryUID, err := uuid.Parse(req.CategoryUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	postUID, err := uuid.Parse(req.PostUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	commentUID, err := uuid.Parse(req.CommentUid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	report, err := s.db.createReport(categoryUID, postUID, commentUID, req.Reason)
	if err != nil {
		return nil, internalError(err)
	}

	return report.SingleReport()
}

// DeleteReport deletes report by ID
func (s *Server) DeleteReport(ctx context.Context, req *pb.DeleteReportRequest) (*pb.DeleteReportResponse, error) {
	uid, err := uuid.Parse(req.Uid)
	if err != nil {
		return nil, statusInvalidUUID
	}

	err = s.db.deleteReport(uid)
	switch err {
	case nil:
		return new(pb.DeleteReportResponse), nil
	case errNotFound:
		return nil, statusReportNotFound
	default:
		return nil, internalError(err)
	}
}
