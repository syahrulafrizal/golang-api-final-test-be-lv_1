package usecase_admin

import (
	"app/domain"
	mongo_model "app/domain/model/mongo"
	request_model "app/domain/model/request"
	helper "app/helpers"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Yureka-Teknologi-Cipta/yureka/helpers"
	"github.com/Yureka-Teknologi-Cipta/yureka/response"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *appUsecase) BlogList(ctx context.Context, claim domain.JWTClaimAdmin, urlQuery url.Values) response.Base {
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

func (u *appUsecase) BlogCreate(ctx context.Context, claim domain.JWTClaimAdmin, payload request_model.BlogRequest, request *http.Request) response.Base {
	ctx, cancel := context.WithTimeout(ctx, u.contextTimeout)
	defer cancel()

	errValidation := make(map[string]string)

	// validating request
	file, uploadedFile, errReq := request.FormFile("file")
	if errReq != nil {
		errValidation["file"] = "file field is required"
	}

	if uploadedFile == nil || file == nil {
		errValidation["file"] = "file field is required"
	} else {
		typeDocument := uploadedFile.Header.Get("Content-Type")

		if !helper.InArrayString(typeDocument, request_model.AllowedMimeTypes) {
			errValidation["file"] = "field file is not valid type"
		}

		fileSize := uploadedFile.Size
		maxFileSize := int64(1 * 1024 * 1024) // 1 MB in bytes

		if fileSize > maxFileSize {
			errValidation["file"] = "file size exceeds the maximum limit of 1 MB"
		}

		defer file.Close()
	}

	if len(request.Form["title"]) == 0 || request.Form["title"][0] == "" {
		errValidation["title"] = "title field is required"
	} else {
		payload.Title = request.Form["title"][0]
	}

	if len(request.Form["content"]) == 0 || request.Form["content"][0] == "" {
		errValidation["content"] = "content field is required"
	} else {
		payload.Content = request.Form["content"][0]
	}

	if len(errValidation) > 0 {
		return response.ErrorValidation(errValidation, "error validation")
	}

	adminID := claim.AdminID

	// check the db
	admin, err := u.mongodbRepo.FetchOneAdmin(ctx, mongo_model.AdminFilter{
		DefaultFilter: mongo_model.DefaultFilter{
			IDStr: &adminID,
		},
	})
	if err != nil {
		return response.Error(500, err.Error())
	}

	if admin == nil {
		return response.Error(400, "admin user not found")
	}

	// generate unique id
	newBlogID := primitive.NewObjectID()

	// Save the file to the media folder
	if _, err := os.Stat("media"); os.IsNotExist(err) {
		err = os.Mkdir("media", os.ModePerm)
		if err != nil {
			errValidation["file"] = "failed to create media directory"
			return response.ErrorValidation(errValidation, "error validation")
		}
	}
	mediaPath := "media/" + newBlogID.Hex() + "_" + uploadedFile.Filename
	out, err := os.Create(mediaPath)
	if err != nil {
		errValidation["file"] = "failed to save file"
	} else {
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			errValidation["file"] = "failed to save file"
		}
	}

	baseURL := request.Host
	if request.TLS != nil {
		baseURL = "https://" + baseURL
	} else {
		baseURL = "http://" + baseURL
	}

	newBlog := mongo_model.Blog{
		ID:        newBlogID,
		Title:     payload.Title,
		Content:   payload.Content,
		Thumbnail: baseURL + "/" + mediaPath,
		CreatedBy: mongo_model.BlogCreatedByNested{
			ID:   admin.ID.Hex(),
			Name: admin.Name,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.mongodbRepo.CreateBlog(ctx, &newBlog)
	if err != nil {
		return response.Error(500, err.Error())
	}

	return response.Success(newBlog)
}
