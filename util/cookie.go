package util

import (
	"net/http"
	"time"
)

func SetCookieLevel(w http.ResponseWriter, r *http.Request, cookievalue string) {

	level := cookievalue
	if level == "" {
		level = "low"
	}
	SetCookie(w, "Level", level)

}

func CheckLevel(r *http.Request) bool {
	level := GetCookie(r, "Level")
	if level == "" || level == "low" {
		return false //set default level to low
	} else if level == "high" {
		return true //level == high
	} else {
		return false // level == low
	}
}

/* cookie setter getter */

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := http.Cookie{
		//Path : "/",
		//Domain : "localhost",
		Name:  name,
		Value: value,
	}
	http.SetCookie(w, &cookie)
var (
    // Fixed: Prepared statement for query plan caching and additional SQL injection protection
    getProfileStmt *sql.Stmt
    stmtMutex      sync.RWMutex
    
    // Fixed: Rate limiter for preventing enumeration attacks
    requestTracker = make(map[string]*requestRecord)
    trackerMutex   sync.RWMutex
)

type requestRecord struct {
    count     int
    firstSeen time.Time
}

// InitPreparedStatements initializes prepared statements for database queries
// This function should be called during application startup after DB connection is established
// Fixed: Implements prepared statement caching for performance and security
func InitPreparedStatements(db *sql.DB) error {
    stmtMutex.Lock()
    defer stmtMutex.Unlock()
    
    getProfileSql := `SELECT p.user_id, p.full_name, p.city, p.phone_number 
                      FROM Profile as p, Users as u 
                      WHERE p.user_id = u.id 
                      AND u.id = ?`
    
    var err error
    getProfileStmt, err = db.Prepare(getProfileSql)
    if err != nil {
        return err
    }
    return nil
}

// CleanupPreparedStatements closes prepared statements during application shutdown
// Fixed: Proper resource cleanup for prepared statements
func CleanupPreparedStatements() {
    stmtMutex.Lock()
    defer stmtMutex.Unlock()
    
    if getProfileStmt != nil {
        getProfileStmt.Close()
    }
}

    id, err := strconv.Atoi(uid)
    if err != nil {
        return errors.New("invalid UID format: must be numeric")
    }
    if id <= 0 || id > 999999999 {
        return errors.New("invalid UID range")
    }
    return nil
}

		http.SetCookie(w, cookie)
	}
}
