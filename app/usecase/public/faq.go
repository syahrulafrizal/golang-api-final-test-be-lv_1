package usecase_public

import (
	mongo_model "app/domain/model/mongo"
	"context"
	"net/url"

	"github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *publicUsecase) FaqList(ctx context.Context, urlQuery url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := helpers.GetLimitOffset(urlQuery)

	options := mongo_model.FaqFilter{
		DefaultFilter: mongo_model.DefaultFilter{
			Limit:  &limit,
			Offset: &offset,
		},
	}

	// count first
	totalDocuments := u.mongodbRepo.CountFaq(ctx, options)
	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// sorting here
	options.Sorts = helpers.GetSorts(urlQuery, mongo_model.FaqAllowedSort)

	// check the db
	cur, err := u.mongodbRepo.FetchFaq(ctx, options)
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
		row := mongo_model.Faq{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("Faq Decode ", err)
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