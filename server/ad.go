package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-content/pkg/handler/adDb"
	"grpc-content/pkg/util"
	pb "grpc-content/proto"
)

type adService struct {
	*pb.UnimplementedAdServer
}

func NewAdService() *adService {
	return &adService{}
}

func checkAppidAndGroupid(appid string, groupid uint32) error {
	if appid == "" {
		return status.Error(codes.InvalidArgument, "缺少appid参数")
	}
	if groupid == 0 {
		return status.Error(codes.InvalidArgument, "缺少groupid参数")
	}
	return nil
}

func (a adService) MigrateAd(ctx context.Context, req *pb.MigrateAdRequest) (*pb.MigrateAdResponse, error) {
	adExample := adDb.Ad{}
	err := adExample.Migrate()
	if err != nil {
		return nil, err
	}
	return &pb.MigrateAdResponse{
		Message: "success",
	}, nil
}

func (a adService) AddAd(ctx context.Context, req *pb.AddAdRequest) (*pb.AddAdResponse, error) {
	if err := checkAppidAndGroupid(req.Appid, req.Groupid); err != nil {
		return nil, err
	}

	adExample := adDb.Ad{
		Appid:       req.Appid,
		Groupid:     req.Groupid,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		Remark:      req.Remark,
		Cover:       req.Cover,
		Url:         req.Url,
		Priority:    req.Priority,
		Display:     req.Display,
		State:       req.State,
	}

	adid, err := adExample.Insert()
	if err != nil {
		return nil, err
	}

	return &pb.AddAdResponse{
		Adid: adid,
	}, nil
}

func (a adService) GetAdList(ctx context.Context, req *pb.GetAdListRequest) (*pb.GetAdListResponse, error) {
	if err := checkAppidAndGroupid(req.Appid, req.Groupid); err != nil {
		return nil, err
	}

	util.FormatPageAndPerPage(&req.Page, &req.PerPage)

	adExample := adDb.Ad{
		Appid:   req.Appid,
		Groupid: req.Groupid,
		Page:    req.Page,
		PerPage: req.PerPage,
		Name: req.Name,
		Title: req.Title,
		Display: req.Display,
		State: req.State,
	}
	record, total, err := adExample.GetAdList()
	if err != nil {
		return nil, err
	}

	adList := make([]*pb.AdListElement, len(record))
	for i := 0; i < len(record); i++ {
		ele := record[i]
		adList[i] = &pb.AdListElement{
			Adid:        uint32(ele.ID),
			Name:        ele.Name,
			Title:       ele.Title,
			Description: ele.Description,
			Remark:      ele.Remark,
			Cover:       ele.Cover,
			Url:         ele.Url,
			Priority:    ele.Priority,
			Display:     ele.Display,
			CreateAt:    util.FormatDateTime(ele.CreatedAt),
		}
	}

	return &pb.GetAdListResponse{
		List:  adList,
		Total: total,
	}, nil
}

func (a adService) EditAd(ctx context.Context, req *pb.EditAdRequest) (*pb.EditAdResponse, error) {
	if err := checkAppidAndGroupid(req.Appid, req.Groupid); err != nil {
		return nil, err
	}
	if req.Adid == 0 {
		return nil, status.Error(codes.InvalidArgument, "缺少adid参数")
	}
	adExample := adDb.Ad{
		Appid:       req.Appid,
		Groupid:     req.Groupid,
		Adid:        req.Adid,
		Name:        req.Name,
		Title:       req.Title,
		Description: req.Description,
		Remark:      req.Remark,
		Cover:       req.Cover,
		Url:         req.Url,
		Priority:    req.Priority,
		Display:     req.Display,
		State:       req.State,
	}
	exists, err := adExample.ExistByID()
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, status.Error(codes.NotFound, "记录不存在")
	}

	var rowsAffected uint32
	rowsAffected, err = adExample.Edit()
	if err != nil {
		return nil, err
	}

	return &pb.EditAdResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (a adService) DeleteAd(ctx context.Context, req *pb.DeleteAdRequest) (*pb.DeleteAdResponse, error) {
	if err := checkAppidAndGroupid(req.Appid, req.Groupid); err != nil {
		return nil, err
	}
	adExample := adDb.Ad{
		Appid: req.Appid,
		Groupid: req.Groupid,
		AdidList: req.AdidList,
	}
	rowsAffected, err := adExample.Delete()
	if err != nil {
		return nil, err
	}

	return &pb.DeleteAdResponse{
		RowsAffected: rowsAffected,
	}, nil
}

func (a adService) GetAdPlacementList(ctx context.Context, req *pb.GetAdPlacementListRequest) (*pb.GetAdPlacementListResponse, error) {
	if err := checkAppidAndGroupid(req.Appid, req.Groupid); err != nil {
		return nil, err
	}

	util.FormatPageAndPerPage(&req.Page, &req.PerPage)

	adExample := adDb.Ad{
		Appid:   req.Appid,
		Groupid: req.Groupid,
		Page:    req.Page,
		PerPage: req.PerPage,
	}
	record, total, err := adExample.GetAdPlacementList()
	if err != nil {
		return nil, err
	}

	adPlacementList := make([]*pb.AdPlacementListElement, len(record))
	for i := 0; i < len(record); i++ {
		ele := record[i]
		adPlacementList[i] = &pb.AdPlacementListElement{
			Adid:        uint32(ele.ID),
			Name:        ele.Name,
			Title:       ele.Title,
			Description: ele.Description,
			Remark:      ele.Remark,
			Cover:       ele.Cover,
			Url:         ele.Url,
			Priority:    ele.Priority,
			Display:     ele.Display,
		}
	}

	return &pb.GetAdPlacementListResponse{
		List:  adPlacementList,
		Total: total,
	}, nil
}

