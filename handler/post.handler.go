package handler

import (
	"net/http"
	"sort"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/KuroNeko6666/prima-api/helper"
	"github.com/KuroNeko6666/prima-api/interfaces/form"
	"github.com/KuroNeko6666/prima-api/interfaces/response"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetPost(ctx *fiber.Ctx) error {
	var user model.User
	var resPosts []response.Posts
	var posts []model.Post

	search := ctx.Query("search", "")
	limit := ctx.QueryInt("limit", 10)
	offset := ctx.QueryInt("offset", 0)

	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})

	if r := database.DB.Model(&posts).
		Preload("User").Preload("Images").
		Preload("Comments", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("User").Limit(5)
		}).
		Preload("Comments.SubComments", func(tx *gorm.DB) *gorm.DB {
			return tx.Preload("User").Limit(5)
		}).
		Preload("Likes", func(tx *gorm.DB) *gorm.DB {
			return tx.Limit(5)
		}).
		Where("caption LIKE ?", "%"+search+"%").
		Limit(limit).Offset(offset).
		Not("user_id = ?", user.ID).Find(&posts); r.Error != nil && r.RowsAffected == 0 {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: r.Error.Error(),
		})
	}

	for _, item := range posts {
		var temp response.Posts
		copier.CopyWithOption(&temp, item, copier.Option{IgnoreEmpty: true})
		temp.LikeCount = database.DB.Model(&item).Association("Likes").Count()
		temp.CommentCount = database.DB.Model(&item).Association("Comments").Count()
		db := database.DB.Model(&model.Follower{}).Where("follower_id = ?", item.UserID).Where("following_id", user.ID).Find(&model.Follower{})
		if db.RowsAffected == 0 {
			temp.FollowStatus = "follow"
		} else {
			temp.FollowStatus = "followed"
		}
		resPosts = append(resPosts, temp)
	}

	sort.Slice(resPosts, func(i, j int) bool {
		return resPosts[i].FollowStatus > resPosts[j].FollowStatus
	})

	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_FINE,
		Data:    resPosts,
	})
}

func CreatePost(ctx *fiber.Ctx) error {
	var form form.Post
	var post model.Post
	var user model.User

	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})

	rawFiles, err := ctx.MultipartForm()

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := helper.Validate().Struct(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	files, r := helper.FilesUpload(rawFiles, "post-images", user)

	if r != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(response.Message{
			Message: config.RES_INTERNAL_SERVER,
		})
	}

	for _, file := range files {
		var postImage model.PostImage
		postImage.Name = file.FileName
		postImage.Url = "http://192.168.100.107:9090/post-images/" + file.FileName
		post.Images = append(post.Images, postImage)
	}

	post.UserID = user.ID
	post.Caption = form.Caption

	if db := database.DB.Model(&post).Create(&post); db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})
}

func UpdatePost(ctx *fiber.Ctx) error {
	var form form.Post
	var post model.Post
	var user model.User

	sessionID := ctx.Cookies("session_id")
	postID := ctx.Query("post_id")

	if r := helper.GetUserLoggedIn(&user, sessionID); r != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_UNAUTHORIZED,
		})
	}

	rawFiles, err := ctx.MultipartForm()

	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := helper.Validate().Struct(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	post.ID = postID

	if db := database.DB.Model(&post).Find(&post); db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if user.ID != post.UserID {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_ERR_ACCESS,
		})
	}

	files, r := helper.FilesUpload(rawFiles, "post-images", user)

	if r != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(response.Message{
			Message: config.RES_INTERNAL_SERVER,
		})
	}

	for _, file := range files {
		var postImage model.PostImage
		postImage.PostID = post.ID
		postImage.Name = file.FileName
		postImage.Url = "http://192.168.100.107:9090/post-images/" + file.FileName
		if db := database.DB.Model(&postImage).Create(&postImage); db.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: db.Error.Error(),
			})
		}
	}

	post.Caption = form.Caption

	if db := database.DB.Model(&post).Updates(&post); db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})
}

func DeletePost(ctx *fiber.Ctx) error {
	var post model.Post
	var postImages []model.PostImage
	var user model.User

	postID := ctx.Query("post_id")
	post.ID = postID

	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})

	if db := database.DB.Model(&post).Preload("User").Preload("Images").Find(&post); db.Error != nil || db.RowsAffected == 0 {

		if db.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: db.Error.Error(),
			})
		}

		return ctx.Status(http.StatusNotFound).JSON(response.Message{
			Message: config.RES_NOT_FOUND,
		})
	}

	if user.ID != post.UserID {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_ERR_ACCESS,
		})
	}

	for _, file := range post.Images {
		if r := helper.FileDelete(file.Name, "post-images"); r != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: r.Error(),
			})
		}
	}

	if r := database.DB.Model(&postImages).Where("post_id = ?", post.ID).Delete(&postImages); r.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: r.Error.Error(),
		})
	}

	if db := database.DB.Model(&post).Delete(&post); db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})
}

func PostLike(ctx *fiber.Ctx) error {
	var user model.User
	var postLike model.PostLike
	var form form.Target

	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	postLike.UserID = user.ID
	postLike.PostID = form.TargetID

	if db := database.DB.Model(&postLike).Find(&postLike); db.Error != nil || db.RowsAffected != 0 {
		if db.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: config.RES_CONFLICT,
			})
		}

		if db := database.DB.Model(&postLike).Delete(&postLike); db.Error != nil || db.RowsAffected == 0 {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: config.RES_CONFLICT,
			})
		}
		return ctx.Status(http.StatusOK).JSON(response.Message{
			Message: config.RES_UNLIKE,
		})
	}

	if db := database.DB.Model(&postLike).Create(&postLike); db.Error != nil || db.RowsAffected == 0 {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: config.RES_CONFLICT,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_LIKE,
	})

}

func Comment(ctx *fiber.Ctx) error {
	var form form.PostComment
	var comment model.PostComment
	var user model.User

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	copier.CopyWithOption(&comment, &form, copier.Option{IgnoreEmpty: true})
	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})
	comment.UserID = user.ID

	if db := database.DB.Model(&comment).Omit("User").Omit("Post").Create(&comment); db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})

}

func SubComment(ctx *fiber.Ctx) error {
	var form form.PostComment
	var comment model.PostComment
	var subComment model.PostComment
	var user model.User

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	copier.CopyWithOption(&subComment, &form, copier.Option{IgnoreEmpty: true})
	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})
	subComment.UserID = user.ID
	comment.ID = form.CommentID

	db := database.DB.Model(&comment).Find(&comment)
	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected == 0 {
		return ctx.Status(http.StatusNotFound).JSON(response.Message{
			Message: config.RES_NOT_FOUND,
		})
	}

	if r := database.DB.Model(&comment).Association("SubComment").Append(&subComment); r != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})

}

func DeleteComment(ctx *fiber.Ctx) error {
	var form form.Target
	var comment model.PostComment
	var user model.User

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	copier.CopyWithOption(&user, ctx.Locals("user"), copier.Option{IgnoreEmpty: true})

	comment.ID = form.TargetID

	db := database.DB.Model(&comment).Preload("Post").Find(&comment)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected == 0 {
		return ctx.Status(http.StatusNotFound).JSON(response.Message{
			Message: config.RES_NOT_FOUND,
		})
	}

	if user.ID != comment.UserID {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_ERR_ACCESS,
		})
	}

	if r := database.DB.Model(&comment).Delete(&comment); r.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})

}
