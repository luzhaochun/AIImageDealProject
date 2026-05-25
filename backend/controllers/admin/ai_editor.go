package admin

import (
	"io"

	"backend/controllers"
	"backend/pkg/ai"
	"backend/pkg/apperr"
	"backend/services"

	beego "github.com/beego/beego/v2/server/web"
)

type AiEditorController struct {
	controllers.BaseController
}

func (c *AiEditorController) Prepare() {
	if c.User() == nil || c.User().Role != "admin" {
		c.Ctx.Output.SetStatus(403)
		c.Data["json"] = map[string]string{"code": "FORBIDDEN", "message": "仅管理员可用"}
		c.ServeJSON()
		c.StopRun()
	}
}

func (c *AiEditorController) Config() {
	mode, _ := beego.AppConfig.String("ai_edit_mode")
	if mode == "" {
		mode = "mock"
	}
	key, _ := beego.AppConfig.String("openai_api_key")
	model, _ := beego.AppConfig.String("openai_image_model")
	if model == "" {
		model = "gpt-image-2"
	}
	studioReady := mode == "openai" && key != ""
	c.Success(map[string]interface{}{
		"mode":         mode,
		"has_key":      key != "",
		"model":        model,
		"studio_ready": studioReady,
		"studio_modes": ai.AllStudioModes(),
	})
}

func (c *AiEditorController) StudioInpaint() {
	c.studioEdit()
}

func (c *AiEditorController) StudioEdit() {
	c.studioEdit()
}

func (c *AiEditorController) studioEdit() {
	imageFile, _, err := c.GetFile("image")
	if err != nil {
		c.Fail(apperr.New("VALIDATION_ERROR", "请上传图片", apperr.ErrValidation))
		return
	}
	defer imageFile.Close()

	imgData, err := io.ReadAll(imageFile)
	if err != nil {
		c.Fail(err)
		return
	}

	var maskData []byte
	if maskFile, _, err := c.GetFile("mask"); err == nil {
		defer maskFile.Close()
		maskData, _ = io.ReadAll(maskFile)
	}

	prompt := c.GetString("prompt")
	mode := c.GetString("mode")
	res, err := (&services.AiEditService{}).StudioEdit(imgData, maskData, prompt, mode)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(res)
}

func (c *AiEditorController) Inpaint() {
	imageFile, _, err := c.GetFile("image")
	if err != nil {
		c.Fail(apperr.New("VALIDATION_ERROR", "请上传图片", apperr.ErrValidation))
		return
	}
	defer imageFile.Close()
	maskFile, _, err := c.GetFile("mask")
	if err != nil {
		c.Fail(apperr.New("VALIDATION_ERROR", "请上传蒙版", apperr.ErrValidation))
		return
	}
	defer maskFile.Close()

	imgData, err := io.ReadAll(imageFile)
	if err != nil {
		c.Fail(err)
		return
	}
	maskData, err := io.ReadAll(maskFile)
	if err != nil {
		c.Fail(err)
		return
	}

	prompt := c.GetString("prompt")
	res, err := (&services.AiEditService{}).Inpaint(imgData, maskData, prompt)
	if err != nil {
		c.Fail(err)
		return
	}
	c.Success(res)
}
