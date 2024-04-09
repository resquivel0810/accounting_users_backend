package main

import (
	"backend/cmd/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type UserPayload struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PWD       string `json:"pwd"`
	Role      string `json:"role"`
	Account   string `json:"account"`
	Picture   string `json:"picture"`
	Status    string `json:"status"`
	Confirmed string `json:"confirmed"`
	Token     string `json:"token"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
}

func (app *application) Register(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	var payload UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errJSON(w, err)
	}
	user.Username = payload.Username
	user.Name = payload.Name
	user.LastName = payload.LastName
	user.Email = payload.Email
	user.PWD = payload.PWD
	user.Role, _ = strconv.Atoi(payload.Role)
	user.Account, _ = strconv.Atoi(payload.Account)
	user.Picture = payload.Picture
	user.Confirmed, _ = strconv.Atoi(payload.Confirmed)

	aux, err1 := app.models.DB.Register(&user)
	if err1 != nil {
		app.errJSON(w, err1)
	}
	err = app.writeJSON(w, http.StatusOK, aux, "user")
	if err != nil {
		app.errJSON(w, err)
	}

}
func (app *application) googleAuth(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	var payload UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errJSON(w, err)
	}
	user.Username = payload.Username
	user.Name = payload.Name
	user.LastName = payload.LastName
	user.Email = payload.Email
	user.PWD = payload.PWD
	user.Role, _ = strconv.Atoi(payload.Role)
	user.Account, _ = strconv.Atoi(payload.Account)
	user.Picture = payload.Picture
	user.Confirmed, _ = strconv.Atoi(payload.Confirmed)
	var id, _ = app.models.DB.VerifyEmail(user.Email)
	if id > 0 {
		aux, _ := app.models.DB.Get(id)
		token, err := app.models.DB.GnerateToken(id, app.config.jwt.secret)
		if err != nil {
			app.errJSON(w, err)
		}
		aux.Token = token
		if aux.Status == 0 {
			app.models.DB.UpdateStatus(aux.Id)
		}
		err = app.models.DB.UpdateToken(aux.Id, aux.Token)
		if err != nil {
			app.errJSON(w, err)
		}
		err = app.writeJSON(w, http.StatusOK, aux, "user")
		if err != nil {
			app.errJSON(w, err)
		}
	} else {
		user.Username = app.models.DB.RandomUser(user.Email)
		token, err := app.models.DB.GnerateToken(id, app.config.jwt.secret)
		aux, err1 := app.models.DB.RegisterGoogle(&user)
		if err1 != nil {
			app.errJSON(w, err1)
		}
		if err != nil {
			app.errJSON(w, err)
		}
		aux.Token = token
		err = app.writeJSON(w, http.StatusOK, aux, "user")
		if err != nil {
			app.errJSON(w, err)
		}
	}

}
func (app *application) Update(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	var payload UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errJSON(w, err)
	}
	user.Id, _ = strconv.Atoi(payload.Id)
	user.Username = payload.Username
	user.Name = payload.Name
	user.LastName = payload.LastName
	user.Email = payload.Email
	user.PWD = payload.PWD
	user.Role, _ = strconv.Atoi(payload.Role)
	user.Account, _ = strconv.Atoi(payload.Account)
	user.Picture = payload.Picture
	user.Confirmed, _ = strconv.Atoi(payload.Confirmed)

	aux, err1 := app.models.DB.Update(&user)
	if err1 != nil {
		app.errJSON(w, err1)
	}
	err = app.writeJSON(w, http.StatusOK, aux, "user")
	if err != nil {
		app.errJSON(w, err)
	}

}
func (app *application) UpdateLog(w http.ResponseWriter, r *http.Request) {
	user := models.User{}
	var payload UserPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errJSON(w, err)
	}
	user.Id, _ = strconv.Atoi(payload.Id)
	user.Username = payload.Username
	user.Name = payload.Name
	user.LastName = payload.LastName
	user.Email = payload.Email
	user.PWD = payload.PWD
	user.Role, _ = strconv.Atoi(payload.Role)
	user.Account, _ = strconv.Atoi(payload.Account)
	user.Picture = payload.Picture
	user.Confirmed, _ = strconv.Atoi(payload.Confirmed)
	user.Token = payload.Token
	app.logger.Println(user.Token)
	valid, err := app.models.DB.ValidateToken(user.Token, user.Id, app.config.jwt.secret)
	if valid {
		aux, err1 := app.models.DB.Update(&user)
		if err1 != nil {
			app.errJSON(w, err1)
		}
		err = app.writeJSON(w, http.StatusOK, aux, "user")
		if err != nil {
			app.errJSON(w, err)
		}
	} else {
		//app.errJSON(w, err)
		err = app.writeJSON(w, http.StatusOK, -1, "Response")
	}

}

type DParams struct {
	Id     string `json:"id"`
	Reason string `json:"reason"`
}

func (app *application) SaveReason(w http.ResponseWriter, r *http.Request) {
	var params DParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	app.logger.Println("Id is: ", params.Id)
	app.logger.Println("Reason is: ", params.Reason)
	id, _ := strconv.Atoi(params.Id)
	res, err := app.models.DB.SaveReason(id, params.Reason)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, res, "Result")
}

type LParams struct {
	Txt string `json:"txt"`
	Pwd string `json:"pwd"`
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	var params LParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	app.logger.Println("Param is: ", params.Txt)
	app.logger.Println("Password is: ", params.Pwd)
	user, err := app.models.DB.Login(params.Txt, params.Pwd, app.config.jwt.secret)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, user, "user")
}
func (app *application) ValidatePWD(w http.ResponseWriter, r *http.Request) {
	var params LParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	app.logger.Println("Param is: ", params.Txt)
	app.logger.Println("Password is: ", params.Pwd)
	user, err := app.models.DB.ValidatePWD(params.Txt, params.Pwd, app.config.jwt.secret)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, user, "user")
}

type EmParams struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Token string `json: "token"`
}

func (app *application) changeEmail(w http.ResponseWriter, r *http.Request) {
	var params EmParams
	var id int
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	app.logger.Println("Id is: ", params.Id)
	app.logger.Println("Email is: ", params.Email)
	id, err = strconv.Atoi(params.Id)
	valid, err := app.models.DB.ValidateToken(params.Token, id, app.config.jwt.secret)
	if valid {
		res, err := app.models.DB.UpdateEmail(id, params.Email)
		if err != nil {
			app.errJSON(w, err)
		}
		err = app.writeJSON(w, http.StatusOK, res, "Response")
	} else {
		//app.errJSON(w, err)
		err = app.writeJSON(w, http.StatusOK, -1, "Response")
	}
}

type Comment struct {
	Id      string `json:"id"`
	IdUser  string `json:"iduser"`
	Comment string `json:"comment"`
	Rate    string `json:"rate"`
	Consent string `json:"consent"`
}

func (app *application) SaveComment(w http.ResponseWriter, r *http.Request) {
	var params Comment
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	app.logger.Println("Id is: ", params.Id)
	app.logger.Println("Id user is: ", params.IdUser)
	app.logger.Println("Comment is: ", params.Comment)
	app.logger.Println("Rate is: ", params.Rate)
	app.logger.Println("Consent is:", params.Consent)
	res, err := app.models.DB.SaveComment(params.IdUser, params.Comment, params.Rate, params.Consent)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, res, "Response")
}
func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errJSON(w, err)
		return
	}
	app.logger.Println("id is ", id)
	user, err := app.models.DB.Get(id)
	err = app.writeJSON(w, http.StatusOK, user, "user")
}
func (app *application) getFolioId(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	title := params.ByName("title")
	user, err := app.models.DB.FolioId(title)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "user")
}
func (app *application) commentWatched(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("id")
	err := app.models.DB.CommentWatched(id)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, "Ok", "Status")
}
func (app *application) genenerateCode(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	size, err := strconv.Atoi(params.ByName("size"))
	if err != nil {
		app.logger.Print(errors.New("invalid size parameter"))
		app.errJSON(w, err)
		return
	}
	code, err := app.models.DB.GenenerateBookCode(size)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, code, "Code")
}
func (app *application) countActiveUser(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.CountActiveusers()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "total")
}
func (app *application) countUser(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.Countusers()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "total")
}
func (app *application) mostWatchedTerms(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.mostWatchedTerms()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "Most watched terms")
}
func (app *application) countCommentsDisplay(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.CountCommentsDisplay()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "total")
}
func (app *application) countCommentsWatched(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.CountCommentsWatched()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "total")
}
func (app *application) countLostUser(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.Countlostusers()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "total")
}
func (app *application) AverageTime(w http.ResponseWriter, r *http.Request) {

	user, err := app.models.DB.AverageTime()
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, user, "total")
}
func (app *application) removeUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))

	err = app.models.DB.RemoveActiveUsers(id)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, "ok", "status")
}
func (app *application) addUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errJSON(w, err)
		return
	}

	err = app.models.DB.AddActiveUsers(id)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, "ok", "status")
}
func (app *application) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.DB.All()
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, users, "users")

}
func (app *application) getComments(w http.ResponseWriter, r *http.Request) {
	comments, err := app.models.DB.Comments()
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, comments, "users")

}
func (app *application) getCommentsCMS(w http.ResponseWriter, r *http.Request) {
	comments, err := app.models.DB.CommentsCMS()
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, comments, "users")

}
func (app *application) getCommentsUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	comments, err := app.models.DB.CommentsUser(id)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, comments, "users")

}

type DisplayParams struct {
	Id      string `json:"id"`
	Display string `json: "display"`
}

func (app *application) commentsDisplay(w http.ResponseWriter, r *http.Request) {
	var params DisplayParams
	err := json.NewDecoder(r.Body).Decode(&params)
	app.logger.Println("display is ", params.Display)
	app.logger.Println("id is ", params.Id)
	err = app.models.DB.CommentsDisplay(params.Id, params.Display)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, "Success", "Status")

}
func (app *application) confirm(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errJSON(w, err)
		return
	}
	app.logger.Println("id is ", id)
	err = app.models.DB.Confirm(id)

	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}
func (app *application) mterm(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errJSON(w, err)
		return
	}
	app.logger.Println("id is ", id)
	err = app.models.DB.Mterm(id)

	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}
func (app *application) mtermlost(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	txt := params.ByName("term")

	app.logger.Println("txt is ", txt)
	err := app.models.DB.MtermLost(txt)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}
func (app *application) msavetime(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	time := params.ByName("time")

	app.logger.Println("time is ", time)
	err := app.models.DB.MsaveTime(time)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}

type DelParams struct {
	Id    string `json:"id"`
	Token string `json: "token"`
}

func (app *application) delete(w http.ResponseWriter, r *http.Request) {
	var params DelParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	id, err := strconv.Atoi(params.Id)
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errJSON(w, err)
		return
	}
	app.logger.Println("id is ", id)
	valid, err := app.models.DB.ValidateToken(params.Token, id, app.config.jwt.secret)
	var status int
	var msg string
	if valid {
		err = app.models.DB.Delete(id, params.Token)
		status = http.StatusOK
		msg = "Success"
	} else {
		status = http.StatusForbidden
		msg = "Error"
		app.errJSON(w, err)
	}

	err = app.writeJSON(w, status, msg, "Status")
}

type BookCodeParams struct {
	Id    string `json:"id"`
	Code  string `json:"code"`
	Token string `json: "token"`
}

func (app *application) ValidateBookCode(w http.ResponseWriter, r *http.Request) {
	var params BookCodeParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	id, err := strconv.Atoi(params.Id)
	app.logger.Println("code is: ", params.Code)
	app.logger.Println("id is: ", params.Id)
	app.logger.Println("token is: ", params.Token)
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errJSON(w, err)
		return
	}
	valid, err := app.models.DB.ValidateToken(params.Token, id, app.config.jwt.secret)
	var status int
	var msg string
	if valid {
		err, st := app.models.DB.ValidateBookCodeUser(id, params.Code)
		status = http.StatusOK
		app.logger.Println(st)
		msg = "Success"
		if err != nil || st == 0 {
			status = http.StatusForbidden
			msg = "Error invalid code parameter"
		}
		app.logger.Println("token correct")
	} else {
		app.logger.Println("token incorrect")
		status = http.StatusForbidden
		msg = "Error"
		app.errJSON(w, err)
	}

	err = app.writeJSON(w, status, msg, "Status")
}
func (app *application) forgotPWD(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	txt := params.ByName("txt")
	app.logger.Println("txt is: ", txt)
	err := app.models.DB.ForgotPassword(txt)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}

func (app *application) idFromCode(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	code := params.ByName("code")
	app.logger.Println("code is: ", code)
	id, err := app.models.DB.GetIdCode(code)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, id, "ID")
}
func (app *application) verifyUser(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	username := params.ByName("user")
	app.logger.Println("code is: ", username)
	id, err := app.models.DB.VerifyUser(username)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, id, "ID")
}
func (app *application) verifyEmail(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	email := params.ByName("email")
	app.logger.Println("code is: ", email)
	id, err := app.models.DB.VerifyEmail(email)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, id, "ID")
}

func (app *application) closeCode(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	code := params.ByName("code")
	app.logger.Println("code is: ", code)
	err := app.models.DB.CloseCode(code)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}

type CParams struct {
	Id  string `json:"id"`
	Pwd string `json:"pwd"`
}
type C2Params struct {
	Id     string `json:"id"`
	Oldpwd string `json:"oldpwd"`
	Newpwd string `json:"newpwd"`
	Token  string `json:"token"`
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	var params CParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	id, err := strconv.Atoi(params.Id)
	pwd := params.Pwd
	app.logger.Println("Param is: ", id)
	app.logger.Println("Password is: ", pwd)

	err = app.models.DB.ChangePassword(id, pwd)
	if err != nil {
		app.errJSON(w, err)
	}
	err = app.writeJSON(w, http.StatusOK, "Success", "Status")
}
func (app *application) changePasswordLog(w http.ResponseWriter, r *http.Request) {
	var params C2Params
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		app.errJSON(w, err)
	}
	id, err := strconv.Atoi(params.Id)
	newpwd := params.Newpwd
	token := params.Token
	app.logger.Println("Param is: ", id)
	app.logger.Println("New password is: ", newpwd)
	valid, err := app.models.DB.ValidateToken(token, id, app.config.jwt.secret)
	if valid {
		err = app.models.DB.ChangePasswordLog(id, newpwd)
		if err != nil {
			app.errJSON(w, err)
		}
		err = app.writeJSON(w, http.StatusOK, "Success", "Status")
	} else {
		//app.errJSON(w, err)
		err = app.writeJSON(w, http.StatusOK, -1, "Response")
	}
}
func (app *application) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(2000)
	file, fileinfo, err := r.FormFile("file")
	f, err := os.OpenFile("/files/"+fileinfo.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		app.errJSON(w, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

}
