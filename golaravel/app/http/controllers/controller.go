package controllers

import (
	"golaravel/app/container"
	"golaravel/app/http/request"
	"golaravel/app/http/response"
	"golaravel/app/validation"
)

type Controller struct {
	Container *container.Container
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) SetContainer(container *container.Container) {
	c.Container = container
}

func (c *Controller) Validate(req *request.Request, rules map[string][]string) (*validation.Validator, bool) {
	data := req.All()
	v := validation.Make(data, rules)

	if v.Fails() {
		return v, false
	}

	return v, true
}

func (c *Controller) ValidateAndRespond(req *request.Request, res *response.Response, rules map[string][]string) bool {
	v, passed := c.Validate(req, rules)
	if !passed {
		res.ValidationError(v.Errors())
		return false
	}
	return true
}

type ResourceController interface {
	Index(*request.Request, *response.Response) error
	Create(*request.Request, *response.Response) error
	Store(*request.Request, *response.Response) error
	Show(*request.Request, *response.Response) error
	Edit(*request.Request, *response.Response) error
	Update(*request.Request, *response.Response) error
	Destroy(*request.Request, *response.Response) error
}

type BaseResourceController struct {
	*Controller
}

func (c *BaseResourceController) Index(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}

func (c *BaseResourceController) Create(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}

func (c *BaseResourceController) Store(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}

func (c *BaseResourceController) Show(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}

func (c *BaseResourceController) Edit(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}

func (c *BaseResourceController) Update(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}

func (c *BaseResourceController) Destroy(req *request.Request, res *response.Response) error {
	return res.Status(501).JSON(map[string]string{"error": "Not Implemented"})
}
