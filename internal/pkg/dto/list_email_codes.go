package dto

import "github.com/fidesy/sdk/common/postgres"

type ListEmailCodesParams struct {
	Filter     ListEmailCodesFilter
	Pagination postgres.Pagination
}

type ListEmailCodesFilter struct {
	IDIn []string
}
