package domain

type PageQueryString struct {
	Search   string `url:"q"`
	Page     int    `url:"page"`
	PageSize int    `url:"pageSize"`
}
