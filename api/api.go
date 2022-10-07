package api

import (
	"reflect"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/althafariq/discusspedia/repository"
	"github.com/gin-gonic/gin"
)

type API struct {
	commentRepo       repository.CommentRepository
	likeRepo          repository.LikeRepository
	notifRepo         repository.NotificationRepository
	postRepo          repository.PostRepository
	userRepo          repository.UserRepository
	categoryRepo      repository.CategoryRepository
	questionnaireRepo repository.QuestionnaireRepository
	router            *gin.Engine
}

func NewAPI(
	commentRepo repository.CommentRepository,
	likeRepo repository.LikeRepository,
	notifRepo repository.NotificationRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
	categoryRepo repository.CategoryRepository,
	questionnaireRepo repository.QuestionnaireRepository,
) API {
		router := gin.Default()
		api := API{
			router:            router,
			commentRepo:       commentRepo,
			likeRepo:          likeRepo,
			notifRepo:         notifRepo,
			postRepo:          postRepo,
			userRepo:          userRepo,
			categoryRepo:      categoryRepo,
			questionnaireRepo: questionnaireRepo,
		}
		// config := cors.DefaultConfig()
		// // config.AllowOrigins = []string{"*"}
		// config.AllowAllOrigins = true
		// config.AllowCredentials = true
		// config.AddAllowHeaders("Authorization")
		router.Use(cors.New(cors.Config{
			AllowOrigins:   []string{"*"},
         AllowMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
         AllowHeaders:   []string{"Authorization"},
         ExposeHeaders:  []string{"Content-Length"},
			AllowCredentials: true,
		}))

	// Untuk validasi request dengan mengembalikan nama dari tag json jika ada
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	router.Static("/media", "./media")

	router.POST("/api/login", api.login)
	router.POST("/api/register", api.register)
	router.GET("/api/category", api.GetAllCategories)

	profileRouter := router.Group("/api/profile", AuthMiddleware())
	{
		profileRouter.GET("/get", api.getProfile)
		profileRouter.PATCH("/patch", api.updateProfile)
		profileRouter.PUT("/avatar", api.changeAvatar)
	}

	router.GET("/api/post", api.readPosts)
	router.GET("/api/post/:id", api.readPost)
	postRouter := router.Group("/api/post", AuthMiddleware())
	{
		postRouter.POST("/", api.createPost)
		postRouter.PUT("/", api.updatePost)
		postRouter.POST("/images/:id", api.uploadPostImages)
		postRouter.DELETE("/:id", api.deletePost)
	}

	router.GET("/api/comments", api.ReadAllComment)
	commentRoutersWithAuth := router.Group("/api/comments", AuthMiddleware())
	{
		commentRoutersWithAuth.POST("/", api.CreateComment)
		commentRoutersWithAuth.PUT("/", api.UpdateComment)
		commentRoutersWithAuth.DELETE("/:id", api.DeleteComment)
	}

	postLikeRouters := router.Group("/api/post/:id/likes", AuthMiddleware())
	{
		postLikeRouters.POST("/", api.CreatePostLike)
		postLikeRouters.DELETE("/", api.DeletePostLike)
	}

	commentLikeRouters := router.Group("/api/comments/:id/likes", AuthMiddleware())
	{
		commentLikeRouters.POST("/", api.CreateCommentLike)
		commentLikeRouters.DELETE("/", api.DeleteCommentLike)
	}

	notifRouter := router.Group("/api/notifications", AuthMiddleware())
	{
		notifRouter.GET("/", api.GetAllNotifications)
		notifRouter.PUT("/read", api.SetReadNotif)
	}

	router.GET("/api/questionnaires", api.ReadAllQuestionnaires)
	router.GET("/api/questionnaires/:id", api.ReadAllQuestionnaireByID)
	questionnaireRoutersWithAuth := router.Group("/api/questionnaires", AuthMiddleware())
	{
		questionnaireRoutersWithAuth.POST("/", api.CreateQuestionnaire)
		questionnaireRoutersWithAuth.PUT("/", api.UpdateQuestionnaire)
		questionnaireRoutersWithAuth.DELETE("/:id", api.DeleteQuestionnaire)
	}

	return api
}

func (api *API) Handler() *gin.Engine {
	return api.router
}

func (api *API) Start() {
	api.Handler().Run()
}
