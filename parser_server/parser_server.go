package parser_server

import (
	"context"
	"parser/config"
	gen "parser/generated"
	"parser/platform"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ParserServer struct {
	p    *platform.Platform
	cnfg config.Config
	gen.UnimplementedParserServiceServer
}

func New(cnfg config.Config) *ParserServer {
	return &ParserServer{
		p:    platform.New(cnfg),
		cnfg: cnfg,
	}
}

func (ps *ParserServer) SearchData(
	ctx context.Context, in *gen.Conditions) (*gen.ParsedResults, error) {
	searchRes, err := ps.p.SearchByCondition(&platform.Condition{
		Keyword:  in.Keyword,
		Type:     in.Type,
		Genres:   in.Genres,
		YearFrom: in.StartYear,
		YearTo:   in.EndYear,
		Coutries: in.Countries,
	}, ps.cnfg.Proxy)

	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"got error %v from Parser service, %v", err)
	}

	var pr gen.ParsedResults
	for i := range searchRes {
		pr.Data = append(pr.Data, &gen.ParsedData{
			Name: searchRes[i].Name,
			Ref:  searchRes[i].Ref,
			Img:  searchRes[i].Img,
		})
	}
	return &pr, nil
}
