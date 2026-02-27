package main

import (
	"fmt"
	"net/http"

	_ "github.com/tidwall/gjson"
	"github.com/julienschmidt/httprouter"

	"github.com/ShiftLeftSecurity/shiftleft-go-demo/setting"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/setup"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/user"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/config"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/util/middleware"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/vulnerability/csa"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/vulnerability/idor"
	pathTraversal "github.com/ShiftLeftSecurity/shiftleft-go-demo/vulnerability/path-traversal"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/vulnerability/sqli"
	"github.com/ShiftLeftSecurity/shiftleft-go-demo/vulnerability/xss"
)

const (
	banner = `
     ÛÛÛÛÛÛÛÛÛ           ÛÛÛÛÛ   ÛÛÛÛÛ ÛÛÛÛÛ   ÛÛÛ   ÛÛÛÛÛ   ÛÛÛÛÛÛÛÛÛ
    ÛÛÛ°°°°°ÛÛÛ         °°ÛÛÛ   °°ÛÛÛ °°ÛÛÛ   °ÛÛÛ  °°ÛÛÛ   ÛÛÛ°°°°°ÛÛÛ
   ÛÛÛ     °°°   ÛÛÛÛÛÛ  °ÛÛÛ    °ÛÛÛ  °ÛÛÛ   °ÛÛÛ   °ÛÛÛ  °ÛÛÛ    °ÛÛÛ
  °ÛÛÛ          ÛÛÛ°°ÛÛÛ °ÛÛÛ    °ÛÛÛ  °ÛÛÛ   °ÛÛÛ   °ÛÛÛ  °ÛÛÛÛÛÛÛÛÛÛÛ
  °ÛÛÛ    ÛÛÛÛÛ°ÛÛÛ °ÛÛÛ °°ÛÛÛ   ÛÛÛ   °°ÛÛÛ  ÛÛÛÛÛ  ÛÛÛ   °ÛÛÛ°°°°°ÛÛÛ
  °°ÛÛÛ  °°ÛÛÛ °ÛÛÛ °ÛÛÛ  °°°ÛÛÛÛÛ°     °°°ÛÛÛÛÛ°ÛÛÛÛÛ°    °ÛÛÛ    °ÛÛÛ
   °°ÛÛÛÛÛÛÛÛÛ °°ÛÛÛÛÛÛ     °°ÛÛÛ         °°ÛÛÛ °°ÛÛÛ      ÛÛÛÛÛ   ÛÛÛÛÛ
     °°°°°°°°°   °°°°°°       °°°           °°°   °°°      °°°°°   °°°°° `
)

//index and set cookie

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	util.SetCookieLevel(w, r, "low") //set cookie Level default to low

	data := make(map[string]interface{})
	data["title"] = "Index"

	util.SafeRender(w, r, "template.index", data)
}

func main() {

	fmt.Println(banner)

	mw := middleware.New()
	router := httprouter.New()

	// Create all module instances (including false positives)
	// Must create both regular and FP modules from same package before variable shadows package name
	userModule := user.New()
	sqlI := sqli.New()
	sqlIFP := sqli.NewFalsePositive()
	xssModule := xss.New()
	xssFP := xss.NewFalsePositive()
	idorModule := idor.New()
	idorFP := idor.NewFalsePositive()
	csaModule := csa.New()
	pathTraversalModule := pathTraversal.New()
	pathTraversalFP := pathTraversal.NewFalsePositiveTraversal()
	setupModule := setup.New()
	settingModule := setting.New()

	router.ServeFiles("/public/*filepath", http.Dir("public/"))
	router.GET("/", mw.LoggingMiddleware(mw.AuthCheck(indexHandler)))
	router.GET("/index", mw.LoggingMiddleware(mw.DetectSQLMap(mw.AuthCheck(indexHandler))))

	// Register all routes (including false positives)
	userModule.SetRouter(router)
	sqlI.SetRouter(router)
	sqlIFP.SetRouter(router)
	xssModule.SetRouter(router)
	xssFP.SetRouter(router)
	idorModule.SetRouter(router)
	idorFP.SetRouter(router)
	csaModule.SetRouter(router)
	pathTraversalModule.SetRouter(router)
	pathTraversalFP.SetRouter(router)
	setupModule.SetRouter(router)
	settingModule.SetRouter(router)

	s := http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	fmt.Printf("Server running at port %s\n", s.Addr)
	fmt.Printf("Open this url %s on your browser to access GoVWA", config.Fullurl)
	fmt.Println("")
	s.ListenAndServe()

}
