package configuration

import (
	"go.dbio.me/cmd/types"
)

func getSort() types.Sort {
	return types.Sort{
		{
			Label: "Newest",
			Value: "newest",
			Icon: types.Icon{
				Label: "fa-sort-amount-down",
				Value: "fa",
			},
		},
		{
			Label: "Oldest",
			Value: "oldest",
			Icon: types.Icon{
				Label: "fa-sort-amount-up",
				Value: "fa",
			},
		},
		{
			Default: true,
			Label:   "Popular",
			Value:   "popular",
			Icon: types.Icon{
				Label: "fa-fire",
				Value: "fa",
			},
		},
		{
			Label: "Trending",
			Value: "trending",
			Icon: types.Icon{
				Label: "fa-chart-line",
				Value: "fa",
			},
		},
		{
			Label: "Ascending",
			Value: "ascending",
			Icon: types.Icon{
				Label: "fa-sort-alpha-down",
				Value: "fa",
			},
		},
		{
			Label: "Descending",
			Value: "descending",
			Icon: types.Icon{
				Label: "fa-sort-alpha-up",
				Value: "fa",
			},
		},
	}
}

func GetSort() types.Sort {
	return getSort()
}
