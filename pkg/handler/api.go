package handler

import (
	"esaccount"
	"esaccount/pkg/util"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

// confidential          godoc
// @Summary      Get all confidential settings
// @Tags         confidential
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200
// @Router       /api/confidential/ [get]
func (h *Handler) getConfidential(c *gin.Context) {
	id, _ := c.Get(userCtx)
	account, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, account.Confidentials)
}

// confidential          godoc
// @Summary      Set confidential settings
// @Tags         confidential
// @Produce      json
// @Param        Authorization    header    string                           true   	"JWT Bearer token (authorization)"
// @Param        confidencial     body      string    true   	"Confidencial settings"
// @Success      200   {object}  esaccount.AuthorizedUser
// @Router       /api/confidential/ [post]
func (h *Handler) setConfidential(c *gin.Context) {
	id, _ := c.Get(userCtx)
	account, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	var input esaccount.ConfidentialSettings
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for k, v := range input {
		account.Confidentials[k] = v
	}
	account.Confidentials["picture_url"] = esaccount.All

	h.repo.SaveAccount(account)

	c.JSON(http.StatusOK, account)
}

// public        godoc
// @Summary      Get all public data of users
// @Tags         confidential
// @Produce      json
// @Param        Authorization    header    string    false   	"JWT Bearer token (authorization)"
// @Param        id  query      int  true  "user id"
// @Success      200
// @Router       /api/public/{id} [get]
func (h *Handler) getPublic(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := h.repo.GetAccount(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	header := c.GetHeader(authorizationHeader)
	allowedFields := []string{}
	isFriend := false
	if header != "" {
		authorized, err := h.services.Authorization.Authorize(header)
		if err == nil {
			isFriend = (slices.Contains(account.Followers, authorized.Id) && slices.Contains(account.Follows, authorized.Id))
		}
	}

	for k, v := range account.Confidentials {
		if (v == esaccount.All) || (isFriend && (v == esaccount.Friends)) {
			allowedFields = append(allowedFields, k)
		}
	}

	res := map[string]any{}
	user, _ := h.services.AuthorizeById(id)

	// auth user
	if slices.Contains(allowedFields, "username") {
		res["username"] = user.Username
	}
	if slices.Contains(allowedFields, "first_name") {
		res["first_name"] = user.FirstName
	}
	if slices.Contains(allowedFields, "last_name") {
		res["last_name"] = user.LastName
	}
	res["picture_url"] = user.LastName
	if slices.Contains(allowedFields, "email") {
		res["email"] = user.LastName
	}

	// account
	if slices.Contains(allowedFields, "description") {
		res["description"] = account.Description
	}
	if slices.Contains(allowedFields, "date_of_birth") {
		res["date_of_birth"] = account.DateOfBirth
	}
	if slices.Contains(allowedFields, "follows") {
		res["follows"] = account.Follows
	}
	if slices.Contains(allowedFields, "followers") {
		res["followers"] = account.Follows
	}

	c.JSON(http.StatusOK, res)
}

type accountResponce struct {
	Description string  `json:"description"`
	DateOfBirth string  `json:"date_of_birth"`
	Follows     []int64 `json:"follows"`
	Followers   []int64 `json:"followers"`
}

// account       godoc
// @Summary      Get account information
// @Tags         account
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200   {object}   handler.accountResponce
// @Router       /api/account [get]
func (h *Handler) getAccount(c *gin.Context) {
	id, _ := c.Get(userCtx)
	account, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, accountResponce{
		Description: account.Description,
		DateOfBirth: account.DateOfBirth,
		Follows:     account.Follows,
		Followers:   account.Followers,
	})
}

// picture       godoc
// @Summary      Get account information
// @Tags         account
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Param		 file			  formData		file	  true 		"Picture"
// @Success      200
// @Router       /api/account/picture [post]
func (h *Handler) picture(c *gin.Context) {
	id, _ := c.Get(userCtx)
	_, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	os.Remove(fmt.Sprintf("./public/usercontent/picture/%d.png", id))
	picFile, err := os.Create(fmt.Sprintf("./public/usercontent/picture/%d.png", id))
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	if _, err := io.Copy(picFile, file); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	header, _ := c.Get(userHeader)
	picUrl := fmt.Sprintf("https://account.easystartup.su/static/picture/%d.png", id)
	_, err = h.services.AuthorizeAndUpdatePicture(header.(string), picUrl)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, fmt.Sprintf("{ \"url\": %s \"\" }", picUrl))
}

type accountInput struct {
	Description string `json:"description"`
	DateOfBirth string `json:"date_of_birth"`
}

// account       godoc
// @Summary      Set account information
// @Tags         account
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Param        account    body    handler.accountInput    true   	"account info"
// @Success      200   {object}   handler.accountInput
// @Router       /api/account [post]
func (h *Handler) editAccount(c *gin.Context) {
	id, _ := c.Get(userCtx)
	account, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	var input accountInput
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account.Description = input.Description
	account.DateOfBirth = input.DateOfBirth
	h.repo.SaveAccount(account)

	c.JSON(http.StatusOK, input)
}

