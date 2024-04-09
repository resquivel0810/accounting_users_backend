package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodGet, "/v1/user/:id", app.getUser)
	router.HandlerFunc(http.MethodGet, "/v1/users/", app.getUsers)
	router.HandlerFunc(http.MethodPost, "/v1/login/", app.Login)
	router.HandlerFunc(http.MethodPost, "/v1/validatepwd/", app.ValidatePWD)
	router.HandlerFunc(http.MethodPost, "/v1/register/", app.Register)
	router.HandlerFunc(http.MethodPost, "/v1/confirm/:id", app.confirm)
	router.HandlerFunc(http.MethodPost, "/v1/cpwd/", app.changePassword)
	router.HandlerFunc(http.MethodGet, "/v1/forgot/:txt", app.forgotPWD)
	router.HandlerFunc(http.MethodGet, "/v1/getid/:code", app.idFromCode)
	router.HandlerFunc(http.MethodGet, "/v1/closecode/:code", app.closeCode)
	router.HandlerFunc(http.MethodPost, "/v1/porfilepic/", app.uploadPhoto)
	router.HandlerFunc(http.MethodPost, "/v1/updateuser/", app.Update)
	router.HandlerFunc(http.MethodPost, "/v1/updateuserlog/", app.UpdateLog)
	router.HandlerFunc(http.MethodPost, "/v1/deleteaccount/", app.delete)
	router.HandlerFunc(http.MethodPost, "/v1/cpwdlog/", app.changePasswordLog)
	router.HandlerFunc(http.MethodGet, "/v1/verifyuser/:user", app.verifyUser)
	router.HandlerFunc(http.MethodGet, "/v1/verifyemail/:email", app.verifyEmail)
	router.HandlerFunc(http.MethodPost, "/v1/changeemail/", app.changeEmail)
	router.HandlerFunc(http.MethodPost, "/v1/googleauth/", app.googleAuth)
	router.HandlerFunc(http.MethodPost, "/v1/savereason/", app.SaveReason)
	router.HandlerFunc(http.MethodGet, "/v1/countusers/", app.countUser)
	router.HandlerFunc(http.MethodGet, "/v1/addactiveuser/:id", app.addUser)
	router.HandlerFunc(http.MethodGet, "/v1/removeactiveuser/:id", app.removeUser)
	router.HandlerFunc(http.MethodGet, "/v1/getactiveuser/", app.countActiveUser)
	router.HandlerFunc(http.MethodGet, "/v1/countlostusers/", app.countLostUser)
	router.HandlerFunc(http.MethodPost, "/v1/mterm/:id", app.mterm)
	router.HandlerFunc(http.MethodPost, "/v1/mtermlost/:term", app.mtermlost)
	router.HandlerFunc(http.MethodGet, "/v1/timesave/:time", app.msavetime)
	router.HandlerFunc(http.MethodGet, "/v1/averagetime/", app.AverageTime)
	router.HandlerFunc(http.MethodPost, "/v1/savecomment/", app.SaveComment)
	router.HandlerFunc(http.MethodGet, "/v1/comments/", app.getComments)
	router.HandlerFunc(http.MethodGet, "/v1/folioid/:title", app.getFolioId)
	router.HandlerFunc(http.MethodGet, "/v1/generatebookcode/:size", app.genenerateCode)
	router.HandlerFunc(http.MethodPost, "/v1/validatebookcode/", app.ValidateBookCode)
	router.HandlerFunc(http.MethodGet, "/v1/commentsUser/:id", app.getCommentsUser)
	router.HandlerFunc(http.MethodGet, "/v1/commentsCMS/", app.getCommentsCMS)
	router.HandlerFunc(http.MethodPost, "/v1/commentsDisplay/", app.commentsDisplay)
	router.HandlerFunc(http.MethodGet, "/v1/countCommentsDisplay/", app.countCommentsDisplay)
	router.HandlerFunc(http.MethodGet, "/v1/countCommentsWatched/", app.countCommentsWatched)
	router.HandlerFunc(http.MethodGet, "/v1/commentWatched/:id", app.commentWatched)
	router.HandlerFunc(http.MethodGet, "/v1/mostWatchedTerms", app.mostWatchedTerms)
	return router
}
