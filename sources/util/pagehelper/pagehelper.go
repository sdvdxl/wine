package pagehelper

type page struct {
	CurPage      int
	TotalPages   int
	TotalRecords int
	PageSize     int
}

func Paging(curPage, pageSize, totalRecords int) page {

	totalPages := totalRecords / pageSize

	if totalRecords%pageSize != 0 {
		totalPages += 1
	}

	return page{CurPage: curPage,
		PageSize:     pageSize,
		TotalRecords: totalRecords,
		TotalPages:   totalPages,
	}
}
