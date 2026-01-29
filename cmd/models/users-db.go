package models

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/pascaldekloe/jwt"
	"github.com/xuri/excelize/v2"
)

type DBModel struct {
	DB *sql.DB
}

func (m DBModel) CreateExcel(codes []string, numcodes int) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Set value of a cell.
	f.SetCellValue("Sheet1", "A1", "Code")
	for i := 2; i <= numcodes+1; i++ {
		aux := strconv.Itoa(i)
		f.SetCellValue("Sheet1", ("A" + aux), codes[i-2])
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs("codes.xlsx"); err != nil {
		fmt.Println(err)
	}
}
func (m DBModel) ForgotPassword(txt string) error {
	id, err := m.GetIdUser(txt)
	if err != nil {
		return err
	}
	if id != "" {
		user, err := m.Get(id)
		if err != nil {
			return err
		}
		rand.Seed(time.Now().UnixNano())
		code := randomString(50)
		var q = "INSERT INTO user_code (id_user,code) values ('" + id + "','" + code + "')"
		_, err = m.DB.Exec(q)

		if err != nil {
			return err
		}
		user.Code = code
		if err := mail_send(user, "Reset Password", "forgot.html"); err != nil {
			log.Printf("Warning: Failed to send password reset email to %s: %v", user.Email, err)
		}
	}
	return nil
}
func (m DBModel) GenenerateBookCode(numcodes int) (*[]string, error) {
	//code := randomCode(15)
	codes := []string{}
	for i := 0; i < numcodes; i++ {
		code := randomCode(15)
		b, err := m.ValidateBookCode(code)
		if err != nil {
			return nil, err
		} else {
			if !b {
				var q = "INSERT INTO bookcode  values (0,md5('" + code + "'),1,0)"
				_, err = m.DB.Exec(q)
				if err != nil {
					return nil, err
				} else {
					codes = append(codes, code)
				}
			}
		}
	}
	m.CreateExcel(codes, numcodes)
	return &codes, nil
}
func (m DBModel) ValidateBookCode(code string) (bool, error) {
	var q = "SELECT code FROM bookcode WHERE code=md5('" + code + "')"
	row, err := m.DB.Query(q)
	if err != nil {
		return false, err
	}
	var aux string
	aux = ""
	for row.Next() {
		err := row.Scan(&aux)
		if err != nil {
			panic(err)
		}
	}
	if aux != "" {
		return true, nil
	}
	return false, nil
}
func (m DBModel) IsBookCodeUsed(code string) (bool, error) {
	var q = "SELECT code FROM bookcode WHERE code=md5('" + code + "') AND status=1"
	row, err := m.DB.Query(q)
	if err != nil {
		return false, err
	}
	var aux string
	aux = ""
	for row.Next() {
		err := row.Scan(&aux)
		if err != nil {
			panic(err)
		}
	}
	if aux != "" {
		return true, nil
	}
	return false, nil
}
func (m DBModel) GetIdCode(code string) (string, error) {
	var q = "SELECT id_user FROM user_code WHERE code='" + code + "' and status=1 and NOW()<expiration"
	row, err := m.DB.Query(q)
	var id string
	var found bool
	if err != nil {
		return "", err
	}
	for row.Next() {
		err := row.Scan(&id)
		if err != nil {
			panic(err)
		}
		found = true
		break
	}
	if !found {
		return "", nil
	}
	return id, nil
}
func (m DBModel) VerifyUser(user string) (string, error) {
	var q = "SELECT id FROM users WHERE username='" + user + "'"
	row, err := m.DB.Query(q)
	var id string
	if err != nil {
		return "", err
	}
	for row.Next() {
		err := row.Scan(&id)
		if err != nil {
			panic(err)
		}

	}

	return id, nil
}
func (m DBModel) VerifyEmail(email string) (string, error) {
	var q = "SELECT id FROM users WHERE email='" + email + "'"
	row, err := m.DB.Query(q)
	var id string
	if err != nil {
		return "", err
	}
	for row.Next() {
		err := row.Scan(&id)
		if err != nil {
			panic(err)
		}

	}

	return id, nil
}
func (m DBModel) CloseCode(code string) error {
	var q = "UPDATE user_code SET status=0 where code='" + code + "'"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil
}
func (m DBModel) CloseAllCodes(id string) error {
	var q = "UPDATE user_code SET status=0 where id_user='" + id + "'"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil
}
func (m DBModel) GetIdUser(txt string) (string, error) {
	var q = "SELECT id FROM users WHERE username='" + txt + "' or email='" + txt + "'"
	row, err := m.DB.Query(q)
	var id string
	if err != nil {
		return "", err
	}
	for row.Next() {
		err := row.Scan(&id)
		if err != nil {
			panic(err)
		}

	}

	return id, nil
}
func (m DBModel) Register(user *User) (*User, error) {
	// Generar UUID para el nuevo usuario
	user.Id = uuid.New().String()
	var q = "INSERT INTO users (id,username,name,last_name,email,pwd,profile_picture_url,status,token,role,account) values ('" + user.Id + "','" + user.Username + "','" + user.Name + "','" + user.LastName + "','" + user.Email + "',MD5('" + user.PWD + "'),'" + user.Picture + "',1,''," + strconv.Itoa(user.Role) + "," + strconv.Itoa(user.Account) + ")"
	_, err := m.DB.Exec(q)

	if err != nil {
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	code := randomString(50)
	q = "INSERT INTO user_code (id_user,code) values ('" + user.Id + "','" + code + "')"
	_, err = m.DB.Exec(q)

	if err != nil {
		return nil, err
	}
	user.Code = code
	if err := mail_send(user, "Confirm your email address", "confirm.html"); err != nil {
		log.Printf("Warning: Failed to send confirmation email to %s: %v", user.Email, err)
		// No retornamos error para que el registro se complete aunque falle el email
	}
	return user, nil
}
func (m DBModel) RegisterGoogle(user *User) (*User, error) {
	// Generar UUID para el nuevo usuario
	user.Id = uuid.New().String()
	var q = "INSERT INTO users (id,username,name,last_name,email,pwd,profile_picture_url,status,token,role,account,email_conf) values ('" + user.Id + "','" + user.Username + "','" + user.Name + "','" + user.LastName + "','" + user.Email + "',MD5('" + user.PWD + "'),'" + user.Picture + "',1,'" + user.Token + "'," + strconv.Itoa(user.Role) + "," + strconv.Itoa(user.Account) + ",1)"
	_, err := m.DB.Exec(q)

	if err != nil {
		return nil, err
	}
	if err := mail_send(user, "Welcome", "welcome.html"); err != nil {
		log.Printf("Warning: Failed to send welcome email to %s: %v", user.Email, err)
	}
	return user, nil
}
func (m DBModel) Update(user *User) (*User, error) {
	var txt = ""
	if user.PWD != "" {
		txt = "pwd=MD5('" + user.PWD + "'),"
	}
	aux, _ := m.Get(user.Id)
	var q = "UPDATE users SET " + txt + " username='" + user.Username + "',name='" + user.Name + "',last_name='" + user.LastName + "',email='" + user.Email + "',profile_picture_url='" + user.Picture + "',status=1 WHERE id='" + user.Id + "'"
	_, err := m.DB.Exec(q)

	if err != nil {
		return nil, err
	}
	if aux.Status == 0 {
		if err := mail_send(user, "Good to see you're back!", "userback.html"); err != nil {
			log.Printf("Warning: Failed to send welcome back email to %s: %v", user.Email, err)
		}
	}
	return user, nil
}
func (m DBModel) UpdateEmail(id string, email string) (string, error) {
	var user *User
	user, _ = m.Get(id)
	if user == nil {
		return "", fmt.Errorf("user not found")
	}
	
	// Guardar el email antiguo antes de actualizar
	oldEmail := user.Email
	
	// Enviar notificación al email antiguo ANTES de actualizar la BD
	if err := mail_send(user, "You've changed your email", "oldemail.html"); err != nil {
		log.Printf("Warning: Failed to send old email notification to %s: %v", oldEmail, err)
	}
	
	// Actualizar email en la base de datos
	var q = "UPDATE users SET email='" + email + "' WHERE id='" + id + "'"
	_, err := m.DB.Exec(q)
	if err != nil {
		return "", err
	}
	
	// Actualizar el objeto user con el nuevo email
	user.Email = email
	
	// Enviar notificación al nuevo email
	if err := mail_send(user, "You've changed your email", "newemail.html"); err != nil {
		log.Printf("Warning: Failed to send new email notification to %s: %v", email, err)
	}
	return id, nil
}
func (m DBModel) SaveComment(iduser string, comment string, rate string, consent string) (int64, error) {

	var q = "INSERT into feedback (id,iduser,rate,comment,consent) VALUES  (0,'" + iduser + "'," + rate + ",'" + comment + "'," + consent + ")"
	res, err := m.DB.Exec(q)
	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastId, nil
}
func (m DBModel) Login(param string, password string, secret string) (*User, error) {
	var q = "SELECT * FROM users WHERE (username='" + param + "' AND pwd=MD5('" + password + "') AND status=1) OR (email='" + param + "' AND pwd=MD5('" + password + "') AND status=1)"
	row, err := m.DB.Query(q)

	if err != nil {
		return nil, err
	}
	var user User
	for row.Next() {
		var id string
		var status, role, account, conf, active int
		var username, name, lastname, email, pwd, pic, token string
		var created, updated time.Time
		err := row.Scan(&id, &username, &name, &lastname, &email, &pwd, &pic, &status, &token, &role, &account, &conf, &created, &updated, &active)
		if err != nil {
			panic(err)
		}
		user.Id = id
		user.Username = username
		user.Name = name
		user.LastName = lastname
		user.Email = email
		user.PWD = pwd
		user.Picture = pic
		user.Status = status
		user.Account = account
		user.Status = status
		user.Role = role
		user.Confirmed = conf
		user.Created = created
		user.Updated = updated
		user.Active = active
	}
	if user.Confirmed == 0 && user.Id != "" {
		m.CloseAllCodes(user.Id)
		rand.Seed(time.Now().UnixNano())
		code := randomString(50)
		q = "INSERT INTO user_code (id_user,code) values ('" + user.Id + "','" + code + "')"
		_, err = m.DB.Exec(q)

		if err != nil {
			return nil, err
		}
		user.Code = code
		if err := mail_send(&user, "Don't forget to confirm your account", "forgetconfirm.html"); err != nil {
			log.Printf("Warning: Failed to send confirmation reminder email to %s: %v", user.Email, err)
		}
	}
	user.Token, _ = m.GnerateToken(user.Id, secret)
	m.SaveToken(user.Id, user.Token)
	return &user, nil
}
func (m DBModel) ValidatePWD(param string, password string, secret string) (*User, error) {
	var q = "SELECT * FROM users WHERE (username='" + param + "' AND pwd=MD5('" + password + "') AND status=1) OR (email='" + param + "' AND pwd=MD5('" + password + "') AND status=1)"
	row, err := m.DB.Query(q)

	if err != nil {
		return nil, err
	}
	var user User
	for row.Next() {
		var id string
		var status, role, account, conf, active int
		var username, name, lastname, email, pwd, pic, token string
		var created, updated time.Time
		err := row.Scan(&id, &username, &name, &lastname, &email, &pwd, &pic, &status, &token, &role, &account, &conf, &created, &updated, &active)
		if err != nil {
			panic(err)
		}
		user.Id = id
		user.Username = username
		user.Name = name
		user.LastName = lastname
		user.Email = email
		user.PWD = pwd
		user.Picture = pic
		user.Status = status
		user.Account = account
		user.Status = status
		user.Role = role
		user.Confirmed = conf
		user.Created = created
		user.Updated = updated
		user.Active = active
	}
	return &user, nil
}
func (m DBModel) AddActiveUsers(id string) error {
	var q = "UPDATE users SET active = 1 where id='" + id + "'"
	_, err := m.DB.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
func (m DBModel) RemoveActiveUsers(id string) error {
	var q = "UPDATE users SET active = 0 where id='" + id + "'"
	_, err := m.DB.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
func (m DBModel) UpdateToken(id string, token string) error {
	var q = "UPDATE users SET token = '" + token + "' where id='" + id + "'"
	_, err := m.DB.Exec(q)
	if err != nil {
		return err
	}
	return nil
}
func (m DBModel) UpdateStatus(id string) error {
	var q = "UPDATE users SET status = 1 where id='" + id + "'"
	_, err := m.DB.Exec(q)
	if err != nil {
		return err
	}
	user, _ := m.Get(id)
	mail_send(user, "Good to see you're back!", "userback.html")
	return nil
}
func (m DBModel) GnerateToken(id string, secret string) (string, error) {
	var claims jwt.Claims
	claims.Subject = id
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(730 * time.Hour))
	claims.Issuer = "accounting-a-z.ch"
	claims.Audiences = []string{"accounting-a-z.ch"}
	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(secret))
	if err != nil {
		return "", err
	}
	return string(jwtBytes), nil
}
func (m DBModel) ValidateToken(tokenString string, id string, secret string) (bool, error) {
	var valid bool
	claims, err := jwt.HMACCheck([]byte(tokenString), []byte(secret))
	if err != nil {
		log.Print("credentials rejected: ", err)
		return false, err
	}
	err = claims.AcceptTemporal(time.Now(), time.Second)
	if err != nil {
		log.Print("credential constraints violated: ", err)
		return false, err
	}
	if claims.Valid(time.Now()) {
		aux, _ := m.CompareDBtoken(id, tokenString)
		if aux {
			valid = true
		} else {
			valid = false
		}
	}
	return valid, nil
}
func (m DBModel) SaveToken(id string, token string) error {
	var q = "UPDATE users SET token='" + token + "' where id='" + id + "'"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil
}
func (m DBModel) Countusers() (int, error) {
	var q = "SELECT count(*) FROM users"
	row, err := m.DB.Query(q)

	if err != nil {
		return 0, err
	}
	var total int
	for row.Next() {

		err := row.Scan(&total)
		if err != nil {
			panic(err)
		}

	}
	//mail_send(&user, "Confirm your email address", "confirm.html")
	return total, nil
}
// func (m DBModel) mostWatchedTerms() (int, error) {
// 	var q = "SELECT idterm, COUNT(idterm) AS MOST_FREQUENT FROM preview.metricsterm GROUP BY idterm ORDER BY COUNT(idterm) DESC LIMIT 10"
// 	row, err := m.DB.Query(q)

// 	if err != nil {
// 		return 0, err
// 	}
// 	var total int
// 	for row.Next() {

// 		err := row.Scan(&total)
// 		if err != nil {
// 			panic(err)
// 		}

// 	}
// 	//mail_send(&user, "Confirm your email address", "confirm.html")
// 	return total, nil
// }
func (m DBModel) CountCommentsWatched() (int, error) {
	var q = "SELECT count(*) FROM feedback WHERE watched = 0"
	row, err := m.DB.Query(q)

	if err != nil {
		return 0, err
	}
	var total int
	for row.Next() {

		err := row.Scan(&total)
		if err != nil {
			panic(err)
		}

	}
	//mail_send(&user, "Confirm your email address", "confirm.html")
	return total, nil
}
func (m DBModel) CountCommentsDisplay() (int, error) {
	var q = "SELECT count(*) FROM feedback WHERE display = 1"
	row, err := m.DB.Query(q)

	if err != nil {
		return 0, err
	}
	var total int
	for row.Next() {

		err := row.Scan(&total)
		if err != nil {
			panic(err)
		}

	}
	//mail_send(&user, "Confirm your email address", "confirm.html")
	return total, nil
}
func (m DBModel) Countlostusers() (int, error) {
	var q = "SELECT count(*) FROM users where status=0"
	row, err := m.DB.Query(q)

	if err != nil {
		return 0, err
	}
	var total int
	for row.Next() {

		err := row.Scan(&total)
		if err != nil {
			panic(err)
		}

	}
	return total, nil
}
func (m DBModel) AverageTime() (float64, error) {
	var q = "SELECT AVG(time) FROM time_used"
	row, err := m.DB.Query(q)

	if err != nil {
		return 0, err
	}
	var total float64
	for row.Next() {

		err := row.Scan(&total)
		if err != nil {
			panic(err)
		}

	}
	return total, nil
}
func (m DBModel) CountActiveusers() (int, error) {
	var q = "SELECT count(*) FROM users WHERE active=1"
	row, err := m.DB.Query(q)

	if err != nil {
		return 0, err
	}
	var total int
	for row.Next() {

		err := row.Scan(&total)
		if err != nil {
			panic(err)
		}

	}
	//mail_send(&user, "Confirm your email address", "confirm.html")
	return total, nil
}
func (m DBModel) Get(id string) (*User, error) {
	var q = "SELECT * FROM users WHERE id='" + id + "'"
	row, err := m.DB.Query(q)

	if err != nil {
		return nil, err
	}
	var user User
	for row.Next() {
		var id string
		var status, role, account, conf, active int
		var username, name, lastname, email, pwd, pic, token string
		var created, updated time.Time
		err := row.Scan(&id, &username, &name, &lastname, &email, &pwd, &pic, &status, &token, &role, &account, &conf, &created, &updated, &active)
		if err != nil {
			panic(err)
		}
		user.Id = id
		user.Username = username
		user.Name = name
		user.LastName = lastname
		user.Email = email
		user.PWD = "secret"
		user.Picture = pic
		user.Confirmed = conf
		user.Account = account
		user.Status = status
		user.Role = role
		user.Created = created
		user.Updated = updated
		user.Active = active

	}
	//mail_send(&user, "Confirm your email address", "confirm.html")
	return &user, nil
}
func (m DBModel) FolioId(title string) (int, error) {
	var db *sql.DB
	var err error
	var id int
	db, err = sql.Open("mysql", "admin:NiL9620C0n@tcp(localhost:3306)/strapi_copy?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	var q = "SELECT id FROM terms WHERE text like '" + title + "'"
	row, err := db.Query(q)
	if err != nil {
		return 0, err
	}
	if row.Next() {
		err = row.Scan(&id)
		if err != nil {
			panic(err)
		}
	}
	return id, nil
}
func (m DBModel) CompareDBtoken(id string, token string) (bool, error) {
	var q = "SELECT * FROM users WHERE token='" + token + "' AND id='" + id + "'"
	row, err := m.DB.Query(q)

	if err != nil {
		return false, err
	}
	var valid bool
	if row.Next() {
		valid = true

	} else {
		valid = false
	}
	return valid, nil
}
func (m DBModel) SaveReason(id string, reason string) (bool, error) {
	var q = "INSERT INTO DelAccount (userid,reason) values ('" + id + "','" + reason + "')"
	_, err := m.DB.Exec(q)

	if err != nil {
		return false, err
	}
	if err != nil {
		panic(err)
	}
	return true, nil
}
func (m DBModel) All() (*[]User, error) {
	row, err := m.DB.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}
	var user User
	users := []User{}
	for row.Next() {
		var id string
		var status, role, account, conf, active int
		var username, name, lastname, email, pwd, pic, token string
		var created, updated time.Time
		err := row.Scan(&id, &username, &name, &lastname, &email, &pwd, &pic, &status, &token, &role, &account, &conf, &created, &updated, &active)
		if err != nil {
			panic(err)
		}
		user.Id = id
		user.Username = username
		user.Name = name
		user.LastName = lastname
		user.Email = email
		user.PWD = "secret"
		user.Picture = pic
		user.Confirmed = conf
		user.Account = account
		user.Status = status
		user.Role = role
		user.Created = created
		user.Updated = updated
		user.Active = active

		users = append(users, user)

	}
	return &users, nil
}
func (m DBModel) Comments() (*[]Feedback, error) {

	row, err := m.DB.Query("SELECT feedback.id,users.id,name,profile_picture_url,rate,comment FROM feedback inner join users ON feedback.iduser = users.id where feedback.display=1 AND feedback.consent=1")

	if err != nil {
		return nil, err
	}
	var comment Feedback
	comments := []Feedback{}
	for row.Next() {
		var id, rate int
		var id_user, name, pic, com string
		err := row.Scan(&id, &id_user, &name, &pic, &rate, &com)
		if err != nil {
			panic(err)
		}
		comment.Id = id
		comment.IdUser = id_user
		comment.Name = name
		comment.Picture = pic
		comment.Rate = rate
		comment.Comment = com

		comments = append(comments, comment)

	}
	return &comments, nil
}
func (m DBModel) CommentsCMS() (*[]Feedback, error) {

	row, err := m.DB.Query("SELECT feedback.id,users.id,name,profile_picture_url,rate,comment,consent,display,feedback.created_at,watched FROM feedback inner join users ON feedback.iduser = users.id ")

	if err != nil {
		return nil, err
	}
	var comment Feedback
	comments := []Feedback{}
	for row.Next() {
		var id, rate, consent, display, watched int
		var id_user string
		var created time.Time
		var name, pic, com string
		err := row.Scan(&id, &id_user, &name, &pic, &rate, &com, &consent, &display, &created, &watched)
		if err != nil {
			panic(err)
		}
		comment.Id = id
		comment.IdUser = id_user
		comment.Name = name
		comment.Picture = pic
		comment.Rate = rate
		comment.Comment = com
		comment.Consent = consent
		comment.Display = display
		comment.Watched = watched
		comment.Created = created

		comments = append(comments, comment)

	}
	return &comments, nil
}
func (m DBModel) CommentsUser(id string) (*[]Feedback, error) {
	row, err := m.DB.Query(("SELECT feedback.id,users.id,name,profile_picture_url,rate,comment FROM feedback inner join users ON feedback.iduser = users.id where feedback.idUser='" + id + "'"))

	if err != nil {
		return nil, err
	}
	var comment Feedback
	comments := []Feedback{}
	for row.Next() {
		var id, rate int
		var id_user, name, pic, com string
		err := row.Scan(&id, &id_user, &name, &pic, &rate, &com)
		if err != nil {
			panic(err)
		}
		comment.Id = id
		comment.IdUser = id_user
		comment.Name = name
		comment.Picture = pic
		comment.Rate = rate
		comment.Comment = com

		comments = append(comments, comment)

	}
	return &comments, nil
}
func (m DBModel) CommentsDisplay(id string, display string) error {
	var q = "UPDATE feedback SET display=" + display + " WHERE id=" + id + " and consent=1"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}

	return nil
}
func (m DBModel) CommentWatched(id string) error {
	var q = "UPDATE feedback SET watched=1 WHERE id=" + id
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}

	return nil
}
func (m DBModel) Confirm(id string) error {
	var q = "UPDATE users SET email_conf=1 where id='" + id + "'"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil

}
func (m DBModel) Mterm(id int) error {
	var q = "INSERT INTO metricsterm (idterm) values (" + strconv.Itoa(id) + ")"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil

}
func (m DBModel) MtermLost(txt string) error {
	var q = "INSERT INTO sugestedterms (term) values ('" + txt + "')"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil

}
func (m DBModel) MsaveTime(time string) error {
	var q = "INSERT INTO time_used (time) values (" + time + ")"
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	return nil

}
func (m DBModel) Delete(id string, token string) error {
	var q = "UPDATE users SET status=0 where id='" + id + "' AND token='" + token + "'"
	var user *User
	_, err := m.DB.Exec(q)

	if err != nil {
		return err
	}
	user, _ = m.Get(id)
	if err := mail_send(user, "We hope you back soon", "userdel.html"); err != nil {
		log.Printf("Warning: Failed to send account deletion email to %s: %v", user.Email, err)
	}
	return nil

}
func (m DBModel) ValidateBookCodeUser(id string, code string) (error, int) {
	b, err := m.ValidateBookCode(code)
	var st int
	if err != nil {
		return err, 0
	}
	if b {
		var q = "UPDATE bookcode SET status = 0,id_user='" + id + "' where code=md5('" + code + "') and status=1"
		stmt, err := m.DB.Exec(q)
		if err != nil {
			return err, 0
		}
		rows, _ := stmt.RowsAffected()
		st = int(rows)
		q = "UPDATE users SET account = 2 where id='" + id + "'"
		_, err = m.DB.Exec(q)
		if err != nil {
			return err, 0
		}
	}
	return nil, st

}
func (m DBModel) ChangePassword(id string, password string) error {
	var q = "UPDATE users SET pwd = MD5('" + password + "') where id='" + id + "'"
	_, err := m.DB.Exec(q)
	var user *User
	if err != nil {
		return err
	}
	user, _ = m.Get(id)
	if err := mail_send(user, "Password Change", "pwdreset.html"); err != nil {
		log.Printf("Warning: Failed to send password reset email to %s: %v", user.Email, err)
	}
	return nil
}
func (m DBModel) ChangePasswordLog(id string, newpwd string) error {
	var q = "UPDATE users SET pwd = MD5('" + newpwd + "') where id='" + id + "'"
	_, err := m.DB.Exec(q)
	var user *User
	if err != nil {
		return err
	}
	user, _ = m.Get(id)
	if err := mail_send(user, "Password Change", "pwdchange.html"); err != nil {
		log.Printf("Warning: Failed to send password change email to %s: %v", user.Email, err)
	}
	return nil
}
func mail_send(user *User, subject string, temp string) error {
	from := mail.Address{"Accounting A-Z", "no-reply@accounting-a-z.ch"}
	to := mail.Address{user.Name, user.Email}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-type"] = `text/html; charset="UTF-8"`
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	t, err := template.ParseFiles(temp)
	if err != nil {
		log.Printf("Error parsing email template %s: %v", temp, err)
		return err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, user)
	if err != nil {
		log.Printf("Error executing email template %s: %v", temp, err)
		return err
	}
	message += buf.String()
	servername := "server29.hostfactory.ch:465"
	host := "server29.hostfactory.ch"
	auth := smtp.PlainAuth("", "no-reply@accounting-a-z.ch", "theStrongestOne1", host)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", servername, tlsConfig)
	if err != nil {
		log.Printf("Error connecting to SMTP server: %v", err)
		return err
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Printf("Error creating SMTP client: %v", err)
		return err
	}
	defer client.Quit()
	err = client.Auth(auth)
	if err != nil {
		log.Printf("Error authenticating SMTP: %v", err)
		return err
	}
	err = client.Mail(from.Address)
	if err != nil {
		log.Printf("Error setting sender: %v", err)
		return err
	}
	err = client.Rcpt(to.Address)
	if err != nil {
		log.Printf("Error setting recipient: %v", err)
		return err
	}
	w, err := client.Data()
	if err != nil {
		log.Printf("Error getting data writer: %v", err)
		return err
	}
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Printf("Error writing email body: %v", err)
		w.Close()
		return err
	}
	err = w.Close()
	if err != nil {
		log.Printf("Error closing data writer: %v", err)
		return err
	}
	return nil
}

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
func (m DBModel) RandomUser(email string) string {
	var num int
	for i := 0; i < len(email); i++ {
		if email[i] == '@' {
			num = i
		}
	}
	name := email[0:num]
	return name + randomString(4)
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}
func randomCode(n int) string {
	var letterRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
