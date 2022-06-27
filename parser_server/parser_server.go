package parser_server

import (
	"context"
	"parser/config"
	gen "parser/generated"
	"parser/platform"
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
	searchRes := ps.p.SearchByCondition(&platform.Condition{
		Keyword:  in.Keyword,
		Type:     in.Type,
		Genres:   in.Genres,
		YearFrom: in.StartYear,
		YearTo:   in.EndYear,
		Coutries: in.Countries,
	}, ps.cnfg.Proxy) //TODO RETURN ERRORS

	var pr gen.ParsedResults
	for i := range searchRes {
		pr.Data = append(pr.Data, &gen.ParsedData{
			Name: searchRes[i].Name,
			Ref:  searchRes[i].Ref,
			Img:  searchRes[i].Img,
		})
	}
	return &pr, nil //TODO ERROR
}
