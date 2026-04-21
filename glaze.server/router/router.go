package router

import (
	"net/http"
	"strings"
	"glaze/internal/user"
	"glaze/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRouter(
	userHandler *user.Handler,
) {
	r = gin.Default()
	//r = gin.New()
	//r.Use(gin.Recovery())

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173" ||
				origin == "https://resourcify-three.vercel.app" ||
				origin == "https://www.linksaver.in" ||
				origin == "https://linksaver.in" ||
				strings.HasPrefix(origin, "chrome-extension://mfnbnegonedhppphoceeomjabelbjnnn")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to resourcify")
	})

	userRouter := r.Group("/user")
	{
		userRouter.POST("/create", userHandler.SignUp)
		userRouter.POST("/login", userHandler.LogIn)
		userRouter.GET("/me", middleware.RequireAuth, userHandler.Me)
	}

	//mailTestingRouter := r.Group("/mail-testing")
	//mailTestingRouter.GET("/", func(c *gin.Context) {
	//
	//	go func() {
	//		config.LoadMailerConfig()
	//		fmt.Println(config.MailerConfig)
	//		m := mail.NewMessage()
	//		m.SetHeader("From", fmt.Sprintf("%s <%s>", config.MailerConfig.SenderName, config.MailerConfig.SenderEmail))
	//		m.SetHeader("To", "anuragdaksh77777@gmail.com")
	//		m.SetHeader("Subject", "hi")
	//
	//		m.SetBody("text/plain", "hello world")
	//		//m.AddAlternative("text/html", htmlBody)
	//
	//		d := mail.NewDialer(config.MailerConfig.Host, config.MailerConfig.Port, config.MailerConfig.Username, config.MailerConfig.Password)
	//		//d.SSL = true
	//		if err := d.DialAndSend(m); err != nil {
	//			log.Println(err)
	//		}
	//	}()
	//
	//	c.JSON(http.StatusOK, gin.H{"msg": "\"hello world\""})

	//})
}

func Start(addr string) error {
	return r.Run(addr)
}