// account       godoc
// @Summary      Get account information by id
// @Tags         account
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Param        id         path    int                     true   	"account id"
// @Success      200   {object}   handler.accountResponce
// @Router       /api/account/{id} [get]
func (h *Handler) getAccountById(c *gin.Context) {
	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || (currentUser.Access != "admin") {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	

	account, err := h.repo.GetAccount(id)
	if err != nil {
		authorized, err := h.services.AuthorizeById(id)
		if (err == nil) && (authorized.Id == id) {
			account, err = h.repo.CreateAccount(id)
		}

		if err != nil {
			util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, accountResponce{
		Description: account.Description,
		DateOfBirth: account.DateOfBirth,
		Follows:     account.Follows,
		Followers:   account.Followers,
	})
}

// account       godoc
// @Summary      Set account information by id
// @Tags         account
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Param        account    body    handler.accountResponce    true   	"account info"
// @Param        id         path    int                     true   	"account id"
// @Success      200   {object}   handler.accountResponce
// @Router       /api/account/{id} [post]
func (h *Handler) editAccountById(c *gin.Context) {
	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || (currentUser.Access != "admin") {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	account, err := h.repo.GetAccount(id)
	if err != nil {
		authorized, err := h.services.AuthorizeById(id)
		if (err == nil) && (authorized.Id == id) {
			account, err = h.repo.CreateAccount(id)
		}
		if err != nil {
			util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	var input accountResponce
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account.Description = input.Description
	account.DateOfBirth = input.DateOfBirth
	account.Followers = input.Followers
	account.Follows = input.Follows
	h.repo.SaveAccount(account)

	c.JSON(http.StatusOK, input)
}

// account         godoc
// @Summary      Delete account by id
// @Tags         account
// @Produce      json
// @Param        id  path      int  true  "id of user"
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200
// @Router       /api/account/{id} [delete]
func (h *Handler) delete(c *gin.Context) {
	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || (currentUser.Access != "admin") {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	_, err = h.repo.GetAccount(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.repo.DeleteAccount(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "{ \"status\": \"ok\" }")
}

// confidential          godoc
// @Summary      Get all confidential settings by id
// @Tags         confidential
// @Produce      json
// @Param        id  path      int  true  "id of user"
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200
// @Router       /api/confidential/{id} [get]
func (h *Handler) getConfidentialById(c *gin.Context) {
	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || (currentUser.Access != "admin") {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.repo.GetAccount(id)
	if err != nil {
		authorized, err := h.services.AuthorizeById(id)
		if (err == nil) && (authorized.Id == id) {
			account, err = h.repo.CreateAccount(id)
		}
		if err != nil {
			util.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}
	c.JSON(http.StatusOK, account.Confidentials)
}

// confidential          godoc
// @Summary      Set confidential settings by id
// @Tags         confidential
// @Produce      json
// @Param        id  path      int  true  "id of user"
// @Param        Authorization    header    string                           true   	"JWT Bearer token (authorization)"
// @Param        confidencial     body      string    true   	"Confidencial settings"
// @Success      200   {object}  esaccount.AuthorizedUser
// @Router       /api/confidential/{id} [post]
func (h *Handler) setConfidentialById(c *gin.Context) {
	current, _ := c.Get(userCtx)
	currentUser, err := h.services.AuthorizeById(current.(int64))
	if (err != nil) || (currentUser.Access != "admin") {
		util.NewErrorResponse(c, http.StatusForbidden, err.Error())
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	account, err := h.repo.GetAccount(id)
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	var input esaccount.ConfidentialSettings
	if err := c.BindJSON(&input); err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	for k, v := range input {
		account.Confidentials[k] = v
	}
	account.Confidentials["picture_url"] = esaccount.All

	h.repo.SaveAccount(account)

	c.JSON(http.StatusOK, account)
}

// following       godoc
// @Summary      Follow to user
// @Tags         following
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200
// @Router       /api/following/follow/{id} [put]
func (h *Handler) follow(c *gin.Context) {
	id, _ := c.Get(userCtx)
	account, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	id, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	followed, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if !slices.Contains(followed.Followers, account.Id) {
		followed.Followers = append(followed.Followers, account.Id)
		h.repo.SaveAccount(followed)
	}

	if !slices.Contains(account.Follows, followed.Id) {
		account.Follows = append(followed.Follows, followed.Id)
		h.repo.SaveAccount(account)
	}

	c.JSON(http.StatusOK, "{ \"status\": \"ok\" }")
}

// following       godoc
// @Summary      Unfollow user
// @Tags         following
// @Produce      json
// @Param        Authorization    header    string    true   	"JWT Bearer token (authorization)"
// @Success      200
// @Router       /api/following/unfollow/{id} [put]
func (h *Handler) unfollow(c *gin.Context) {
	id, _ := c.Get(userCtx)
	account, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	id, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		util.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	followed, err := h.repo.GetAccount(id.(int64))
	if err != nil {
		util.NewErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	followed.Followers = removeElement(followed.Followers, account.Id)
	account.Follows = removeElement(followed.Follows, followed.Id)

	h.repo.SaveAccount(followed)
	h.repo.SaveAccount(account)

	c.JSON(http.StatusOK, "{ \"status\": \"ok\" }")
}

func removeElement(slice []int64, value int64) []int64 {
	// Find the index of the value
	for i, v := range slice {
		if v == value {
			// Remove the element by slicing around it
			return append(slice[:i], slice[i+1:]...)
		}
	}
	// If the value is not found, return the original slice
	return slice
}
