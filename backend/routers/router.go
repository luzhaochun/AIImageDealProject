package routers

import (
	"backend/controllers"
	"backend/controllers/admin"
	"backend/controllers/auth"
	"backend/controllers/review"
	"backend/controllers/workspace"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/health", &controllers.HealthController{}, "get:Get"),

		beego.NSNamespace("/auth",
			beego.NSRouter("/login", &auth.Controller{}, "post:Login"),
			beego.NSRouter("/me", &auth.Controller{}, "get:Me"),
		),

		beego.NSNamespace("/categories",
			beego.NSRouter("/", &workspace.CategoriesController{}, "get:List"),
		),

		beego.NSNamespace("/imports",
			beego.NSRouter("/", &admin.ImportController{}, "post:Post;get:Get"),
		),

		beego.NSNamespace("/images",
			beego.NSRouter("/claim", &workspace.ImageController{}, "post:Claim"),
			beego.NSRouter("/", &workspace.ImageController{}, "get:List"),
			beego.NSRouter("/:id", &workspace.ImageController{}, "get:Get"),
			beego.NSRouter("/:id/discard", &workspace.ImageController{}, "post:Discard"),
			beego.NSRouter("/:id/start-editing", &workspace.ImageController{}, "post:StartEditing"),
			beego.NSRouter("/:id/versions", &workspace.EditorController{}, "get:Versions;put:SaveVersion"),
			beego.NSRouter("/:id/submit", &workspace.EditorController{}, "post:Submit"),
			beego.NSRouter("/:id/reviews", &review.Controller{}, "get:ListImageReviews;post:Submit"),
		),

		beego.NSNamespace("/reviews",
			beego.NSRouter("/queue", &review.Controller{}, "get:Queue"),
			beego.NSRouter("/stats", &review.Controller{}, "get:Stats"),
			beego.NSRouter("/second-pool", &review.Controller{}, "get:SecondPool"),
			beego.NSRouter("/claim-second", &review.Controller{}, "post:ClaimSecond"),
		),

		beego.NSNamespace("/admin/ai",
			beego.NSRouter("/config", &admin.AiEditorController{}, "get:Config"),
			beego.NSRouter("/inpaint", &admin.AiEditorController{}, "post:Inpaint"),
			beego.NSRouter("/studio-inpaint", &admin.AiEditorController{}, "post:StudioInpaint"),
		),

		beego.NSNamespace("/archives",
			beego.NSRouter("/", &review.ArchiveController{}, "get:List"),
			beego.NSRouter("/export", &review.ArchiveController{}, "post:Export"),
			beego.NSRouter("/exports", &review.ArchiveController{}, "get:ExportList"),
			beego.NSRouter("/:id/complete", &review.ArchiveController{}, "post:Complete"),
		),
	)
	beego.AddNamespace(ns)
}
