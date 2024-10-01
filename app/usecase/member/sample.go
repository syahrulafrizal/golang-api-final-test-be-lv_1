package usecase_member

import (
	"app/domain"
	mongo_model "app/domain/model/mongo"
	app_helper "app/helpers"
	"context"
	"net/url"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
)

func (u *appUsecase) SampleUserList(ctx context.Context, claim domain.JWTClaimUser, urlQuery url.Values) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	page, limit, offset := helpers.GetLimitOffset(urlQuery)

	options := mongo_model.UserFilter{
		DefaultFilter: mongo_model.DefaultFilter{
			Limit:  &limit,
			Offset: &offset,
		},
	}

	if urlQuery.Get("username") != "" {
		options.Username = app_helper.StringPointer(urlQuery.Get("username"))
	}

	// ---------- createdAt filter ----------
	var start, end time.Time
	if t, e := time.Parse("2006-01-02 15:04:05", urlQuery.Get("created_at_start")); e == nil {
		start = t
	}

	if t, e := time.Parse("2006-01-02 15:04:05", urlQuery.Get("created_at_end")); e == nil {
		end = t
	}

	if !start.IsZero() && !end.IsZero() {
		options.CreatedAtRange = &mongo_model.DatetimeRange{
			Start: start,
			End:   end,
		}
	} else if !start.IsZero() {
		options.CreatedAtGte = &start
	} else if !end.IsZero() {
		options.CreatedAtLte = &end
	}
	// ---------- end createdAt filter ----------

	// count first
	totalDocuments := u.mongodbRepo.CountUser(ctx, options)
	if totalDocuments == 0 {
		return response.Success(response.List{
			List:  []interface{}{},
			Page:  page,
			Limit: limit,
			Total: totalDocuments,
		})
	}

	// sorting here
	options.Sorts = helpers.GetSorts(urlQuery, mongo_model.UserAllowedSort)

	// check the db
	cur, err := u.mongodbRepo.FetchUser(ctx, options)
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
		row := mongo_model.User{}
		err := cur.Decode(&row)
		if err != nil {
			logrus.Error("User Decode ", err)
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

func (u *appUsecase) SampleUserDetail(ctx context.Context, claim domain.JWTClaimUser, id string) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	// check the db
	user, err := u.mongodbRepo.FetchOneUser(ctx, mongo_model.UserFilter{
		DefaultFilter: mongo_model.DefaultFilter{
			IDStr: app_helper.StringPointer(id),
		},
	})
	if err != nil || user == nil {
		return response.Error(404, "User not found")
	}

	return response.Success(user)
}
