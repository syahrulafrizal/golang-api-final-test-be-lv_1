package usecase_public

import (
	mongo_model "app/domain/model/mongo"
	"context"
	"net/url"

	"github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *publicUsecase) BlogList(ctx context.Context, urlQuery url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := helpers.GetLimitOffset(urlQuery)

	options := mongo_model.BlogFilter{
		DefaultFilter: mongo_model.DefaultFilter{
			Limit:  &limit,
			Offset: &offset,
		},
	}

	// count first
	totalDocuments := u.mongodbRepo.CountBlog(ctx, options)
	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// sorting here
	options.Sorts = helpers.GetSorts(urlQuery, mongo_model.BlogAllowedSort)

	// check the db
	cur, err := u.mongodbRepo.FetchBlog(ctx, options)
	if err != nil {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}
	defer cur.Close(ctx)

	list := make([]interface{}, 0)
	for cur.Next(ctx) {
		row := mongo_model.Blog{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Blog Decode ", err)
			return response.Success(response.List{
				List:  []interface{}{},
				Page:  page,
				Limit: limit,
				Total: totalDocuments,
			})
		}

		list = append(list, row)
	}

	return response.Success(response.List{
		List:  list,
		Page:  page,
		Limit: limit,
		Total: totalDocuments,
	})
}
